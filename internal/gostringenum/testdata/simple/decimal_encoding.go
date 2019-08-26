package simple

// Code generated. DO NOT EDIT.
import (
	"unsafe"
	"fmt"
)

// ParseDecimal parses a Decimal from its string representation.
func ParseDecimal(str string) (Decimal, error) {
	switch str {
	case decimalStringOne:
		return DecimalOne, nil
	case decimalStringTwo:
		return DecimalTwo, nil
	case decimalStringThree:
		return DecimalThree, nil
	case decimalStringFour:
		return DecimalFour, nil
	case decimalStringFive:
		return DecimalFive, nil
	case decimalStringSix:
		return DecimalSix, nil
	case decimalStringSeven:
		return DecimalSeven, nil
	case decimalStringEight:
		return DecimalEight, nil
	case decimalStringNine:
		return DecimalNine, nil
	default:
		return 0, fmt.Errorf("not a Decimal: %q", str)
	}
}

// String returns the string representation of a Decimal.
//
// Will panic if the value is not a valid Decimal.
func (decimal Decimal) String() string {
	switch decimal {
	case DecimalOne:
		return decimalStringOne
	case DecimalTwo:
		return decimalStringTwo
	case DecimalThree:
		return decimalStringThree
	case DecimalFour:
		return decimalStringFour
	case DecimalFive:
		return decimalStringFive
	case DecimalSix:
		return decimalStringSix
	case DecimalSeven:
		return decimalStringSeven
	case DecimalEight:
		return decimalStringEight
	case DecimalNine:
		return decimalStringNine
	default:
		panic(fmt.Errorf("not a Decimal: %d", decimal))
	}
}

// GoString implements fmt.GoStringer.
func (decimal Decimal) GoString() string {
	switch decimal {
	case DecimalOne:
		return decimalStringOne
	case DecimalTwo:
		return decimalStringTwo
	case DecimalThree:
		return decimalStringThree
	case DecimalFour:
		return decimalStringFour
	case DecimalFive:
		return decimalStringFive
	case DecimalSix:
		return decimalStringSix
	case DecimalSeven:
		return decimalStringSeven
	case DecimalEight:
		return decimalStringEight
	case DecimalNine:
		return decimalStringNine
	default:
		return fmt.Sprintf("Decimal{invalid %d}", decimal)
	}
}

// MarshalText implements encoding.TextMarshaler.
//
// The returned byte slice is read-only and writing to it will panic.
func (decimal Decimal) MarshalText() ([]byte, error) {
	var str string
	switch decimal {
	case DecimalOne:
		str = decimalStringOne
	case DecimalTwo:
		str = decimalStringTwo
	case DecimalThree:
		str = decimalStringThree
	case DecimalFour:
		str = decimalStringFour
	case DecimalFive:
		str = decimalStringFive
	case DecimalSix:
		str = decimalStringSix
	case DecimalSeven:
		str = decimalStringSeven
	case DecimalEight:
		str = decimalStringEight
	case DecimalNine:
		str = decimalStringNine
	default:
		return nil, fmt.Errorf("Decimal{invalid %d}", decimal)
	}
	return *(*[]byte)(unsafe.Pointer(&str)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (decimal *Decimal) UnmarshalText(text []byte) error {
	var err error
	*decimal, err = ParseDecimal(*(*string)(unsafe.Pointer(&text)))
	return err
}
