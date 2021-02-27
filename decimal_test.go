package big

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type equalExample struct {
	value    Decimal
	expected string
}

type booleanExample struct {
	value    bool
	expected bool
}

func validateEqExamples(t *testing.T, examples ...equalExample) {
	for _, ex := range examples {
		assert.EqualValues(t, ex.expected, ex.value.String())
	}
}

func validateBoolExamples(t *testing.T, examples ...booleanExample) {
	for _, ex := range examples {
		if ex.expected {
			assert.True(t, ex.value)
		} else {
			assert.False(t, ex.value)
		}
	}
}

func TestNewDecimal(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		d := NewDecimal(math.Pi)

		assert.EqualValues(t, math.Pi, d.Float())
	})

	t.Run("NaN", func(t *testing.T) {
		d := NewDecimal(math.NaN())

		assert.Nil(t, d.fl)
	})
}

func TestNewFromString(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewFromString("1.87"),
			expected: "1.87",
		},
		equalExample{
			value:    NewFromString("NaN"),
			expected: "NaN",
		},
	)
}

func TestNewFromInt(t *testing.T) {
	d := NewFromInt(1)

	assert.EqualValues(t, "1", d.String())
}

func TestMaxSlice(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value: MaxSlice(NewDecimal(-100),
				NewDecimal(100),
				NewDecimal(0)),
			expected: "100",
		},
		equalExample{
			value:    MaxSlice(NewDecimal(100), NaN),
			expected: "NaN",
		},
		equalExample{
			value:    MaxSlice(),
			expected: "0",
		},
	)
}

func TestMinSlice(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value: MinSlice(NewDecimal(-100),
				NewDecimal(100),
				NewDecimal(0)),
			expected: "-100",
		},
		equalExample{
			value:    MinSlice(NewDecimal(100), NaN),
			expected: "NaN",
		},
		equalExample{
			value:    MinSlice(),
			expected: "0",
		},
	)
}

func TestDecimal_Add(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(3.14).Add(NewDecimal(2)),
			expected: "5.14",
		},
		equalExample{
			value:    NaN.Add(NewDecimal(2)),
			expected: "NaN",
		},
		equalExample{
			value:    NewDecimal(1).Add(NaN),
			expected: "NaN",
		},
	)
}

func TestDecimal_Sub(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(3.14).Sub(NewDecimal(2)),
			expected: "1.14",
		},
		equalExample{
			value:    NaN.Sub(NewDecimal(1)),
			expected: "NaN",
		},
		equalExample{
			value:    ONE.Sub(NaN),
			expected: "NaN",
		},
	)
}

func TestDecimal_Mul(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(3.14).Mul(TEN),
			expected: "31.4",
		},
		equalExample{
			value:    NaN.Mul(TEN),
			expected: "NaN",
		},
		equalExample{
			value:    TEN.Mul(NaN),
			expected: "NaN",
		},
	)
}

func TestDecimal_Div(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(3.14).Div(TEN),
			expected: "0.314",
		},
		equalExample{
			value:    TEN.Div(NaN),
			expected: "NaN",
		},
		equalExample{
			value:    NaN.Div(TEN),
			expected: "NaN",
		},
	)
}

func TestDecimal_Neg(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    TEN.Neg(),
			expected: "-10",
		},
		equalExample{
			value:    NaN.Neg(),
			expected: "NaN",
		},
	)
}

func TestDecimal_Abs(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    TEN.Abs(),
			expected: "10",
		},
		equalExample{
			value:    NewFromString("-10").Abs(),
			expected: "10",
		},
		equalExample{
			value:    NaN.Abs(),
			expected: "NaN",
		},
	)
}

func TestDecimal_Frac(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    TEN.Frac(0.5),
			expected: "5",
		},
		equalExample{
			value:    TEN.Frac(1.5),
			expected: "15",
		},
		equalExample{
			value:    NaN.Frac(0.5),
			expected: "NaN",
		},
	)
}

func TestDecimal_EQ(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    NewDecimal(182.1921).EQ(NewDecimal(182.1921)),
			expected: true,
		},
		booleanExample{
			value:    NaN.EQ(NaN),
			expected: false,
		},
	)
}

func TestDecimal_GT(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    NewDecimal(182.1921).GT(NewDecimal(182.1921)),
			expected: false,
		},
		booleanExample{
			value:    NewDecimal(182.1920).GT(NewDecimal(182.1921)),
			expected: false,
		},
		booleanExample{
			value:    NewDecimal(182.1921).GT(NewDecimal(182.1920)),
			expected: true,
		},
		booleanExample{
			value:    NaN.GT(NaN),
			expected: false,
		},
	)
}

func TestDecimal_GTE(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    NewDecimal(182.1921).GTE(NewDecimal(182.1921)),
			expected: true,
		},
		booleanExample{
			value:    NewDecimal(182.1920).GTE(NewDecimal(182.1921)),
			expected: false,
		},
		booleanExample{
			value:    NewDecimal(182.1921).GTE(NewDecimal(182.1920)),
			expected: true,
		},
		booleanExample{
			value:    NaN.GTE(NaN),
			expected: false,
		},
	)
}

