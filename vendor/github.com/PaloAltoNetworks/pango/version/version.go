// Package version contains a version number struct that pango uses to make
// decisions on the specific structs to use when sending XML to the PANOS
// device.
package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Number is the version number struct.
type Number struct {
	Major, Minor, Patch int
	Suffix              string
}

// Gte tests if this version number is greater than or equal to the argument.
func (v Number) Gte(o Number) bool {
	if v.Major != o.Major {
		return v.Major > o.Major
	}

	if v.Minor != o.Minor {
		return v.Minor > o.Minor
	}

	return v.Patch >= o.Patch
}

// String returns the version number as a string.
func (v Number) String() string {
	if v.Suffix == "" {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else {
		return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Suffix)
	}
}

// New returns a version number from the given string.
func New(version string) (Number, error) {
	parts := strings.Split(version, ".")[:3]

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Number{}, fmt.Errorf("Major %s is not a number: %s", parts[0], err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Number{}, fmt.Errorf("Minor %s is not a number: %s", parts[0], err)
	}

	var patch_str string
	var suffix string
	patch_parts := strings.Split(parts[2], "-")
	if len(patch_parts) == 1 {
		patch_str = parts[2]
		suffix = ""
	} else if len(patch_parts) == 2 {
		patch_str = patch_parts[0]
		suffix = patch_parts[1]
	} else {
		return Number{}, fmt.Errorf("Patch %s is not formatted as expected", parts[2])
	}
	patch, err := strconv.Atoi(patch_str)
	if err != nil {
		return Number{}, fmt.Errorf("Patch %s is not a number: %s", patch_str, err)
	}

	return Number{major, minor, patch, suffix}, nil
}
