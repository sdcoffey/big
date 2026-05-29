package big

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"unicode"
)

const minPrecision uint = 256

var (
	flZero = *big.NewFloat(0)

	// NaN == Not a Number
	NaN = Decimal{nan: true}

	// ZERO == 0
	ZERO = NewFromString("0")

	// ONE == 1
	ONE = NewFromString("1")

	// TEN == 10
	TEN = NewFromString("10")

	// MarshalQuoted - can toggle this to true to marshal values as strings
	MarshalQuoted = false
)

// Decimal is the main exported type. It is a simple, immutable wrapper around a *big.Float
type Decimal struct {
	fl  *big.Float
	nan bool
}

// NewDecimal creates a new Decimal type from a float value.
func NewDecimal(val float64) Decimal {
	if math.IsNaN(val) {
		return NaN
	}

	fl := newFloat(53)
	fl.SetFloat64(val)

	return Decimal{
		fl: fl,
	}
}

// NewFromString creates a new Decimal type from a string value.
func NewFromString(str string) Decimal {
	bfl := newFloat(decimalPrecision(str))

	if _, _, err := bfl.Parse(str, 10); err != nil {
		return NaN
	}

	return Decimal{fl: bfl}
}

// NewFromInt creates a new Decimal type from an int value
func NewFromInt(dec int) Decimal {
	fl := newFloat(intPrecision(dec))
	fl.SetInt64(int64(dec))
	return Decimal{fl: fl}
}

func newFloat(precision uint) *big.Float {
	if precision < minPrecision {
		precision = minPrecision
	}

	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven)
}

func zeroDecimal() Decimal {
	return NewFromInt(0)
}

func oneDecimal() Decimal {
	return NewFromInt(1)
}

func decimalPrecision(str string) uint {
	digits := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			digits++
		}
	}

	if digits == 0 {
		return minPrecision
	}

	return uint(digits)*4 + 16
}

func intPrecision(dec int) uint {
	if dec == 0 {
		return minPrecision
	}

	if dec < 0 {
		magnitude := uint64(-(dec + 1)) + 1
		return uint(bits.Len64(magnitude)) + 1
	}

	return uint(bits.Len64(uint64(dec))) + 1
}

// MaxSlice returns the max of a slice of decimals
func MaxSlice(decimals ...Decimal) Decimal {
	if anyNan(decimals...) {
		return NaN
	} else if len(decimals) == 0 {
		return zeroDecimal()
	}

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
	if anyNan(decimals...) {
		return NaN
	} else if len(decimals) == 0 {
		return zeroDecimal()
	}

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
	return nanGuard(func() Decimal {
		return Decimal{fl: resultFloat(1, d, addend).Add(d.value(), addend.value())}
	}, d, addend)
}

// Sub subtracts another decimal instance from this Decimal instance.
func (d Decimal) Sub(subtrahend Decimal) Decimal {
	return nanGuard(func() Decimal {
		return Decimal{fl: resultFloat(1, d, subtrahend).Sub(d.value(), subtrahend.value())}
	}, d, subtrahend)
}

// Mul multiplies another decimal instance with this Decimal instance.
func (d Decimal) Mul(factor Decimal) Decimal {
	return nanGuard(func() Decimal {
		return Decimal{fl: newFloat(sumPrecision(d, factor)).Mul(d.value(), factor.value())}
	}, d, factor)
}

// Div divides this Decimal by the denominator passed.
func (d Decimal) Div(denominator Decimal) Decimal {
	return nanGuard(func() Decimal {
		return Decimal{fl: resultFloat(0, d, denominator).Quo(d.value(), denominator.value())}
	}, d, denominator)
}

// Frac returns another Decimal instance representing this Decimal multiplied by the
// provided float.
func (d Decimal) Frac(f float64) Decimal {
	fractionFactor := NewDecimal(f)

	return nanGuard(func() Decimal {
		return d.Mul(NewDecimal(f))
	}, d, fractionFactor)
}

// Neg returns this Decimal multiplied by -1.
func (d Decimal) Neg() Decimal {
	return nanGuard(func() Decimal {
		return d.Mul(NewDecimal(-1))
	}, d)
}

// Abs returns the absolute value of this Decimal
func (d Decimal) Abs() Decimal {
	if d.LT(zeroDecimal()) {
		return d.Neg()
	}

	return d
}

// Pow returns the decimal to the inputted power
func (d Decimal) Pow(exp int) Decimal {
	return nanGuard(func() Decimal {
		if exp == 0 {
			return oneDecimal()
		}

		if exp < 0 {
			if d.IsZero() {
				return NaN
			}

			return oneDecimal().Div(d.Pow(-exp))
		}

		x := oneDecimal()
		base := Decimal{fl: d.cpy()}

		for exp > 0 {
			if exp%2 == 1 {
				x = x.Mul(base)
			}

			exp /= 2
			if exp > 0 {
				base = base.Mul(base)
			}
		}

		return x
	}, d)
}

