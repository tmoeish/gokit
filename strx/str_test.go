package strx_test

import (
	"testing"

	"github.com/tmoeish/gokit/strx"
)

func TestIsEmpty(t *testing.T) {
	if !strx.IsEmpty("") {
		t.Fatal("IsEmpty('')")
	}
	if strx.IsEmpty("x") {
		t.Fatal("IsEmpty('x')")
	}
}

func TestIsBlank(t *testing.T) {
	if !strx.IsBlank("   ") {
		t.Fatal("IsBlank spaces")
	}
	if strx.IsBlank("a") {
		t.Fatal("IsBlank 'a'")
	}
}

func TestCoalesce(t *testing.T) {
	if strx.Coalesce("", "  ", "hello") != "hello" {
		t.Fatal("Coalesce")
	}
}

func TestTruncate(t *testing.T) {
	if strx.Truncate("hello world", "...", 5) != "hello..." {
		t.Fatal("Truncate")
	}
	if strx.Truncate("hi", "...", 10) != "hi" {
		t.Fatal("Truncate no-op")
	}
}

func TestCapitalize(t *testing.T) {
	if strx.Capitalize("hello") != "Hello" {
		t.Fatal("Capitalize")
	}
	if strx.Uncapitalize("Hello") != "hello" {
		t.Fatal("Uncapitalize")
	}
}

func TestReverse(t *testing.T) {
	if strx.Reverse("hello") != "olleh" {
		t.Fatal("Reverse")
	}
	if strx.Reverse("日本語") != "語本日" {
		t.Fatal("Reverse unicode")
	}
}

func TestBeforeAfter(t *testing.T) {
	if strx.Before("hello.world", ".") != "hello" {
		t.Fatal("Before")
	}
	if strx.After("hello.world", ".") != "world" {
		t.Fatal("After")
	}
}

func TestBetween(t *testing.T) {
	if strx.Between("[hello]", "[", "]") != "hello" {
		t.Fatal("Between")
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"hello world", "hello_world"},
		{"already_snake", "already_snake"},
	}
	for _, tc := range tests {
		if got := strx.ToSnakeCase(tc.in); got != tc.want {
			t.Fatalf("ToSnakeCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"hello_world", "helloWorld"},
		{"snake_case_test", "snakeCaseTest"},
	}
	for _, tc := range tests {
		if got := strx.ToCamelCase(tc.in); got != tc.want {
			t.Fatalf("ToCamelCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	if strx.ToPascalCase("hello_world") != "HelloWorld" {
		t.Fatal("ToPascalCase")
	}
}

func TestToKebabCase(t *testing.T) {
	if strx.ToKebabCase("helloWorld") != "hello-world" {
		t.Fatal("ToKebabCase")
	}
}

func TestSlugify(t *testing.T) {
	if strx.Slugify("Hello World!") != "hello-world" {
		t.Fatal("Slugify")
	}
}

func TestIsNumeric(t *testing.T) {
	if !strx.IsNumeric("12345") {
		t.Fatal("IsNumeric")
	}
	if strx.IsNumeric("123a5") {
		t.Fatal("IsNumeric with alpha")
	}
}

func TestContainsAny(t *testing.T) {
	if !strx.ContainsAny("hello world", "xyz", "world") {
		t.Fatal("ContainsAny")
	}
}

func TestContainsAll(t *testing.T) {
	if !strx.ContainsAll("hello world", "hello", "world") {
		t.Fatal("ContainsAll")
	}
	if strx.ContainsAll("hello world", "hello", "foo") {
		t.Fatal("ContainsAll missing")
	}
}

func TestSplitAndTrim(t *testing.T) {
	got := strx.SplitAndTrim(" a , b ,  ,c ", ",")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("SplitAndTrim: %v", got)
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	if strx.LongestCommonPrefix("flower", "flow", "flight") != "fl" {
		t.Fatal("LongestCommonPrefix")
	}
}

func TestInitials(t *testing.T) {
	if strx.Initials("John Doe") != "JD" {
		t.Fatal("Initials")
	}
}
