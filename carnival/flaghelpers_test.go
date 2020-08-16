package carnival_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mcenirm-go/wip-cautious-carnival/carnival"
	"github.com/stretchr/testify/assert"
)

var _ carnival.ListOfUintsForFlag

func TestListOfUintsForFlag(t *testing.T) {
	x := carnival.NewListOfUintsForFlag(1, 2, 3, 4, 5)
	t.Run("initial values", func(t *testing.T) {
		assert.Equal(t, "1,2,3,4,5", x.String())
		want := []uint{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(want, x.Values()) {
			t.Errorf("values are not equal, want %v, got %v", want, x.Values())
		}
	})
	t.Run("appending values", func(t *testing.T) {
		// first time should clear initial values
		err := x.Set("9,8,7")
		assert.NoError(t, err)
		assert.Equal(t, "9,8,7", x.String())
		// subsequent times appends values
		err = x.Set("6,5,4")
		assert.NoError(t, err)
		assert.Equal(t, "9,8,7,6,5,4", x.String())
	})
	for _, badString := range []string{"", "-1", "potato"} {
		t.Run(fmt.Sprintf("bad string: %q", badString), func(t *testing.T) {
			err := x.Set(badString)
			assert.Error(t, err)
		})
	}
}
