package vm

import (
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []VmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	}

	runVmTests(t, tests)
}
