package main

import "testing"

func TestIsBalance(t *testing.T) {
	testcases := []struct {
		in       string
		expected bool
	}{
		{"", true},
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
		{"([{}])", true},
		{"([)]", false},
		{"(", false},
		{")", false},
	}

	for _, testcase := range testcases {
		got := IsBalanced(testcase.in)

		if testcase.expected != got {
			t.Fatalf("IsBalanced(%q) = %v, want %v", testcase.in, got, testcase.expected)
		}
	}
}
