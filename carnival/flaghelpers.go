package carnival

import (
	"fmt"
	"strconv"
	"strings"
)

// ListOfUintsForFlag is a custom flag.Value for holding a list of unsigned integers.
type ListOfUintsForFlag struct {
	values []uint
	setYet bool
}

// NewListOfUintsForFlag creates a new ListOfUintsForFlag with initial values.
func NewListOfUintsForFlag(values ...uint) *ListOfUintsForFlag {
	return &ListOfUintsForFlag{values: values}
}

// String returns the comma-separated values.
func (lou *ListOfUintsForFlag) String() string {
	foo := make([]string, len(lou.values))
	for i, u := range lou.values {
		foo[i] = strconv.Itoa(int(u))
	}
	return fmt.Sprint(strings.Join(foo, ","))
}

// Set appends values in a comma-separated string.
// The first call clears any default values.
func (lou *ListOfUintsForFlag) Set(value string) error {
	if !lou.setYet {
		lou.values = []uint{}
		lou.setYet = true
	}
	asStrings := strings.Split(value, ",")
	for _, asString := range asStrings {
		u64, err := strconv.ParseUint(asString, 0, 32)
		if err != nil {
			return err
		}
		lou.values = append(lou.values, uint(u64))
	}
	return nil
}

// Values returns a copy of the values.
func (lou *ListOfUintsForFlag) Values() []uint {
	values := make([]uint, len(lou.values))
	copy(values, lou.values)
	return values
}
