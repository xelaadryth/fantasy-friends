package rgapi

import (
	"testing"
)

func TestNormalizeGameName(t *testing.T) {
	expected := "sometestname"
	actual := NormalizeGameName("Some Test Name")
	if actual != expected {
		t.Error("Normalized:", actual, "Expected:", expected)
	}
}
