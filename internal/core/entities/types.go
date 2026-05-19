package entities

import (
    "database/sql/driver"
    "fmt"
)

// AmountInKobo represents monetary values in kobo (1 NGN = 100 kobo)
type AmountInKobo int64

// Scan implements sql.Scanner interface
func (a *AmountInKobo) Scan(value interface{}) error {
    if value == nil {
        *a = 0
        return nil
    }
    
    switch v := value.(type) {
    case int64:
        *a = AmountInKobo(v)
    case int32:
        *a = AmountInKobo(v)
    case int:
        *a = AmountInKobo(v)
    case float64:
        *a = AmountInKobo(v)
    default:
        return fmt.Errorf("cannot scan type %T into AmountInKobo", value)
    }
    return nil
}

// Value implements driver.Valuer interface
func (a AmountInKobo) Value() (driver.Value, error) {
    return int64(a), nil
}

// ToNaira converts kobo to naira (float64)
func (a AmountInKobo) ToNaira() float64 {
    return float64(a) / 100
}

// FromNaira converts naira to kobo
func FromNaira(amount float64) AmountInKobo {
    return AmountInKobo(amount * 100)
}