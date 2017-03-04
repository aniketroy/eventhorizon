package metaevents

import (
	"testing"
)

func EqualString(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Fatalf("exp=%v; got=%v", expected, actual)
	}
}

func EqualInt(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Fatalf("exp=%v; got=%v", expected, actual)
	}
}

func TestUnknown(t *testing.T) {
	isMeta, line, event := Parse(".poop {\"foo\": \"bar\"}")

	if !isMeta {
		t.Fatalf("Expecting is meta event")
	}

	_, castSucceeded := event.(Rotated)

	if castSucceeded {
		t.Fatalf("Casting must not succeed")
	}

	EqualString(t, line, ".poop {\"foo\": \"bar\"}")
}

func TestRegularText(t *testing.T) {
	isMeta, line, _ := Parse("foobar")

	if isMeta {
		t.Fatalf("Must not be meta line")
	}

	EqualString(t, line, "foobar")
}

func TestEmptyLine(t *testing.T) {
	isMeta, line, _ := Parse("")

	if isMeta {
		t.Fatalf("Must not be meta line")
	}

	EqualString(t, line, "")
}

func TestDotEscapedRegularLine(t *testing.T) {
	isMeta, line, _ := Parse("\\.Rotated")

	if isMeta {
		t.Fatalf("Must not be meta line")
	}

	EqualString(t, line, ".Rotated")
}

func TestBackslashEscapedRegularLine(t *testing.T) {
	isMeta, line, _ := Parse("\\\\foo")

	if isMeta {
		t.Fatalf("Must not be meta line")
	}

	EqualString(t, line, "\\foo")
}

func TestNotMetaEvent(t *testing.T) {
	isMeta, line, _ := Parse("yes oh hai")

	if isMeta {
		t.Fatalf("Must not be detected as meta event")
	}

	EqualString(t, line, "yes oh hai")
}
