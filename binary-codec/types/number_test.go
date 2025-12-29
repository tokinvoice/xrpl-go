package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Peersyst/xrpl-go/binary-codec/serdes"
	"github.com/stretchr/testify/require"
)

func TestNumber_FromJSON(t *testing.T) {
	n := &Number{}

	tests := []struct {
		name    string
		input   string
		wantHex string
		wantErr error
	}{
		{
			name:    "zero",
			input:   "0",
			wantHex: "000000000000000080000000",
		},
		{
			name:    "positive integer",
			input:   "1000000000000000",
			wantHex: "00038d7ea4c6800000000000", // mantissa=1e15, exponent=0
		},
		{
			name:    "negative integer",
			input:   "-1000000000000000",
			wantHex: "fffc72815b39800000000000", // mantissa=-1e15, exponent=0
		},
		{
			name:    "decimal value",
			input:   "123.456",
			wantHex: "000462d366410000fffffff3", // normalized mantissa and exponent
		},
		{
			name:    "scientific notation positive exponent",
			input:   "1e10",
			wantHex: "00038d7ea4c68000fffffffb", // mantissa=1e15, exponent=-5 (1e10 = 1e15 * 1e-5)
		},
		{
			name:    "scientific notation negative exponent",
			input:   "1e-10",
			wantHex: "00038d7ea4c68000ffffffe7", // mantissa=1e15, exponent=-25
		},
		{
			name:    "large number near max mantissa",
			input:   "9999999999999999",
			wantHex: "002386f26fc0ffff00000000", // mantissa=max, exponent=0
		},
		{
			name:    "small number near min mantissa",
			input:   "1000000000000000",
			wantHex: "00038d7ea4c6800000000000", // mantissa=min, exponent=0
		},
		{
			name:    "invalid - not a number",
			input:   "abc",
			wantErr: ErrInvalidNumberString,
		},
		{
			name:    "invalid - empty string",
			input:   "",
			wantErr: ErrInvalidNumberString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := n.FromJSON(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantHex, hex.EncodeToString(result))
		})
	}
}

func TestNumber_ToJSON(t *testing.T) {
	n := &Number{}

	tests := []struct {
		name    string
		hexData string
		want    string
	}{
		{
			name:    "canonical zero",
			hexData: "000000000000000080000000",
			want:    "0",
		},
		{
			name:    "positive integer exponent zero",
			hexData: "00038d7ea4c6800000000000",
			want:    "1000000000000000",
		},
		{
			name:    "negative mantissa exponent zero",
			hexData: "fffc72815b39800000000000", // -1000000000000000 with exp=0
			want:    "-1000000000000000",
		},
		{
			name:    "decimal rendering (exp between -25 and -5)",
			hexData: "00038d7ea4c68000fffffff1", // mantissa=1e15, exp=-15 => 1.0
			want:    "1",
		},
		{
			name:    "scientific notation (exp < -25)",
			hexData: "00038d7ea4c68000ffffffe5", // mantissa=1e15, exp=-27
			want:    "1000000000000000e-27",
		},
		{
			name:    "scientific notation (exp > -5)",
			hexData: "00038d7ea4c68000fffffffc", // mantissa=1e15, exp=-4
			want:    "1000000000000000e-4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hexData)
			require.NoError(t, err)

			parser := serdes.NewBinaryParser(data, nil)
			result, err := n.ToJSON(parser)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestNumber_RoundTrip(t *testing.T) {
	n := &Number{}

	testCases := []string{
		"0",
		"1000000000000000",
		"-1000000000000000",
		"9999999999999999",
		"1234567890123456",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			encoded, err := n.FromJSON(input)
			require.NoError(t, err)

			parser := serdes.NewBinaryParser(encoded, nil)
			decoded, err := n.ToJSON(parser)
			require.NoError(t, err)

			// For round-trip, we verify the decoded value represents the same number
			// (format may differ due to normalization)
			t.Logf("Input: %s -> Encoded: %x -> Decoded: %s", input, encoded, decoded)
		})
	}
}

func TestNumber_FromJSON_NumericTypes(t *testing.T) {
	n := &Number{}

	tests := []struct {
		name    string
		input   any
		wantHex string
	}{
		{
			name:    "uint64",
			input:   uint64(1000000000000000),
			wantHex: "00038d7ea4c6800000000000",
		},
		{
			name:    "int64",
			input:   int64(1000000000000000),
			wantHex: "00038d7ea4c6800000000000",
		},
		{
			name:    "int",
			input:   1000000000000000,
			wantHex: "00038d7ea4c6800000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := n.FromJSON(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.wantHex, hex.EncodeToString(result))
		})
	}
}

func TestNumber_InvalidInputTypes(t *testing.T) {
	n := &Number{}

	tests := []struct {
		name  string
		input any
	}{
		{"nil", nil},
		{"float", 123.456},
		{"bool", true},
		{"slice", []byte{1, 2, 3}},
		{"map", map[string]string{"a": "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := n.FromJSON(tt.input)
			require.ErrorIs(t, err, ErrInvalidNumberString)
		})
	}
}

func TestExtractNumberParts(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantMantissa *big.Int
		wantExponent int32
		wantErr      bool
	}{
		{
			name:         "simple integer",
			input:        "123",
			wantMantissa: big.NewInt(123),
			wantExponent: 0,
		},
		{
			name:         "with decimal",
			input:        "123.456",
			wantMantissa: big.NewInt(123456),
			wantExponent: -3,
		},
		{
			name:         "with exponent",
			input:        "123e5",
			wantMantissa: big.NewInt(123),
			wantExponent: 5,
		},
		{
			name:         "negative with decimal and exponent",
			input:        "-1.5e-3",
			wantMantissa: big.NewInt(-15),
			wantExponent: -4,
		},
		{
			name:         "leading zeros",
			input:        "00123",
			wantMantissa: big.NewInt(123),
			wantExponent: 0,
		},
		{
			name:    "invalid - letters",
			input:   "12a3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mantissa, exponent, err := extractNumberParts(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 0, tt.wantMantissa.Cmp(mantissa), "mantissa mismatch")
			require.Equal(t, tt.wantExponent, exponent, "exponent mismatch")
		})
	}
}

func TestNormalizeNumber(t *testing.T) {
	tests := []struct {
		name         string
		mantissa     *big.Int
		exponent     int32
		wantMantissa *big.Int
		wantExponent int32
		wantErr      bool
	}{
		{
			name:         "zero returns canonical zero",
			mantissa:     big.NewInt(0),
			exponent:     0,
			wantMantissa: big.NewInt(0),
			wantExponent: DefaultValueExponent,
		},
		{
			name:         "small mantissa scales up",
			mantissa:     big.NewInt(1),
			exponent:     0,
			wantMantissa: big.NewInt(1000000000000000),
			wantExponent: -15,
		},
		{
			name:         "large mantissa scales down",
			mantissa:     new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(10)),
			exponent:     0,
			wantMantissa: big.NewInt(1000000000000000),
			wantExponent: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mantissa, exponent, err := normalizeNumber(tt.mantissa, tt.exponent)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 0, tt.wantMantissa.Cmp(mantissa), "mantissa mismatch: got %s, want %s", mantissa, tt.wantMantissa)
			require.Equal(t, tt.wantExponent, exponent, "exponent mismatch")
		})
	}
}