// Sqrt returns the decimal's square root
func (d Decimal) Sqrt() Decimal {
	return nanGuard(func() Decimal {
		if d.LT(zeroDecimal()) {
			return NaN
		}

		return Decimal{fl: d.cpy().Sqrt(d.cpy())}
	}, d)
}

// EQ returns true if this Decimal exactly equals the provided decimal.
func (d Decimal) EQ(other Decimal) bool {
	if anyNan(d, other) {
		return false
	}

	return d.Cmp(other) == 0
}

// LT returns true if this decimal is less than the provided decimal.
func (d Decimal) LT(other Decimal) bool {
	if anyNan(d, other) {
		return false
	}

	return d.Cmp(other) < 0
}

// LTE returns true if this decimal is less or equal to the provided decimal.
func (d Decimal) LTE(other Decimal) bool {
	if anyNan(d, other) {
		return false
	}

	return d.Cmp(other) <= 0
}

// GT returns true if this decimal is greater than the provided decimal.
func (d Decimal) GT(other Decimal) bool {
	if anyNan(d, other) {
		return false
	}

	return d.Cmp(other) > 0
}

// GTE returns true if this decimal is greater than or equal to the provided decimal.
func (d Decimal) GTE(other Decimal) bool {
	if anyNan(d, other) {
		return false
	}

	return d.Cmp(other) >= 0
}

// Cmp will return 1 if this decimal is greater than the provided, 0 if they are the same, and -1 if it is less.
func (d Decimal) Cmp(other Decimal) int {
	if d.NaN() && other.NaN() {
		return 0
	}

	if d.NaN() {
		return -1
	}

	if other.NaN() {
		return 1
	}

	return d.value().Cmp(other.value())
}

// Float will return this Decimal as a float value.
// Note that there may be some loss of precision in this operation.
func (d Decimal) Float() float64 {
	if d.NaN() {
		return math.NaN()
	}

	f, _ := d.value().Float64()
	return f
}

// Zero will return true if this Decimal is equal to 0.
// Deprecated: Use IsZero instead
func (d Decimal) Zero() bool {
	return d.IsZero()
}

// NaN returns true if the underlying is not a valid number
func (d Decimal) NaN() bool {
	return d.nan
}

// IsZero will return true if this Decimal is equal to 0.
func (d Decimal) IsZero() bool {
	if d.NaN() {
		return false
	}

	return d.value().Cmp(&flZero) == 0
}

func (d Decimal) String() string {
	if d.NaN() {
		return "NaN"
	}

	return d.value().String()
}

// FormattedString returns the string value of the number to the requested precision
func (d Decimal) FormattedString(places int) string {
	if d.NaN() {
		return d.String()
	}

	return d.value().Text('f', places)
}

// MarshalJSON implements the json.Marshaler interface
func (d Decimal) MarshalJSON() ([]byte, error) {
	if MarshalQuoted {
		return []byte("\"" + d.String() + "\""), nil
	}

	if d.NaN() {
		return []byte("null"), nil
	}

	return d.value().MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Decimal) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)

	if isQuoted(b) {
		b = b[1 : len(b)-1]
	}

	if bytes.Equal(b, []byte("null")) || bytes.Equal(b, []byte("NaN")) {
		*d = NaN
		return nil
	}

	fl := newFloat(decimalPrecision(string(b)))
	if _, _, err := fl.Parse(string(b), 10); err != nil {
		return err
	}

	*d = Decimal{fl: fl}
	return nil
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
	switch src := src.(type) {
	case string:
		return d.UnmarshalJSON([]byte(src))
	case []byte:
		return d.UnmarshalJSON(src)
	case nil:
		*d = NaN
		return nil
	default:
		return errors.New(fmt.Sprint("Passed value ", src, " should be a string"))
	}
}

func (d Decimal) cpy() *big.Float {
	val := d.value()
	cpy := newFloat(val.Prec())
	return cpy.Copy(val)
}

func resultFloat(extra uint, decimals ...Decimal) *big.Float {
	return newFloat(maxPrecision(decimals...) + extra)
}

func maxPrecision(decimals ...Decimal) uint {
	precision := minPrecision
	for _, decimal := range decimals {
		if decimal.NaN() {
			continue
		}

		if valuePrecision := decimal.value().Prec(); valuePrecision > precision {
			precision = valuePrecision
		}
	}

	return precision
}

func sumPrecision(decimals ...Decimal) uint {
	precision := uint(0)
	for _, decimal := range decimals {
		if decimal.NaN() {
			continue
		}

		precision += decimal.value().Prec()
	}

	return precision
}

func (d Decimal) value() *big.Float {
	if d.fl != nil {
		return d.fl
	}

	return newFloat(minPrecision)
}

func anyNan(decimals ...Decimal) bool {
	for _, decimal := range decimals {
		if decimal.NaN() {
			return true
		}
	}

	return false
}

func nanGuard(yeildFunc func() Decimal, decimals ...Decimal) Decimal {
	if anyNan(decimals...) {
		return NaN
	}

	return yeildFunc()
}
