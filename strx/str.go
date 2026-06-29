// Package strx provides string utilities.
// Inspired by Apache Commons Lang StringUtils and Java Guava's Strings.
package strx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsEmpty reports whether s has zero length.
func IsEmpty(s string) bool { return len(s) == 0 }

// IsNotEmpty reports whether s has non-zero length.
func IsNotEmpty(s string) bool { return len(s) > 0 }

// IsBlank reports whether s contains only whitespace (or is empty).
func IsBlank(s string) bool { return strings.TrimSpace(s) == "" }

// IsNotBlank reports whether s contains at least one non-whitespace character.
func IsNotBlank(s string) bool { return !IsBlank(s) }

// Coalesce returns the first non-blank string.
func Coalesce(strs ...string) string {
	for _, s := range strs {
		if IsNotBlank(s) {
			return s
		}
	}
	return ""
}

// Default returns s if not blank, otherwise def.
func Default(s, def string) string {
	if IsBlank(s) {
		return def
	}
	return s
}

// Truncate truncates s to at most n runes, appending suffix if truncated.
func Truncate(s, suffix string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + suffix
}

// TruncateBytes truncates s to at most n bytes. Avoids splitting multi-byte runes.
func TruncateBytes(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:n])
		if r != utf8.RuneError {
			break
		}
		n -= size
	}
	return s[:n]
}

// PadLeft pads s on the left with pad until its rune length is width.
func PadLeft(s, pad string, width int) string {
	runes := []rune(s)
	padRunes := []rune(pad)
	if len(padRunes) == 0 || len(runes) >= width {
		return s
	}
	var b strings.Builder
	for i := len(runes); i < width; i += len(padRunes) {
		b.WriteString(pad)
	}
	result := b.String()
	result = string([]rune(result)[:width-len(runes)]) + s
	return result
}

// PadRight pads s on the right with pad until its rune length is width.
func PadRight(s, pad string, width int) string {
	runes := []rune(s)
	padRunes := []rune(pad)
	if len(padRunes) == 0 || len(runes) >= width {
		return s
	}
	var b strings.Builder
	b.WriteString(s)
	for b.Len() < width*len(pad) && len([]rune(b.String())) < width {
		b.WriteString(pad)
	}
	return string([]rune(b.String())[:width])
}

// PadCenter pads s on both sides with pad until its rune length is width.
func PadCenter(s, pad string, width int) string {
	runes := []rune(s)
	missing := width - len(runes)
	if missing <= 0 {
		return s
	}
	leftPad := missing / 2
	rightPad := missing - leftPad
	return PadRight(PadLeft(s, pad, len(runes)+leftPad), pad, width) + strings.Repeat(pad, rightPad-0)[:rightPad]
}

// Capitalize returns s with the first rune uppercased.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// Uncapitalize returns s with the first rune lowercased.
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

// Reverse returns the reversed string (rune-aware).
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CountOccurrences counts non-overlapping occurrences of sub in s.
func CountOccurrences(s, sub string) int {
	if sub == "" {
		return 0
	}
	return strings.Count(s, sub)
}

// Before returns the substring before the first occurrence of sep.
func Before(s, sep string) string {
	i := strings.Index(s, sep)
	if i < 0 {
		return s
	}
	return s[:i]
}

// After returns the substring after the first occurrence of sep.
func After(s, sep string) string {
	i := strings.Index(s, sep)
	if i < 0 {
		return ""
	}
	return s[i+len(sep):]
}

// BeforeLast returns the substring before the last occurrence of sep.
func BeforeLast(s, sep string) string {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return s
	}
	return s[:i]
}

// AfterLast returns the substring after the last occurrence of sep.
func AfterLast(s, sep string) string {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return ""
	}
	return s[i+len(sep):]
}

// Between returns the substring between start and end.
func Between(s, start, end string) string {
	a := After(s, start)
	if a == "" {
		return ""
	}
	return Before(a, end)
}

// Wrap wraps s with prefix and suffix.
func Wrap(s, prefix, suffix string) string {
	return prefix + s + suffix
}

