package types

import (
	"encoding/binary"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/Peersyst/xrpl-go/binary-codec/types/interfaces"
)

const (
	// NumberBytesLength is the number of bytes for a Number type (8 mantissa + 4 exponent).
	NumberBytesLength = 12
	// DefaultValueExponent is the exponent used for canonical zero.
	DefaultValueExponent = -2147483648
	// MinExponent is the minimum allowed exponent.
	MinExponent = -32768
	// MaxExponent is the maximum allowed exponent.
	MaxExponent = 32768
)

var (
	// MinMantissa is the minimum normalized mantissa.
	MinMantissa = big.NewInt(1000000000000000)
	// MaxMantissa is the maximum normalized mantissa.
	MaxMantissa = big.NewInt(9999999999999999)
	// ErrInvalidNumberString is returned when the input string cannot be parsed as a number.
	ErrInvalidNumberString = errors.New("invalid number string")
	// ErrMantissaExponentTooLarge is returned when the number cannot be normalized.
	ErrMantissaExponentTooLarge = errors.New("mantissa and exponent are too large")
	// numberRegex matches decimal/float/scientific number strings.
	numberRegex = regexp.MustCompile(`^([-+]?)([0-9]+)(?:\.([0-9]+))?(?:[eE]([+-]?[0-9]+))?$`)
)

// Number represents an XRPL Number type used in XLS-65/66 for vault and loan amounts.
// It is encoded as 12 bytes: 8-byte signed mantissa + 4-byte signed exponent, both big-endian.
type Number struct{}

// FromJSON parses a number string or numeric type and returns the 12-byte serialized representation.
func (n *Number) FromJSON(json any) ([]byte, error) {
	var val string
	switch v := json.(type) {
	case string:
		val = v
	case uint64:
		val = strconv.FormatUint(v, 10)
	case int64:
		val = strconv.FormatInt(v, 10)
	case int:
		val = strconv.Itoa(v)
	default:
		return nil, ErrInvalidNumberString
	}

	mantissa, exponent, err := extractNumberParts(val)
	if err != nil {
		return nil, err
	}

	// Normalize
	mantissa, exponent, err = normalizeNumber(mantissa, exponent)
	if err != nil {
		return nil, err
	}

	// Encode to 12 bytes
	bytes := make([]byte, NumberBytesLength)
	writeNumberInt64BE(bytes, mantissa, 0)
	writeNumberInt32BE(bytes, exponent, 8)

	return bytes, nil
}

// ToJSON converts the 12-byte serialized Number back to a string.
func (n *Number) ToJSON(p interfaces.BinaryParser, opts ...int) (any, error) {
	b, err := p.ReadBytes(NumberBytesLength)
	if err != nil {
		return nil, err
	}

	mantissa := readNumberInt64BE(b, 0)
	exponent := readNumberInt32BE(b, 8)

	// Canonical zero
	if mantissa.Sign() == 0 && exponent == DefaultValueExponent {
		return "0", nil
	}

	if exponent == 0 {
		return mantissa.String(), nil
	}

	// Use scientific notation for small/large exponents
	if exponent < -25 || exponent > -5 {
		return mantissa.String() + "e" + strconv.Itoa(int(exponent)), nil
	}

	// Decimal rendering for -25 <= exp <= -5
	return renderNumberDecimal(mantissa, exponent), nil
}

// extractNumberParts parses a number string into mantissa, exponent, and sign.
func extractNumberParts(val string) (*big.Int, int32, error) {
	match := numberRegex.FindStringSubmatch(val)
	if match == nil {
		return nil, 0, ErrInvalidNumberString
	}

	sign := match[1]
	intPart := match[2]
	fracPart := match[3]
	expPart := match[4]

	// Remove leading zeros
	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}

	mantissaStr := intPart
	var exponent int32 = 0

	if fracPart != "" {
		mantissaStr += fracPart
		exponent -= int32(len(fracPart))
	}

	if expPart != "" {
		exp, err := strconv.ParseInt(expPart, 10, 32)
		if err != nil {
			return nil, 0, err
		}
		exponent += int32(exp)
	}

	mantissa := new(big.Int)
	mantissa.SetString(mantissaStr, 10)
	if sign == "-" {
		mantissa.Neg(mantissa)
	}

	return mantissa, exponent, nil
}

// normalizeNumber adjusts mantissa and exponent to XRPL constraints.
func normalizeNumber(mantissa *big.Int, exponent int32) (*big.Int, int32, error) {
	// Handle zero
	if mantissa.Sign() == 0 {
		return big.NewInt(0), DefaultValueExponent, nil
	}

	m := new(big.Int).Abs(mantissa)
	isNegative := mantissa.Sign() < 0
	ten := big.NewInt(10)

	// Scale up if too small
	for m.Cmp(MinMantissa) < 0 && exponent > MinExponent {
		exponent--
		m.Mul(m, ten)
	}

	// Scale down if too large
	for m.Cmp(MaxMantissa) > 0 {
		if exponent >= MaxExponent {
			return nil, 0, ErrMantissaExponentTooLarge
		}
		exponent++
		m.Div(m, ten)
	}

	if isNegative {
		m.Neg(m)
	}

	return m, exponent, nil
}

// writeNumberInt64BE writes a big.Int as a signed 64-bit big-endian integer.
func writeNumberInt64BE(b []byte, val *big.Int, offset int) {
	i64 := val.Int64()
	binary.BigEndian.PutUint64(b[offset:], uint64(i64))
}

// writeNumberInt32BE writes an int32 as big-endian.
func writeNumberInt32BE(b []byte, val int32, offset int) {
	binary.BigEndian.PutUint32(b[offset:], uint32(val))
}

// readNumberInt64BE reads a signed 64-bit big-endian integer as big.Int.
func readNumberInt64BE(b []byte, offset int) *big.Int {
	u64 := binary.BigEndian.Uint64(b[offset:])
	i64 := int64(u64)
	return big.NewInt(i64)
}

// readNumberInt32BE reads a signed 32-bit big-endian integer.
func readNumberInt32BE(b []byte, offset int) int32 {
	u32 := binary.BigEndian.Uint32(b[offset:])
	return int32(u32)
}

// renderNumberDecimal renders a mantissa and exponent as a decimal string.
func renderNumberDecimal(mantissa *big.Int, exponent int32) string {
	isNegative := mantissa.Sign() < 0
	m := new(big.Int).Abs(mantissa)
	mantissaStr := m.String()

	// Pad with zeros
	padPrefix := 27
	padSuffix := 23
	rawValue := strings.Repeat("0", padPrefix) + mantissaStr + strings.Repeat("0", padSuffix)

	offset := int(exponent) + 43
	if offset < 0 {
		offset = 0
	}
	if offset > len(rawValue) {
		offset = len(rawValue)
	}

	integerPart := strings.TrimLeft(rawValue[:offset], "0")
	if integerPart == "" {
		integerPart = "0"
	}
	fractionPart := strings.TrimRight(rawValue[offset:], "0")

	result := integerPart
	if fractionPart != "" {
		result += "." + fractionPart
	}
	if isNegative {
		result = "-" + result
	}
	return result
}

