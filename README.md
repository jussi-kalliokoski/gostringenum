# gostringenum


[![CI status](https://github.com/jussi-kalliokoski/gostringenum/workflows/CI/badge.svg)](https://github.com/jussi-kalliokoski/gostringenum/actions)

`gostringenum` is a utility for generating encoding code for go enums that have string representations.

`gostringenum` aims to provide an efficient solution for small enumerations, but does not concern itself with enums that dozens of values. `gostringenum` also optionally generates tests for the generated code.

## Defining Enums

`gostringenum` supports two different types of enums:

### Enum With an Unknown Zero Value

An enum with an unknown zero value is a type of enum where the zero value signifies an unknown value, and the encoding/decoding can never fail, but defaults to the zero value instead. To define an enum with an unknown zero value, simply define an exported zero value:

```go
type MyEnum int

const (
  MyEnumUnknown MyEnum = iota
  MyEnumA
  MyEnumB
)

var (
  myEnumStringUnknown = "unknown"
  myEnumStringA       = "a"
  myEnumStringB       = "b"
)
```

### Enum With an Invalid Zero Value

An enum with an invalid zero value is a type of enum where any value other than the ones defined is considered invalid, and the encoding/decoding will fail upon encountering these values. To define an enum with an invalid zero value, don't define an exported zero value:

```go
type MyEnum int

const (
  myEnumInvalid MyEnum = iota
  MyEnumA
  MyEnumB
)

var (
  myEnumStringA = "a"
  myEnumStringB = "b"
)
```

## Usage

Recommended usage is to embed as a script to your project and use go:generate. This provides you with minimal setup, an easy upgrade path using standard toolchain, and a more obvious indication of where the code comes from.

```go
// ./myenum.go

package mypkg

type MyEnum int

const (
  myEnumInvalid MyEnum = iota
  MyEnumA
  MyEnumB
)

var (
  myEnumStringA = "a"
  myEnumStringB = "b"
)

//go:generate go run generate/gostringenum.go -- -test-file . MyEnum myenum_encoding.go
```

```go
// ./generate/gostringenum.go

package main

import "github.com/jussi-kalliokoski/gostringenum/cli"

func main() {
	cli.Run()
}
```

## License

MIT License. See [LICENSE](LICENSE) for more details.
