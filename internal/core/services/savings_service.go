package services

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type SavingsService struct {
    savingsGoalRepo     interfaces.SavingsGoalRepository
    savingsContributionRepo interfaces.SavingsContributionRepository
    roundupRepo         interfaces.AutoRoundupRepository
    walletRepo          interfaces.WalletRepository
    transactionRepo     interfaces.TransactionRepository
    userRepo            interfaces.UserRepository
    db                  *gorm.DB
}

func NewSavingsService(
    savingsGoalRepo interfaces.SavingsGoalRepository,
    savingsContributionRepo interfaces.SavingsContributionRepository,
    roundupRepo interfaces.AutoRoundupRepository,
    walletRepo interfaces.WalletRepository,
    transactionRepo interfaces.TransactionRepository,
    userRepo interfaces.UserRepository,
    db *gorm.DB,
) *SavingsService {
    return &SavingsService{
        savingsGoalRepo:     savingsGoalRepo,
        savingsContributionRepo: savingsContributionRepo,
        roundupRepo:         roundupRepo,
        walletRepo:          walletRepo,
        transactionRepo:    transactionRepo,
        userRepo:           userRepo,
        db:                 db,
    }
}

func (s *SavingsService) CreateGoal(ctx context.Context, userID uuid.UUID, req *request.CreateSavingsGoalRequest) (*response.SavingsGoalResponse, error) {
    startDate := time.Now()
    targetDate := startDate.AddDate(0, 0, req.DurationDays)
    
    goal := &entities.SavingsGoal{
        ID:              uuid.New(),
        UserID:          userID,
        Name:            req.Name,
        TargetAmount:    entities.AmountInKobo(req.TargetAmount),
        CurrentAmount:   0,
        InterestRate:    0,
        DurationDays:    req.DurationDays,
        StartDate:       startDate,
        TargetDate:      targetDate,
        IsAutoDebit:     req.IsAutoDebit,
        AutoDebitAmount: entities.AmountInKobo(req.AutoDebitAmount),
        AutoDebitDay:    req.AutoDebitDay,
        Status:          "active",
        CreatedAt:       startDate,
    }
    
    if err := s.savingsGoalRepo.Create(ctx, goal); err != nil {
        return nil, err
    }
    
    return s.mapGoalToResponse(goal), nil
}

