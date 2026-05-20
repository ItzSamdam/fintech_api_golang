package response

import (
    "time"
    "github.com/google/uuid"
)

type SavingsGoalResponse struct {
    ID              uuid.UUID `json:"id"`
    Name            string    `json:"name"`
    TargetAmount    int64     `json:"target_amount"`     // In kobo
    TargetAmountNaira float64 `json:"target_amount_naira"`
    CurrentAmount   int64     `json:"current_amount"`    // In kobo
    CurrentAmountNaira float64 `json:"current_amount_naira"`
    ProgressPercent float64   `json:"progress_percent"`
    InterestRate    float64   `json:"interest_rate"`
    DurationDays    int       `json:"duration_days"`
    StartDate       time.Time `json:"start_date"`
    TargetDate      time.Time `json:"target_date"`
    DaysRemaining   int       `json:"days_remaining"`
    IsAutoDebit     bool      `json:"is_auto_debit"`
    AutoDebitAmount int64     `json:"auto_debit_amount,omitempty"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
}

type SavingsContributionResponse struct {
    ID              uuid.UUID `json:"id"`
    Amount          int64     `json:"amount"`           // In kobo
    AmountNaira     float64   `json:"amount_naira"`
    InterestEarned  int64     `json:"interest_earned"`  // In kobo
    ContributionDate time.Time `json:"contribution_date"`
    IsAutoDebit     bool      `json:"is_auto_debit"`
}

type RoundupStatusResponse struct {
    IsActive        bool   `json:"is_active"`
    SavingsGoalID   string `json:"savings_goal_id"`
    Multiplier      int    `json:"multiplier"`
    MaxDailyAmount  int64  `json:"max_daily_amount"`
    TotalRoundup    int64  `json:"total_roundup"`
    TotalRoundupNaira float64 `json:"total_roundup_naira"`
}