package main

import "testing"

func TestCapitalizeWord(t *testing.T) {
	var tests = []struct {
		input, result string
	}{
		{"foo", "Foo"},
		{"fooBar", "FooBar"},
		{"österreich", "Österreich"},
		{"Österreich", "Österreich"},
		{"über", "Über"},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := capitalizeWord(test.input)
			if test.result != res {
				t.Fatalf("wrong result for %q: want %q, got %q", test.input, test.result, res)
			}
		})
	}
}
