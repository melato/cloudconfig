package cloudconfig

import (
	"embed"
	"io/fs"
	"testing"
)

//go:embed test/*.yaml
var testFS embed.FS

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

func TestConfigCommandScript(t *testing.T) {
	var c any
	c = "echo"
	s, isScript := CommandScript(c)
	if !isScript || s != "echo" {
		t.Fail()
	}
	c = []string{"echo", "a"}
	args, isArgs := CommandArgs(c)
	if !isArgs || !stringSliceEquals([]string{"echo", "a"}, args) {
		t.Fail()
	}

	c = []any{"echo", 1}
	args, isArgs = CommandArgs(c)
	if !isArgs || !stringSliceEquals([]string{"echo", "1"}, args) {
		t.Fail()
	}
}

func TestHasCloudConfigComment(t *testing.T) {
	data, err := fs.ReadFile(testFS, "test/comment.yaml")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !FirstLineIs(data, Comment) {
		t.Fatalf("did not find comment")
	}
	data, err = fs.ReadFile(testFS, "test/no-comment.yaml")
	if FirstLineIs(data, Comment) {
		t.Fatalf("shoudl not have found comment")
	}
}
