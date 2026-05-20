package request

type CreateSavingsGoalRequest struct {
    Name            string `json:"name" validate:"required,max=100"`
    TargetAmount    int64  `json:"target_amount" validate:"required,min=100000"` // In kobo (min 1000 NGN)
    DurationDays    int    `json:"duration_days" validate:"required,min=7,max=3650"` // 7 days to 10 years
    IsAutoDebit     bool   `json:"is_auto_debit"`
    AutoDebitAmount int64  `json:"auto_debit_amount" validate:"required_if=IsAutoDebit true,min=5000"`
    AutoDebitDay    int    `json:"auto_debit_day" validate:"required_if=IsAutoDebit true,min=1,max=28"`
}

type ContributeToGoalRequest struct {
    GoalID   string `json:"goal_id" validate:"required"`
    Amount   int64  `json:"amount" validate:"required,min=5000"` // In kobo
    SourceWalletID string `json:"source_wallet_id" validate:"required"`
}

type UpdateSavingsGoalRequest struct {
    Name            string `json:"name" validate:"omitempty,max=100"`
    IsAutoDebit     *bool  `json:"is_auto_debit"`
    AutoDebitAmount int64  `json:"auto_debit_amount" validate:"required_if=IsAutoDebit true,min=5000"`
    AutoDebitDay    int    `json:"auto_debit_day" validate:"required_if=IsAutoDebit true,min=1,max=28"`
}

type ActivateRoundupRequest struct {
    GoalID          string `json:"goal_id" validate:"required"`
    Multiplier      int    `json:"multiplier" validate:"required,min=1,max=10"` // Round up to nearest X naira
    MaxDailyAmount  int64  `json:"max_daily_amount" validate:"required,min=5000"` // In kobo
}