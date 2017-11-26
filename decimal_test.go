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

	t.Run("String", func(t *testing.T) {
		f := NewDecimal(1.3419)
		assert.EqualValues(t, "1.3419", f.String())
	})

	t.Run("ZERO", func(t *testing.T) {
		f := ZERO
		f = f.Add(ONE)

		assert.EqualValues(t, 1, f.Float())
		assert.EqualValues(t, 0, ZERO.Float())
		assert.EqualValues(t, 1, ONE.Float())
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
