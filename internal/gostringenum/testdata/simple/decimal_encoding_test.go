package simple_test

// Code generated. DO NOT EDIT.
import (
	"testing"
	"."
)

func TestDecimal(t *testing.T) {
	t.Run("ParseDecimal", func(t *testing.T) {
		t.Run("DecimalOne", func(t *testing.T) {
			expected := simple.DecimalOne
			received, err := simple.ParseDecimal("one")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalTwo", func(t *testing.T) {
			expected := simple.DecimalTwo
			received, err := simple.ParseDecimal("two")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalThree", func(t *testing.T) {
			expected := simple.DecimalThree
			received, err := simple.ParseDecimal("three")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFour", func(t *testing.T) {
			expected := simple.DecimalFour
			received, err := simple.ParseDecimal("four")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFive", func(t *testing.T) {
			expected := simple.DecimalFive
			received, err := simple.ParseDecimal("five")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSix", func(t *testing.T) {
			expected := simple.DecimalSix
			received, err := simple.ParseDecimal("six")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSeven", func(t *testing.T) {
			expected := simple.DecimalSeven
			received, err := simple.ParseDecimal("seven")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalEight", func(t *testing.T) {
			expected := simple.DecimalEight
			received, err := simple.ParseDecimal("eight")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalNine", func(t *testing.T) {
			expected := simple.DecimalNine
			received, err := simple.ParseDecimal("nine")
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("Invalid value", func(t *testing.T) {
			_, err := simple.ParseDecimal(string([]byte{0}))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})
	t.Run("String", func(t *testing.T) {
		t.Run("DecimalOne", func(t *testing.T) {
			expected := "one"
			received := simple.DecimalOne.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalTwo", func(t *testing.T) {
			expected := "two"
			received := simple.DecimalTwo.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalThree", func(t *testing.T) {
			expected := "three"
			received := simple.DecimalThree.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalFour", func(t *testing.T) {
			expected := "four"
			received := simple.DecimalFour.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalFive", func(t *testing.T) {
			expected := "five"
			received := simple.DecimalFive.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalSix", func(t *testing.T) {
			expected := "six"
			received := simple.DecimalSix.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalSeven", func(t *testing.T) {
			expected := "seven"
			received := simple.DecimalSeven.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalEight", func(t *testing.T) {
			expected := "eight"
			received := simple.DecimalEight.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalNine", func(t *testing.T) {
			expected := "nine"
			received := simple.DecimalNine.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("Invalid value", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic")
				}
			}()
			var decimal simple.Decimal = ^0
			_ = decimal.String()
		})
	})
	t.Run("GoString", func(t *testing.T) {
		t.Run("DecimalOne", func(t *testing.T) {
			expected := "one"
			received := simple.DecimalOne.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalTwo", func(t *testing.T) {
			expected := "two"
			received := simple.DecimalTwo.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalThree", func(t *testing.T) {
			expected := "three"
			received := simple.DecimalThree.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalFour", func(t *testing.T) {
			expected := "four"
			received := simple.DecimalFour.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalFive", func(t *testing.T) {
			expected := "five"
			received := simple.DecimalFive.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalSix", func(t *testing.T) {
			expected := "six"
			received := simple.DecimalSix.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalSeven", func(t *testing.T) {
			expected := "seven"
			received := simple.DecimalSeven.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalEight", func(t *testing.T) {
			expected := "eight"
			received := simple.DecimalEight.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("DecimalNine", func(t *testing.T) {
			expected := "nine"
			received := simple.DecimalNine.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("Invalid value", func(t *testing.T) {
			var decimal simple.Decimal = ^0
			_ = decimal.GoString()
		})
	})
	t.Run("MarshalText", func(t *testing.T) {
		t.Run("DecimalOne", func(t *testing.T) {
			expected := []byte("one")
			received, err := simple.DecimalOne.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalTwo", func(t *testing.T) {
			expected := []byte("two")
			received, err := simple.DecimalTwo.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalThree", func(t *testing.T) {
			expected := []byte("three")
			received, err := simple.DecimalThree.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFour", func(t *testing.T) {
			expected := []byte("four")
			received, err := simple.DecimalFour.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFive", func(t *testing.T) {
			expected := []byte("five")
			received, err := simple.DecimalFive.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSix", func(t *testing.T) {
			expected := []byte("six")
			received, err := simple.DecimalSix.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSeven", func(t *testing.T) {
			expected := []byte("seven")
			received, err := simple.DecimalSeven.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalEight", func(t *testing.T) {
			expected := []byte("eight")
			received, err := simple.DecimalEight.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalNine", func(t *testing.T) {
			expected := []byte("nine")
			received, err := simple.DecimalNine.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("Invalid value", func(t *testing.T) {
			var decimal simple.Decimal
			if _, err := decimal.MarshalText(); err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		t.Run("DecimalOne", func(t *testing.T) {
			expected := simple.DecimalOne
			var received simple.Decimal
			err := received.UnmarshalText([]byte("one"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalTwo", func(t *testing.T) {
			expected := simple.DecimalTwo
			var received simple.Decimal
			err := received.UnmarshalText([]byte("two"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalThree", func(t *testing.T) {
			expected := simple.DecimalThree
			var received simple.Decimal
			err := received.UnmarshalText([]byte("three"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFour", func(t *testing.T) {
			expected := simple.DecimalFour
			var received simple.Decimal
			err := received.UnmarshalText([]byte("four"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalFive", func(t *testing.T) {
			expected := simple.DecimalFive
			var received simple.Decimal
			err := received.UnmarshalText([]byte("five"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSix", func(t *testing.T) {
			expected := simple.DecimalSix
			var received simple.Decimal
			err := received.UnmarshalText([]byte("six"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalSeven", func(t *testing.T) {
			expected := simple.DecimalSeven
			var received simple.Decimal
			err := received.UnmarshalText([]byte("seven"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalEight", func(t *testing.T) {
			expected := simple.DecimalEight
			var received simple.Decimal
			err := received.UnmarshalText([]byte("eight"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("DecimalNine", func(t *testing.T) {
			expected := simple.DecimalNine
			var received simple.Decimal
			err := received.UnmarshalText([]byte("nine"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("Invalid value", func(t *testing.T) {
			var decimal simple.Decimal
			if err := decimal.UnmarshalText([]byte{0}); err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})
}
