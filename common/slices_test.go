package common

import "testing"

func TestIncludesElemInArrayTrue(t *testing.T) {
	arr := []string{"one", "two"}

	result := Includes(arr, "two")
	if result != true {
		t.Error("Should be in there")
	}
}

func TestIncludesElemNotInArrayFalse(t *testing.T) {
	arr := []string{"one", "two"}

	result := Includes(arr, "three")
	if result != false {
		t.Error("Should NOT be in there")
	}
}
