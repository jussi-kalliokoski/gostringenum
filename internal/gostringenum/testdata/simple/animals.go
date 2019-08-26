package simple

// Animal ...
type Animal int

// ...
const (
	AnimalUnknown Animal = iota
	AnimalZebra
	AnimalGiraffe
	AnimalCoyote
)

var (
	animalStringUnknown = "unknown"
	animalStringZebra   = "zebra"
	animalStringGiraffe = "giraffe"
	animalStringCoyote  = "coyote"
)
