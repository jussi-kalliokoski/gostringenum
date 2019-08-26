package simple_test

// Code generated. DO NOT EDIT.
import (
	"testing"
	"."
)

func TestAnimal(t *testing.T) {
	t.Run("AnimalFromString", func(t *testing.T) {
		t.Run("AnimalUnknown", func(t *testing.T) {
			expected := simple.AnimalUnknown
			received := simple.AnimalFromString("unknown")
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalZebra", func(t *testing.T) {
			expected := simple.AnimalZebra
			received := simple.AnimalFromString("zebra")
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalGiraffe", func(t *testing.T) {
			expected := simple.AnimalGiraffe
			received := simple.AnimalFromString("giraffe")
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalCoyote", func(t *testing.T) {
			expected := simple.AnimalCoyote
			received := simple.AnimalFromString("coyote")
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
	})
	t.Run("String", func(t *testing.T) {
		t.Run("AnimalUnknown", func(t *testing.T) {
			expected := "unknown"
			received := simple.AnimalUnknown.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalZebra", func(t *testing.T) {
			expected := "zebra"
			received := simple.AnimalZebra.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalGiraffe", func(t *testing.T) {
			expected := "giraffe"
			received := simple.AnimalGiraffe.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalCoyote", func(t *testing.T) {
			expected := "coyote"
			received := simple.AnimalCoyote.String()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
	})
	t.Run("GoString", func(t *testing.T) {
		t.Run("AnimalUnknown", func(t *testing.T) {
			expected := "unknown"
			received := simple.AnimalUnknown.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalZebra", func(t *testing.T) {
			expected := "zebra"
			received := simple.AnimalZebra.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalGiraffe", func(t *testing.T) {
			expected := "giraffe"
			received := simple.AnimalGiraffe.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
		t.Run("AnimalCoyote", func(t *testing.T) {
			expected := "coyote"
			received := simple.AnimalCoyote.GoString()
			if expected != received {
				t.Fatalf("expected %q, got %q", expected, received)
			}
		})
	})
	t.Run("MarshalText", func(t *testing.T) {
		t.Run("AnimalUnknown", func(t *testing.T) {
			expected := []byte("unknown")
			received, err := simple.AnimalUnknown.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalZebra", func(t *testing.T) {
			expected := []byte("zebra")
			received, err := simple.AnimalZebra.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalGiraffe", func(t *testing.T) {
			expected := []byte("giraffe")
			received, err := simple.AnimalGiraffe.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalCoyote", func(t *testing.T) {
			expected := []byte("coyote")
			received, err := simple.AnimalCoyote.MarshalText()
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if string(expected) != string(received) {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		t.Run("AnimalUnknown", func(t *testing.T) {
			expected := simple.AnimalUnknown
			var received simple.Animal
			err := received.UnmarshalText([]byte("unknown"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalZebra", func(t *testing.T) {
			expected := simple.AnimalZebra
			var received simple.Animal
			err := received.UnmarshalText([]byte("zebra"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalGiraffe", func(t *testing.T) {
			expected := simple.AnimalGiraffe
			var received simple.Animal
			err := received.UnmarshalText([]byte("giraffe"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
		t.Run("AnimalCoyote", func(t *testing.T) {
			expected := simple.AnimalCoyote
			var received simple.Animal
			err := received.UnmarshalText([]byte("coyote"))
			if err != nil {
				t.Fatalf("expected no error, got %#v", err)
			}
			if expected != received {
				t.Fatalf("expected %#v, got %#v", expected, received)
			}
		})
	})
}
