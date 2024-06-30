package quick

import (
	"testing"
	"testing/quick"
)

func CheckProperties(t *testing.T, properties map[string]any) {

	t.Parallel()

	for name, property := range properties {
		name, property := name, property
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			err := quick.Check(property, nil)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
