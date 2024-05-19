package quick

import (
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func CheckProperties(t *testing.T, properties map[string]any) {

	t.Parallel()

	for name, property := range properties {
		name, property := name, property
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			err := quick.Check(property, nil)
			require.NoError(t, err)
		})
	}
}
