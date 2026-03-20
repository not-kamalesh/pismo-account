package amount

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Decimal represents money safely
type Decimal struct {
	Value    decimal.Decimal
	Currency string
}

// Currency exponent map (number of decimal places)
var CurrencyExponent = map[string]int{
	"USD": 2,
	"EUR": 2,
	"INR": 2,
	"JPY": 0,
	"KWD": 3,
}

func NewAmount(value string, currency string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{
		Value:    roundToCurrencyExponent(d, currency),
		Currency: currency,
	}, nil
}

func NewAmountFromFloat(f float64, currency string) Decimal {
	d := decimal.NewFromFloat(f)
	return Decimal{
		Value:    roundToCurrencyExponent(d, currency),
		Currency: currency,
	}
}

func roundToCurrencyExponent(d decimal.Decimal, currency string) decimal.Decimal {
	exp, ok := CurrencyExponent[currency]
	if !ok {
		exp = 2
	}
	return d.Round(int32(exp))
}

func (a Decimal) String() string {
	return a.Value.String()
}

// ToMinorUnit converts amount to smallest currency unit (e.g: paise for INR)
func (a Decimal) ToMinorUnit() int64 {
	exp, ok := CurrencyExponent[a.Currency]
	if !ok {
		exp = 2
	}

	multiplier := decimal.New(1, int32(exp)) // 10^exp
	minor := a.Value.Mul(multiplier)

	return minor.IntPart() // safe since already normalized
}

func (a Decimal) Add(b Decimal) (Decimal, error) {
	if a.Currency != b.Currency {
		return Decimal{}, fmt.Errorf("currency mismatch")
	}
	return Decimal{
		Value:    roundToCurrencyExponent(a.Value.Add(b.Value), a.Currency),
		Currency: a.Currency,
	}, nil
}

func (a Decimal) Sub(b Decimal) (Decimal, error) {
	if a.Currency != b.Currency {
		return Decimal{}, fmt.Errorf("currency mismatch")
	}
	return Decimal{
		Value:    roundToCurrencyExponent(a.Value.Sub(b.Value), a.Currency),
		Currency: a.Currency,
	}, nil
}