func TestDecimal_LT(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    NewDecimal(182.1921).LT(NewDecimal(182.1921)),
			expected: false,
		},
		booleanExample{
			value:    NewDecimal(182.1920).LT(NewDecimal(182.1921)),
			expected: true,
		},
		booleanExample{
			value:    NewDecimal(182.1921).LT(NewDecimal(182.1920)),
			expected: false,
		},
		booleanExample{
			value:    NaN.LT(NaN),
			expected: false,
		},
	)
}

func TestDecimal_LTE(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    NewDecimal(182.1921).LTE(NewDecimal(182.1921)),
			expected: true,
		},
		booleanExample{
			value:    NewDecimal(182.1920).LTE(NewDecimal(182.1921)),
			expected: true,
		},
		booleanExample{
			value:    NewDecimal(182.1921).LTE(NewDecimal(182.1920)),
			expected: false,
		},
		booleanExample{
			value:    NaN.LTE(NaN),
			expected: false,
		},
	)
}

func TestDecimal_Cmp(t *testing.T) {
	assert.EqualValues(t, 0, ONE.Cmp(ONE))
	assert.EqualValues(t, 1, TEN.Cmp(ONE))
	assert.EqualValues(t, -1, ONE.Cmp(TEN))
	assert.EqualValues(t, 0, NaN.Cmp(NaN))
}

func TestDecimal_Float(t *testing.T) {
	assert.EqualValues(t, 1.13, NewDecimal(1.13).Float())
	assert.True(t, math.IsNaN(NaN.Float()))
}

func TestDecimal_Pow(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(8).Pow(2),
			expected: "64",
		},
		equalExample{
			value:    TEN.Pow(0),
			expected: "1",
		},
		equalExample{
			value:    NaN.Pow(2),
			expected: "NaN",
		},
	)
}

func TestDecimal_Sqrt(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(64).Sqrt(),
			expected: "8",
		},
		equalExample{
			value:    NaN.Sqrt(),
			expected: "NaN",
		},
	)
}

func TestDecimal_String(t *testing.T) {
	validateEqExamples(t,
		equalExample{
			value:    NewDecimal(1.13),
			expected: "1.13",
		},
		equalExample{
			value:    NaN,
			expected: "NaN",
		},
	)
}

func TestDecimal_FormattedString(t *testing.T) {
	assert.EqualValues(t, "3.1416", NewDecimal(math.Pi).FormattedString(4))
	assert.EqualValues(t, "NaN", NaN.FormattedString(4))
}

func TestDecimal_IsZero(t *testing.T) {
	validateBoolExamples(t,
		booleanExample{
			value:    ZERO.IsZero(),
			expected: true,
		},
		booleanExample{
			value:    ONE.IsZero(),
			expected: false,
		},
		booleanExample{
			value:    NaN.IsZero(),
			expected: false,
		},
	)
}
func TestDecimal_Json(t *testing.T) {
	type jsonType struct {
		Decimal Decimal `json:"decimal"`
	}

	t.Run("MarshalJSON - quoted", func(t *testing.T) {
		MarshalQuoted = true
		tmpStruct := jsonType{
			Decimal: NewFromString("3.1419"),
		}
		marshaled, err := json.Marshal(tmpStruct)

		assert.NoError(t, err)
		assert.Equal(t, `{"decimal":"3.1419"}`, string(marshaled))
		MarshalQuoted = false
	})

	t.Run("MarshalJSON - unquoted", func(t *testing.T) {
		tmpStruct := jsonType{
			Decimal: NewFromString("3.1419"),
		}
		marshaled, err := json.Marshal(tmpStruct)

		assert.NoError(t, err)
		assert.Equal(t, `{"decimal":3.1419}`, string(marshaled))
	})

	t.Run("UnmarshalJSON - unquoted", func(t *testing.T) {
		var ts jsonType

		d := `{"decimal":3.1419}`
		err := json.Unmarshal([]byte(d), &ts)

		assert.NoError(t, err)
		assert.Equal(t, "3.1419", ts.Decimal.String())
	})

	t.Run("UnmarshalJSON - quoted", func(t *testing.T) {
		var ts jsonType

		d := `{"decimal":"3.1419"}`
		err := json.Unmarshal([]byte(d), &ts)

		assert.NoError(t, err)
		assert.Equal(t, "3.1419", ts.Decimal.String())
	})
}

func TestDecimal_Sql(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		d := ONE
		value, err := d.Value()

		assert.NoError(t, err)
		assert.Equal(t, `1`, value)
	})

	t.Run("Scan", func(t *testing.T) {
		var d Decimal

		data := `1.23`
		err := d.Scan(data)

		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())
	})

	t.Run("Scan []byte", func(t *testing.T) {
		var d Decimal

		data := `1.23`
		err := d.Scan([]byte(data))

		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())
	})

	t.Run("Scan returns error when src is not string", func(t *testing.T) {
		var d Decimal

		data := 1.23
		err := d.Scan(data)

		assert.NotNil(t, err)
		assert.EqualValues(t, "Passed value 1.23 should be a string", err.Error())
	})
}
