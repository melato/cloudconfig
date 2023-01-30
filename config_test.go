package cloudinit

import (
	"testing"
)

func stringSliceEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if s != b[i] {
			return false
		}
	}
	return true
}

func TestConfigMerge(t *testing.T) {
	var x, y Config
	x.Packages = []string{"a", "b"}
	y.Packages = []string{"c", "b"}
	x.Merge(&y)
	if !stringSliceEquals(x.Packages, []string{"a", "b", "c"}) {
		t.Fail()
	}
}

func TestConfigMerge2(t *testing.T) {
	var c, x, y Config
	x.Packages = []string{"a", "b"}
	y.Packages = []string{"c", "b"}
	c.Merge(&x)
	c.Merge(&y)
	if !stringSliceEquals(c.Packages, []string{"a", "b", "c"}) {
		t.Fail()
	}
}