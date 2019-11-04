package big

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimal(t *testing.T) {
	t.Run("NewFromString", func(t *testing.T) {
		d := NewFromString("1.87")

		assert.EqualValues(t, "1.87", d.String())
	})

	t.Run("Add", func(t *testing.T) {
		f1 := NewDecimal(3.14)
		f2 := NewDecimal(2)

		assert.EqualValues(t, "5.14", f1.Add(f2).String())
	})

	t.Run("Sub", func(t *testing.T) {
		f1 := NewDecimal(3.14)
		f2 := NewDecimal(2)

		assert.EqualValues(t, "1.14", f1.Sub(f2).String())
	})

	t.Run("Mul", func(t *testing.T) {
		f1 := NewDecimal(3.14)
		f2 := NewDecimal(2)

		assert.EqualValues(t, "6.28", f1.Mul(f2).String())
		assert.EqualValues(t, "3.14", f1.String())
		assert.EqualValues(t, "2", f2.String())
	})

	t.Run("Div", func(t *testing.T) {
		f1 := NewDecimal(3.14)
		f2 := NewDecimal(2)

		assert.EqualValues(t, "1.57", f1.Div(f2).String())
	})

	t.Run("Neg", func(t *testing.T) {
		f1 := NewDecimal(3.14)

		assert.EqualValues(t, "-3.14", f1.Neg().String())
	})

	t.Run("Abs", func(t *testing.T) {
		f1 := NewDecimal(3.14)
		assert.EqualValues(t, 3.14, f1.Abs().Float())

		f2 := NewDecimal(-3.14)
		assert.EqualValues(t, 3.14, f2.Abs().Float())
	})

	t.Run("Frac", func(t *testing.T) {
		f1 := NewDecimal(3.14)

		assert.EqualValues(t, 1.57, f1.Frac(0.5).Float())
	})

	t.Run("EQ", func(t *testing.T) {
		f1 := NewDecimal(182.1921)
		f2 := NewDecimal(182.1921)

		assert.True(t, f2.EQ(f1))
	})

	t.Run("GT", func(t *testing.T) {
		f1 := NewDecimal(1.3419)
		f2 := NewDecimal(13419)

		assert.True(t, f2.GT(f1))
	})

	t.Run("GTE", func(t *testing.T) {
		f1 := NewDecimal(1.3419)
		f2 := NewDecimal(1.3419)

		assert.True(t, f2.GTE(f1))
	})

	t.Run("LT", func(t *testing.T) {
		f1 := NewDecimal(1.3419)
		f2 := NewDecimal(13419)

		assert.True(t, f1.LT(f2))
	})

	t.Run("LTE", func(t *testing.T) {
		f1 := NewDecimal(1.3419)
		f2 := NewDecimal(1.3419)

		assert.True(t, f1.LTE(f2))
	})

	t.Run("Cmp", func(t *testing.T) {
		f1 := NewDecimal(1.3419)
		f2 := NewDecimal(13419)

		assert.EqualValues(t, 1, f2.Cmp(f1))
	})

	t.Run("Float", func(t *testing.T) {
		f := NewDecimal(1.3419)
		assert.EqualValues(t, 1.3419, f.Float())
	})

	t.Run("Pow", func(t *testing.T) {
		t.Run("when exp is 0", func(t *testing.T) {
			f := NewDecimal(8)
			assert.EqualValues(t, 1, f.Pow(0).Float())
		})

		t.Run("when exp is positive", func(t *testing.T) {
			f := NewDecimal(8)
			assert.EqualValues(t, 512, f.Pow(3).Float())
		})
	})

	t.Run("Sqrt", func(t *testing.T) {
		f := NewDecimal(64)
		assert.EqualValues(t, 8, f.Sqrt().Float())
	})

	t.Run("String", func(t *testing.T) {
		f := NewDecimal(1.3419)
		assert.EqualValues(t, "1.3419", f.String())
	})

	t.Run("String - zero", func(t *testing.T) {
		f := Decimal{}
		assert.EqualValues(t, "0", f.String())

	})

	t.Run("FormattedString", func(t *testing.T) {
		f := NewDecimal(1.3419)

		assert.EqualValues(t, "1.3419", f.FormattedString(4))
		assert.EqualValues(t, "1.341900", f.FormattedString(6))
		assert.EqualValues(t, "1.3", f.FormattedString(1))
	})

	t.Run("ZERO", func(t *testing.T) {
		f := ZERO
		f = f.Add(ONE)

		assert.EqualValues(t, 1, f.Float())
		assert.EqualValues(t, 0, ZERO.Float())
		assert.EqualValues(t, 1, ONE.Float())
	})

	t.Run("IsZero -- when nil", func(t *testing.T) {
		f := Decimal{}

		assert.True(t, f.IsZero())
	})

	t.Run("IsZero -- when zero", func(t *testing.T) {
		f := NewFromString("0")

		assert.True(t, f.IsZero())
	})

	t.Run("IsZero -- when not zero", func(t *testing.T) {
		f := NewFromString("1")

		assert.False(t, f.IsZero())
	})
}

func TestDecimal_Json(t *testing.T) {
	type jsonType struct {
		Decimal Decimal `json:"decimal"`
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		tmpStruct := jsonType{
			Decimal: ONE,
		}
		marshaled, err := json.Marshal(tmpStruct)

		assert.NoError(t, err)
		assert.Equal(t, `{"decimal":1}`, string(marshaled))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var ts jsonType

		d := `{"decimal":1.23}`
		err := json.Unmarshal([]byte(d), &ts)

		assert.NoError(t, err)
		assert.Equal(t, "1.23", ts.Decimal.String())
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
