package simple

// Code generated. DO NOT EDIT.
import "unsafe"

// AnimalFromString parses a Animal from its string representation.
func AnimalFromString(str string) Animal {
	switch str {
	case animalStringZebra:
		return AnimalZebra
	case animalStringGiraffe:
		return AnimalGiraffe
	case animalStringCoyote:
		return AnimalCoyote
	default:
		return AnimalUnknown
	}
}

// String returns the string representation of a Animal.
func (animal Animal) String() string {
	switch animal {
	case AnimalZebra:
		return animalStringZebra
	case AnimalGiraffe:
		return animalStringGiraffe
	case AnimalCoyote:
		return animalStringCoyote
	default:
		return animalStringUnknown
	}
}

// GoString implements fmt.GoStringer.
func (animal Animal) GoString() string {
	return animal.String()
}

// MarshalText implements encoding.TextMarshaler.
//
// The returned byte slice is read-only and writing to it will panic.
func (animal Animal) MarshalText() ([]byte, error) {
	str := animal.String()
	return *(*[]byte)(unsafe.Pointer(&str)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (animal *Animal) UnmarshalText(text []byte) error {
	*animal = AnimalFromString(*(*string)(unsafe.Pointer(&text)))
	return nil
}