func (s *SavingsService) ContributeToGoal(ctx context.Context, userID uuid.UUID, req *request.ContributeToGoalRequest) (*response.SavingsContributionResponse, error) {
    goalID, err := uuid.Parse(req.GoalID)
    if err != nil {
        return nil, err
    }
    
    goal, err := s.savingsGoalRepo.GetByID(ctx, goalID)
    if err != nil {
        return nil, err
    }
    
    if goal == nil {
        return nil, errors.New("goal not found")
    }
    
    if goal.UserID != userID {
        return nil, errors.New("unauthorized to contribute to this goal")
    }
    
    if goal.Status != "active" {
        return nil, errors.New("goal is not active")
    }
    
    if goal.CurrentAmount >= goal.TargetAmount {
        return nil, errors.New("goal already completed")
    }
    
    // Get wallet
    wallet, err := s.walletRepo.GetByUserIDForUpdate(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    if wallet.Balance < entities.AmountInKobo(req.Amount) {
        return nil, errors.New("insufficient balance")
    }
    
    reference := s.generateReference("SAV")
    
    var contribution *entities.SavingsContribution
    
    err = s.db.Transaction(func(tx *gorm.DB) error {
        // Debit wallet
        if err := s.walletRepo.Debit(ctx, wallet.ID, req.Amount, reference); err != nil {
            return err
        }
        
        // Create transaction
        transaction := &entities.Transaction{
            ID:            uuid.New(),
            Reference:     reference,
            WalletID:      wallet.ID,
            UserID:        userID,
            Type:          "debit",
            Category:      "savings",
            Amount:        entities.AmountInKobo(req.Amount),
            Fee:           0,
            VAT:           0,
            TotalAmount:   entities.AmountInKobo(req.Amount),
            BalanceBefore: wallet.Balance,
            BalanceAfter:  entities.AmountInKobo(wallet.Balance) - entities.AmountInKobo(req.Amount),
            Status:        "success",
            Description:   fmt.Sprintf("Contribution to savings goal: %s", goal.Name),
            CreatedAt:     time.Now(),
            CompletedAt:   timeNow(),
        }
        
        if err := s.transactionRepo.Create(ctx, transaction); err != nil {
            return err
        }
        
        // Create contribution
        contribution = &entities.SavingsContribution{
            ID:              uuid.New(),
            SavingsGoalID:   goalID,
            TransactionID:   transaction.ID,
            Amount:          entities.AmountInKobo(req.Amount),
            InterestEarned:  0,
            ContributionDate: time.Now(),
            IsAutoDebit:     false,
        }
        
        if err := s.savingsContributionRepo.Create(ctx, contribution); err != nil {
            return err
        }
        
        // Update goal current amount
        newAmount := entities.AmountInKobo(goal.CurrentAmount) + entities.AmountInKobo(req.Amount)
        if err := s.savingsGoalRepo.UpdateCurrentAmount(ctx, goalID, int64(newAmount)); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return &response.SavingsContributionResponse{
        ID:              contribution.ID,
        Amount:          int64(contribution.Amount),
        AmountNaira:     float64(contribution.Amount) / 100,
        InterestEarned:  int64(contribution.InterestEarned),
        ContributionDate: contribution.ContributionDate,
        IsAutoDebit:     contribution.IsAutoDebit,
    }, nil
}

func (s *SavingsService) GetGoals(ctx context.Context, userID uuid.UUID) ([]response.SavingsGoalResponse, error) {
    goals, err := s.savingsGoalRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    responses := make([]response.SavingsGoalResponse, len(goals))
    for i, goal := range goals {
        responses[i] = *s.mapGoalToResponse(&goal)
    }
    
    return responses, nil
}

func (s *SavingsService) GetGoal(ctx context.Context, userID uuid.UUID, goalID string) (*response.SavingsGoalResponse, error) {
    id, err := uuid.Parse(goalID)
    if err != nil {
        return nil, err
    }
    
    goal, err := s.savingsGoalRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if goal == nil {
        return nil, errors.New("goal not found")
    }
    
    if goal.UserID != userID {
        return nil, errors.New("unauthorized to view this goal")
    }
    
    return s.mapGoalToResponse(goal), nil
}

func (s *SavingsService) UpdateGoal(ctx context.Context, userID uuid.UUID, goalID string, req *request.UpdateSavingsGoalRequest) error {
    id, err := uuid.Parse(goalID)
    if err != nil {
        return err
    }
    
    goal, err := s.savingsGoalRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    if goal == nil {
        return errors.New("goal not found")
    }
    
    if goal.UserID != userID {
        return errors.New("unauthorized to update this goal")
    }
    
    if req.Name != "" {
        goal.Name = req.Name
    }
    
    if req.IsAutoDebit != nil {
        goal.IsAutoDebit = *req.IsAutoDebit
        goal.AutoDebitAmount = entities.AmountInKobo(req.AutoDebitAmount)
        goal.AutoDebitDay = req.AutoDebitDay
    }
    
    return s.savingsGoalRepo.Update(ctx, goal)
}

func (s *SavingsService) DeleteGoal(ctx context.Context, userID uuid.UUID, goalID string) error {
    id, err := uuid.Parse(goalID)
    if err != nil {
        return err
    }
    
    goal, err := s.savingsGoalRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    if goal == nil {
        return errors.New("goal not found")
    }
    
    if goal.UserID != userID {
        return errors.New("unauthorized to delete this goal")
    }
    
    // If there's money in the goal, withdraw it first
    if goal.CurrentAmount > 0 {
        // Withdraw funds back to wallet
        // ... implementation
    }
    
    return s.savingsGoalRepo.Delete(ctx, id)
}

func (s *SavingsService) ActivateRoundup(ctx context.Context, userID uuid.UUID, req *request.ActivateRoundupRequest) error {
    goalID, err := uuid.Parse(req.GoalID)
    if err != nil {
        return err
    }
    
    goal, err := s.savingsGoalRepo.GetByID(ctx, goalID)
    if err != nil {
        return err
    }
    
    if goal == nil {
        return errors.New("goal not found")
    }
    
    if goal.UserID != userID {
        return errors.New("unauthorized to link this goal")
    }
    
    // Check if roundup already exists
    existing, err := s.roundupRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if existing != nil {
        // Update existing
        existing.SavingsGoalID = goalID
        existing.IsActive = true
        existing.Multiplier = req.Multiplier
        existing.MaxDailyAmount = entities.AmountInKobo(req.MaxDailyAmount)
        return s.roundupRepo.Update(ctx, existing)
    }
    
    // Create new roundup
    roundup := &entities.AutoRoundup{
        ID:            uuid.New(),
        UserID:        userID,
        SavingsGoalID: goalID,
        IsActive:      true,
        Multiplier:    req.Multiplier,
        MaxDailyAmount: entities.AmountInKobo(req.MaxDailyAmount),
        TotalRoundup:  0,
    }
    
    return s.roundupRepo.Create(ctx, roundup)
}

func (s *SavingsService) DeactivateRoundup(ctx context.Context, userID uuid.UUID) error {
    return s.roundupRepo.Deactivate(ctx, userID)
}

func (s *SavingsService) GetRoundupStatus(ctx context.Context, userID uuid.UUID) (*response.RoundupStatusResponse, error) {
    roundup, err := s.roundupRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if roundup == nil {
        return &response.RoundupStatusResponse{
            IsActive: false,
        }, nil
    }
    
    return &response.RoundupStatusResponse{
        IsActive:        roundup.IsActive,
        SavingsGoalID:   roundup.SavingsGoalID.String(),
        Multiplier:      roundup.Multiplier,
        MaxDailyAmount:  int64(roundup.MaxDailyAmount),
        TotalRoundup:    int64(roundup.TotalRoundup),
        TotalRoundupNaira: float64(roundup.TotalRoundup) / 100,
    }, nil
}

func (s *SavingsService) mapGoalToResponse(goal *entities.SavingsGoal) *response.SavingsGoalResponse {
    progressPercent := float64(goal.CurrentAmount) / float64(goal.TargetAmount) * 100
    daysRemaining := int(goal.TargetDate.Sub(time.Now()).Hours() / 24)
    if daysRemaining < 0 {
        daysRemaining = 0
    }
    
    return &response.SavingsGoalResponse{
        ID:                 goal.ID,
        Name:               goal.Name,
        TargetAmount:       int64(goal.TargetAmount),
        TargetAmountNaira:  float64(goal.TargetAmount) / 100,
        CurrentAmount:      int64(goal.CurrentAmount),
        CurrentAmountNaira: float64(goal.CurrentAmount) / 100,
        ProgressPercent:    progressPercent,
        InterestRate:       goal.InterestRate,
        DurationDays:       goal.DurationDays,
        StartDate:          goal.StartDate,
        TargetDate:         goal.TargetDate,
        DaysRemaining:      daysRemaining,
        IsAutoDebit:        goal.IsAutoDebit,
        AutoDebitAmount:    int64(goal.AutoDebitAmount),
        Status:             goal.Status,
        CreatedAt:          goal.CreatedAt,
    }
}

func (s *SavingsService) generateReference(prefix string) string {
    return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// func timeNow() *time.Time {
//     now := time.Now()
//     return &now
// }