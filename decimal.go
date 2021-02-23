package big

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

var (
	flZero = *big.NewFloat(0)

	ZERO = NewDecimal(0)
	ONE  = NewDecimal(1)
	TEN  = NewDecimal(10)

	MarshalQuoted = false
)

// Decimal is the main exported type. It is a simple, immutable wrapper around a *big.Float
type Decimal struct {
	fl *big.Float
}

// NewDecimal creates a new Decimal type from a float value.
func NewDecimal(fl float64) Decimal {
	return Decimal{big.NewFloat(fl)}
}

// NewFromString creates a new Decimal type from a string value.
func NewFromString(str string) Decimal {
	bfl := new(big.Float)
	bfl.UnmarshalText([]byte(str))
	return Decimal{bfl}
}

// NewFromInt creates a new Decimal type from an int value
func NewFromInt(dec int) Decimal {
	return Decimal{big.NewFloat(float64(dec))}
}

// MaxSlice returns the max of a slice of decimals
func MaxSlice(decimals ...Decimal) Decimal {
	initial := NewFromString("-Inf")

	for _, decimal := range decimals {
		if decimal.GT(initial) {
			initial = decimal
		}
	}

	return initial
}

// MinSlice returns the min of a slice of decimals
func MinSlice(decimals ...Decimal) Decimal {
	initial := NewFromString("Inf")
	for _, decimal := range decimals {
		if decimal.LT(initial) {
			initial = decimal
		}
	}

	return initial
}

// Add adds a decimal instance to another Decimal instance.
func (d Decimal) Add(addend Decimal) Decimal {
	return Decimal{d.cpy().Add(d.fl, addend.fl)}
}

// Sub subtracts another decimal instance from this Decimal instance.
func (d Decimal) Sub(subtrahend Decimal) Decimal {
	return Decimal{d.cpy().Sub(d.fl, subtrahend.fl)}
}

// Mul multiplies another decimal instance with this Decimal instance.
func (d Decimal) Mul(factor Decimal) Decimal {
	return Decimal{d.cpy().Mul(d.fl, factor.fl)}
}

// Div divides this Decimal by the denominator passed.
func (d Decimal) Div(denominator Decimal) Decimal {
	return Decimal{d.cpy().Quo(d.fl, denominator.fl)}
}

// Frac returns another Decimal instance representing this Decimal multiplied by the
// provided float.
func (d Decimal) Frac(f float64) Decimal {
	return d.Mul(NewDecimal(f))
}

// Neg returns this Decimal multiplied by -1.
func (d Decimal) Neg() Decimal {
	return d.Mul(NewDecimal(-1))
}

// Abs returns the absolute value of this Decimal
func (d Decimal) Abs() Decimal {
	if d.LT(ZERO) {
		return d.Mul(ONE.Neg())
	}

	return d
}

// Pow returns the decimal to the inputted power
func (d Decimal) Pow(exp int) Decimal {
	if exp == 0 {
		return ONE
	}

	x := Decimal{d.cpy()}

	for i := 1; i < exp; i++ {
		x = x.Mul(d)
	}

	return x
}

// Sqrt returns the decimal's square root
func (d Decimal) Sqrt() Decimal {
	return Decimal{d.cpy().Sqrt(d.cpy())}
}

// EQ returns true if this Decimal exactly equals the provided decimal.
func (d Decimal) EQ(other Decimal) bool {
	return d.Cmp(other) == 0
}

// LT returns true if this decimal is less than the provided decimal.
func (d Decimal) LT(other Decimal) bool {
	return d.Cmp(other) < 0
}

// LTE returns true if this decimal is less or equal to the provided decimal.
func (d Decimal) LTE(other Decimal) bool {
	return d.Cmp(other) <= 0
}

// GT returns true if this decimal is greater than the provided decimal.
func (d Decimal) GT(other Decimal) bool {
	return d.Cmp(other) > 0
}

// GTE returns true if this decimal is greater than or equal to the provided decimal.
func (d Decimal) GTE(other Decimal) bool {
	return d.Cmp(other) >= 0
}

// Cmp will return 1 if this decimal is greater than the provided, 0 if they are the same, and -1 if it is less.
func (d Decimal) Cmp(other Decimal) int {
	return d.fl.Cmp(other.fl)
}

// Float will return this Decimal as a float value.
// Note that there may be some loss of precision in this operation.
func (d Decimal) Float() float64 {
	f, _ := d.fl.Float64()
	return f
}

// Zero will return true if this Decimal is equal to 0.
// Note: Zero is deprecated. Use IsZero instead
func (d Decimal) Zero() bool {
	return d.IsZero()
}

// IsZero will return true if this Decimal is equal to 0.
func (d Decimal) IsZero() bool {
	return d.fl == nil || d.fl.Cmp(&flZero) == 0
}

func (d Decimal) String() string {
	if d.fl == nil {
		d.fl = new(big.Float)
	}

	return d.fl.String()
}

func (d Decimal) FormattedString(places int) string {
	format := "%." + fmt.Sprint(places) + "f"
	fl := d.Float()
	return fmt.Sprintf(format, fl)
}

// MarshalJSON implements the json.Marshaler interface
func (d Decimal) MarshalJSON() ([]byte, error) {
	if MarshalQuoted {
		return []byte("\"" + d.String() + "\""), nil
	}

	return d.fl.MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Decimal) UnmarshalJSON(b []byte) error {
	if d.fl == nil {
		d.fl = big.NewFloat(0)
	}

	if isQuoted(b) {
		b = b[1 : len(b)-1]
	}

	return d.fl.UnmarshalText(b)
}

func isQuoted(b []byte) bool {
	quoteByte := byte('"')
	return len(b) > 0 && b[0] == quoteByte && b[len(b)-1] == quoteByte
}

// Value implements the sql.Valuer interface
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface
func (d *Decimal) Scan(src interface{}) error {
	switch src.(type) {
	case string:
		return json.Unmarshal([]byte(src.(string)), d)
	case []byte:
		return json.Unmarshal([]byte(src.([]byte)), d)
	default:
		return errors.New(fmt.Sprint("Passed value ", src, " should be a string"))
	}
}

func (d Decimal) cpy() *big.Float {
	cpy := new(big.Float)
	return cpy.Copy(d.fl)
}