// Unwrap removes prefix and suffix from s if present.
func Unwrap(s, prefix, suffix string) string {
	if strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix) {
		return s[len(prefix) : len(s)-len(suffix)]
	}
	return s
}

// ContainsAny reports whether s contains any of the given substrings.
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether s contains all of the given substrings.
func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// Words splits s into words by whitespace and punctuation.
func Words(s string) []string {
	var words []string
	var word strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			word.WriteRune(r)
		} else if word.Len() > 0 {
			words = append(words, word.String())
			word.Reset()
		}
	}
	if word.Len() > 0 {
		words = append(words, word.String())
	}
	return words
}

// ToCamelCase converts a string to camelCase (lowerCamelCase).
// Handles spaces, underscores, hyphens, and existing case.
func ToCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	for i, w := range words {
		if i == 0 {
			b.WriteString(strings.ToLower(w))
		} else {
			b.WriteString(titleWord(w))
		}
	}
	return b.String()
}

// ToPascalCase converts a string to PascalCase (UpperCamelCase).
func ToPascalCase(s string) string {
	words := splitWords(s)
	var b strings.Builder
	for _, w := range words {
		b.WriteString(titleWord(w))
	}
	return b.String()
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(s string) string {
	words := splitWords(s)
	lowers := make([]string, len(words))
	for i, w := range words {
		lowers[i] = strings.ToLower(w)
	}
	return strings.Join(lowers, "_")
}

// ToKebabCase converts a string to kebab-case.
func ToKebabCase(s string) string {
	words := splitWords(s)
	lowers := make([]string, len(words))
	for i, w := range words {
		lowers[i] = strings.ToLower(w)
	}
	return strings.Join(lowers, "-")
}

// ToScreamingSnakeCase converts a string to SCREAMING_SNAKE_CASE.
func ToScreamingSnakeCase(s string) string {
	return strings.ToUpper(ToSnakeCase(s))
}

// Slugify converts a string into a URL-friendly slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r) || r == '_' || r == '-':
			if b.Len() > 0 {
				last := []rune(b.String())
				if last[len(last)-1] != '-' {
					b.WriteRune('-')
				}
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// Initials returns the initials of each word in s.
func Initials(s string) string {
	words := Words(s)
	var b strings.Builder
	for _, w := range words {
		r, _ := utf8.DecodeRuneInString(w)
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// Repeat repeats s n times (alias for strings.Repeat).
func Repeat(s string, n int) string {
	return strings.Repeat(s, n)
}

// IsNumeric reports whether every rune in s is a decimal digit.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha reports whether every rune in s is a letter.
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric reports whether every rune in s is a letter or digit.
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// SplitAndTrim splits s by sep and trims whitespace from each part.
// Empty parts are removed.
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// LongestCommonPrefix returns the longest common prefix of all strings.
func LongestCommonPrefix(strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// splitWords splits s into words for case conversion.
// Handles camelCase, PascalCase, snake_case, kebab-case, and spaces.
func splitWords(s string) []string {
	var words []string
	var word strings.Builder

	runes := []rune(s)
	for i, r := range runes {
		switch {
		case r == '_' || r == '-' || r == ' ' || r == '.' || r == '/':
			if word.Len() > 0 {
				words = append(words, word.String())
				word.Reset()
			}
		case unicode.IsUpper(r) && i > 0:
			// Check for transitions like "camelCase" or "XMLParser"
			prev := runes[i-1]
			if unicode.IsLower(prev) || (i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
				if word.Len() > 0 {
					words = append(words, word.String())
					word.Reset()
				}
			}
			word.WriteRune(r)
		default:
			word.WriteRune(r)
		}
	}
	if word.Len() > 0 {
		words = append(words, word.String())
	}
	return words
}

// titleWord uppercases the first rune of w and lowercases the rest.
func titleWord(w string) string {
	if w == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(w)
	return string(unicode.ToUpper(r)) + strings.ToLower(w[size:])
}
