package amount

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		currency string
		expected Decimal
	}{
		{
			name:     "USD rounding",
			input:    "100.567",
			currency: "USD",
			expected: Decimal{
				Value:    decimal.RequireFromString("100.57"),
				Currency: "USD",
			},
		},
		{
			name:     "INR rounding",
			input:    "10.999",
			currency: "INR",
			expected: Decimal{
				Value:    decimal.RequireFromString("11.00"),
				Currency: "INR",
			},
		},
		{
			name:     "JPY no decimals",
			input:    "100.99",
			currency: "JPY",
			expected: Decimal{
				Value:    decimal.RequireFromString("101"),
				Currency: "JPY",
			},
		},
		{
			name:     "KWD 3 decimals",
			input:    "1.2349",
			currency: "KWD",
			expected: Decimal{
				Value:    decimal.RequireFromString("1.235"),
				Currency: "KWD",
			},
		},
		{
			name:     "Default fallback",
			input:    "5.678",
			currency: "XXX",
			expected: Decimal{
				Value:    decimal.RequireFromString("5.68"),
				Currency: "XXX",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amt, err := NewAmount(tt.input, tt.currency)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, amt)
		})
	}
}

func TestToMinorUnit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		currency string
		expected int64
	}{
		{
			name:     "INR paise",
			input:    "100.57",
			currency: "INR",
			expected: 10057,
		},
		{
			name:     "USD cents",
			input:    "10.25",
			currency: "USD",
			expected: 1025,
		},
		{
			name:     "JPY no minor unit",
			input:    "100",
			currency: "JPY",
			expected: 100,
		},
		{
			name:     "KWD 3 decimals",
			input:    "1.234",
			currency: "KWD",
			expected: 1234,
		},
		{
			name:     "Fallback currency",
			input:    "5.68",
			currency: "XXX",
			expected: 568,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amt, _ := NewAmount(tt.input, tt.currency)
			result := amt.ToMinorUnit()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		a           string
		b           string
		currencyA   string
		currencyB   string
		expected    Decimal
		expectError bool
	}{
		{
			name:      "simple add",
			a:         "10.50",
			b:         "2.25",
			currencyA: "USD",
			currencyB: "USD",
			expected: Decimal{
				Value:    decimal.RequireFromString("12.75"),
				Currency: "USD",
			},
		},
		{
			name:      "rounding add",
			a:         "1.005",
			b:         "0.005",
			currencyA: "USD",
			currencyB: "USD",
			expected: Decimal{
				Value:    decimal.RequireFromString("1.02"),
				Currency: "USD",
			},
		},
		{
			name:        "currency mismatch",
			a:           "10",
			b:           "5",
			currencyA:   "USD",
			currencyB:   "EUR",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewAmount(tt.a, tt.currencyA)
			b, _ := NewAmount(tt.b, tt.currencyB)

			result, err := a.Add(b)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name        string
		a           string
		b           string
		currencyA   string
		currencyB   string
		expected    Decimal
		expectError bool
	}{
		{
			name:      "simple subtraction",
			a:         "10.50",
			b:         "2.25",
			currencyA: "USD",
			currencyB: "USD",
			expected: Decimal{
				Value:    decimal.RequireFromString("8.25"),
				Currency: "USD",
			},
		},
		{
			name:      "negative result",
			a:         "5.00",
			b:         "10.00",
			currencyA: "USD",
			currencyB: "USD",
			expected: Decimal{
				Value:    decimal.RequireFromString("-5.00"),
				Currency: "USD",
			},
		},
		{
			name:        "currency mismatch",
			a:           "10",
			b:           "5",
			currencyA:   "USD",
			currencyB:   "INR",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewAmount(tt.a, tt.currencyA)
			b, _ := NewAmount(tt.b, tt.currencyB)

			result, err := a.Sub(b)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
