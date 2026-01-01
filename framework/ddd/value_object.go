package ddd

// ValueObject is the marker interface for value objects.
// Value objects are immutable and are compared by their attributes rather than identity.
type ValueObject interface {
	// Equals checks if this value object is equal to another.
	Equals(other ValueObject) bool
}

// BaseValueObject provides common functionality for value objects.
// Embed this in your value objects and override Equals as needed.
type BaseValueObject struct{}

// Equals default implementation - always returns false.
// Concrete value objects should override this method.
func (b BaseValueObject) Equals(other ValueObject) bool {
	return false
}

// StringValue is a simple string-based value object.
type StringValue string

// String returns the string representation.
func (s StringValue) String() string {
	return string(s)
}

// Equals checks if two StringValues are equal.
func (s StringValue) Equals(other ValueObject) bool {
	if o, ok := other.(StringValue); ok {
		return s == o
	}
	return false
}

// Money is an example value object representing monetary values.
type Money struct {
	Amount   int64  // Amount in smallest currency unit (e.g., cents)
	Currency string // ISO 4217 currency code
}

// NewMoney creates a new Money value object.
func NewMoney(amount int64, currency string) Money {
	return Money{Amount: amount, Currency: currency}
}

// Equals checks if two Money values are equal.
func (m Money) Equals(other ValueObject) bool {
	if o, ok := other.(Money); ok {
		return m.Amount == o.Amount && m.Currency == o.Currency
	}
	return false
}

// Add adds two Money values. Panics if currencies don't match.
func (m Money) Add(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot add money with different currencies")
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}
}

// Subtract subtracts another Money value. Panics if currencies don't match.
func (m Money) Subtract(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot subtract money with different currencies")
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}
}

// IsZero checks if the amount is zero.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive checks if the amount is positive.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// IsNegative checks if the amount is negative.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}
