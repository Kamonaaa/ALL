package main

import "testing"

func run(input string) string {
	return format(process(tokenize(input)))
}

func TestProjectCases(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Simply add 42 (hex) and 10 (bin) and you will see the result is 68.", "Simply add 66 and 2 and you will see the result is 68."},
		{"Ready set go (up)!", "Ready set GO!"},
		{"I should stop SHOUTING (low)", "I should stop shouting"},
		{"welcome to brooklyn bridge (cap)", "welcome to brooklyn Bridge"},
		{"This is so exciting (up, 2)", "This is so EXCITING"},
		{"hello world (up,2)", "HELLO WORLD"},
	}
	for _, tc := range cases {
		got := run(tc.in)
		if got != tc.want {
			t.Errorf("input=%q\nwant=%q\ngot =%q", tc.in, tc.want, got)
		}
	}
}

func TestPunctuation(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello , world !", "hello, world!"},
		{"I was thinking ... You were right", "I was thinking... You were right"},
		{"what ! ?", "what!?"},
		{"Wait . . . now", "Wait... now"},
	}
	for _, tc := range cases {
		if got := run(tc.in); got != tc.want {
			t.Errorf("%q => %q (want %q)", tc.in, got, tc.want)
		}
	}
}

func TestQuotesAndApostrophes(t *testing.T) {
	cases := []struct{ in, want string }{
		{"I am ' awesome ' !", "I am 'awesome'!"},
		{"He said ' I am ready now '", "He said 'I am ready now'"},
		{"don't stop", "don't stop"},
		{"rock ' n ' roll", "rock 'n' roll"},
	}
	for _, tc := range cases {
		if got := run(tc.in); got != tc.want {
			t.Errorf("%q => %q (want %q)", tc.in, got, tc.want)
		}
	}
}

func TestArticles(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a apple", "an apple"},
		{"A hour", "An hour"},
		{"a banana", "a banana"},
		{"a honest man", "an honest man"},
	}
	for _, tc := range cases {
		if got := run(tc.in); got != tc.want {
			t.Errorf("%q => %q (want %q)", tc.in, got, tc.want)
		}
	}
}

func TestSafetyAndInvalidInputs(t *testing.T) {
	cases := []struct{ in, want string }{
		{"(up) hello", "hello"},
		{"GG (hex)", "GG"},
		{"102 (bin)", "102"},
		{"hello world (dance,4)", "hello world (dance,4)"},
		{"hello world (up,0)", "hello WORLD"},
		{"hello world (up,-3)", "hello WORLD"},
		{"hello (up", "hello (up"},
	}
	for _, tc := range cases {
		if got := run(tc.in); got != tc.want {
			t.Errorf("%q => %q (want %q)", tc.in, got, tc.want)
		}
	}
}

func TestUnicode(t *testing.T) {
	cases := []struct{ in, want string }{
		{"éclair (cap)", "Éclair"},
		{"straße (up)", "STRAßE"},
	}
	for _, tc := range cases {
		if got := run(tc.in); got != tc.want {
			t.Errorf("%q => %q (want %q)", tc.in, got, tc.want)
		}
	}
}
