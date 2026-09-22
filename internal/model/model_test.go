package model

import (
	"strings"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	s, err := RandomString(32, charsetAlnum)
	if err != nil {
		t.Fatalf("RandomString failed: %v", err)
	}
	if len(s) != 32 {
		t.Errorf("expected length 32, got %d", len(s))
	}
}

func TestRandomStringLength(t *testing.T) {
	for _, length := range []int{1, 10, 100, 1000} {
		s, err := RandomString(length, charsetLower)
		if err != nil {
			t.Fatalf("RandomString(%d) failed: %v", length, err)
		}
		if len(s) != length {
			t.Errorf("expected length %d, got %d", length, len(s))
		}
	}
}

func TestRandomStringInvalidLength(t *testing.T) {
	_, err := RandomString(0, charsetAlnum)
	if err == nil {
		t.Error("expected error for length 0")
	}
	_, err = RandomString(-1, charsetAlnum)
	if err == nil {
		t.Error("expected error for negative length")
	}
	_, err = RandomString(10001, charsetAlnum)
	if err == nil {
		t.Error("expected error for length > 10000")
	}
}

func TestRandomStringEmptyCharset(t *testing.T) {
	_, err := RandomString(10, "")
	if err == nil {
		t.Error("expected error for empty charset")
	}
}

func TestResolveCharset(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{"lower", charsetLower},
		{"upper", charsetUpper},
		{"digits", charsetDigits},
		{"symbols", charsetSymbols},
		{"hex", charsetHex},
		{"alpha", charsetAlpha},
		{"alnum", charsetAlnum},
		{"", charsetAlnum}, // default
		{"abc", "abc"},     // custom
	}
	for _, c := range cases {
		got, err := ResolveCharset(c.name)
		if err != nil {
			t.Fatalf("ResolveCharset(%q) failed: %v", c.name, err)
		}
		if got != c.expected {
			t.Errorf("ResolveCharset(%q) = %q, want %q", c.name, got, c.expected)
		}
	}
}

func TestRandomStringCharset(t *testing.T) {
	s, err := RandomString(100, charsetDigits)
	if err != nil {
		t.Fatalf("RandomString with digits failed: %v", err)
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			t.Errorf("expected only digits, got %c", c)
		}
	}
}

func TestRandomInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		n, err := RandomInt(1, 10)
		if err != nil {
			t.Fatalf("RandomInt failed: %v", err)
		}
		if n < 1 || n > 10 {
			t.Errorf("expected 1-10, got %d", n)
		}
	}
}

func TestRandomIntMinMax(t *testing.T) {
	n, err := RandomInt(5, 5)
	if err != nil {
		t.Fatalf("RandomInt(5,5) failed: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5, got %d", n)
	}
}

func TestRandomIntInvalidRange(t *testing.T) {
	_, err := RandomInt(10, 5)
	if err == nil {
		t.Error("expected error for min > max")
	}
}

func TestRandomFloat(t *testing.T) {
	f, err := RandomFloat()
	if err != nil {
		t.Fatalf("RandomFloat failed: %v", err)
	}
	if f < 0.0 || f >= 1.0 {
		t.Errorf("expected [0,1), got %f", f)
	}
}

func TestRandomFloatRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		f, err := RandomFloatRange(10.0, 20.0)
		if err != nil {
			t.Fatalf("RandomFloatRange failed: %v", err)
		}
		if f < 10.0 || f >= 20.0 {
			t.Errorf("expected [10,20), got %f", f)
		}
	}
}

func TestGeneratePassword(t *testing.T) {
	pw, err := GeneratePassword(DefaultPasswordOptions())
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	if len(pw) != 16 {
		t.Errorf("expected length 16, got %d", len(pw))
	}
}

func TestGeneratePasswordLength(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Length = 32
	pw, err := GeneratePassword(opts)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	if len(pw) != 32 {
		t.Errorf("expected length 32, got %d", len(pw))
	}
}

func TestGeneratePasswordTooShort(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Length = 3
	_, err := GeneratePassword(opts)
	if err == nil {
		t.Error("expected error for length < 4")
	}
}

func TestGeneratePasswordNoClasses(t *testing.T) {
	opts := PasswordOptions{Length: 16}
	_, err := GeneratePassword(opts)
	if err == nil {
		t.Error("expected error when no character classes enabled")
	}
}

func TestGeneratePasswordOnlyLower(t *testing.T) {
	opts := PasswordOptions{Length: 20, Lower: true}
	pw, err := GeneratePassword(opts)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("expected length 20, got %d", len(pw))
	}
	for _, c := range pw {
		if c < 'a' || c > 'z' {
			t.Errorf("expected only lowercase, got %c", c)
		}
	}
}

func TestRandomHex(t *testing.T) {
	h, err := RandomHex(32)
	if err != nil {
		t.Fatalf("RandomHex failed: %v", err)
	}
	if len(h) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("expected 64 hex chars, got %d", len(h))
	}
}

func TestRandomHexInvalid(t *testing.T) {
	_, err := RandomHex(0)
	if err == nil {
		t.Error("expected error for 0 bytes")
	}
	_, err = RandomHex(4097)
	if err == nil {
		t.Error("expected error for > 4096 bytes")
	}
}

func TestRandomBytes(t *testing.T) {
	b, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("RandomBytes failed: %v", err)
	}
	if len(b) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("expected 32 hex chars, got %d", len(b))
	}
}

func TestLoremWords(t *testing.T) {
	s := LoremWords(10)
	words := strings.Fields(s)
	if len(words) != 10 {
		t.Errorf("expected 10 words, got %d", len(words))
	}
}

func TestLoremWordsZero(t *testing.T) {
	s := LoremWords(0)
	if s == "" {
		t.Error("expected at least 1 word for count=0")
	}
}

func TestLoremSentence(t *testing.T) {
	s := LoremSentence()
	if !strings.HasSuffix(s, ".") {
		t.Error("expected sentence to end with period")
	}
}

func TestLoremSentences(t *testing.T) {
	s := LoremSentences(3)
	sentences := strings.Split(s, ". ")
	// At least 3 sentences (last one ends with .)
	if len(sentences) < 3 {
		t.Errorf("expected at least 3 sentences, got %d", len(sentences))
	}
}

func TestLoremParagraph(t *testing.T) {
	s := LoremParagraph()
	if s == "" {
		t.Error("expected non-empty paragraph")
	}
}

func TestLoremParagraphs(t *testing.T) {
	s := LoremParagraphs(3)
	paragraphs := strings.Split(s, "\n\n")
	if len(paragraphs) != 3 {
		t.Errorf("expected 3 paragraphs, got %d", len(paragraphs))
	}
}

func TestChoice(t *testing.T) {
	items := []string{"a", "b", "c"}
	c, err := Choice(items)
	if err != nil {
		t.Fatalf("Choice failed: %v", err)
	}
	found := false
	for _, item := range items {
		if c == item {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Choice returned %q, not in items", c)
	}
}

func TestChoiceEmpty(t *testing.T) {
	_, err := Choice([]string{})
	if err == nil {
		t.Error("expected error for empty items")
	}
}

func TestSample(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	s, err := Sample(items, 3)
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if len(s) != 3 {
		t.Errorf("expected 3 items, got %d", len(s))
	}
}

func TestSampleMoreThanAvailable(t *testing.T) {
	items := []string{"a", "b"}
	s, err := Sample(items, 5)
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if len(s) != 2 {
		t.Errorf("expected 2 items (capped), got %d", len(s))
	}
}

func TestShuffle(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	shuffled := Shuffle(items)
	if len(shuffled) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(shuffled))
	}
	// Check all items are still present
	itemSet := make(map[string]bool)
	for _, item := range shuffled {
		itemSet[item] = true
	}
	for _, item := range items {
		if !itemSet[item] {
			t.Errorf("item %q missing after shuffle", item)
		}
	}
}

func TestRandomBool(t *testing.T) {
	// Just verify it doesn't panic and returns true or false
	for i := 0; i < 100; i++ {
		_ = RandomBool()
	}
}

func TestRandomColor(t *testing.T) {
	c := RandomColor()
	if !strings.HasPrefix(c, "#") {
		t.Errorf("expected hex color starting with #, got %s", c)
	}
	if len(c) != 7 {
		t.Errorf("expected 7 chars (#RRGGBB), got %d", len(c))
	}
}

func TestRandomColorRGB(t *testing.T) {
	c := RandomColorRGB()
	if !strings.HasPrefix(c, "rgb(") {
		t.Errorf("expected rgb() format, got %s", c)
	}
}

func TestRandomColorHSL(t *testing.T) {
	c := RandomColorHSL()
	if !strings.HasPrefix(c, "hsl(") {
		t.Errorf("expected hsl() format, got %s", c)
	}
}

func TestCoinFlip(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := CoinFlip()
		if result != "heads" && result != "tails" {
			t.Errorf("expected heads or tails, got %s", result)
		}
	}
}

func TestDiceRoll(t *testing.T) {
	rolls, err := DiceRoll(3, 6)
	if err != nil {
		t.Fatalf("DiceRoll failed: %v", err)
	}
	if len(rolls) != 3 {
		t.Errorf("expected 3 rolls, got %d", len(rolls))
	}
	for _, r := range rolls {
		if r < 1 || r > 6 {
			t.Errorf("expected 1-6, got %d", r)
		}
	}
}

func TestDiceRollInvalid(t *testing.T) {
	_, err := DiceRoll(0, 6)
	if err == nil {
		t.Error("expected error for count=0")
	}
	_, err = DiceRoll(1, 0)
	if err == nil {
		t.Error("expected error for sides=0")
	}
	_, err = DiceRoll(101, 6)
	if err == nil {
		t.Error("expected error for count > 100")
	}
}

func TestRandomDate(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	d, err := RandomDate(start, end)
	if err != nil {
		t.Fatalf("RandomDate failed: %v", err)
	}
	if d.Before(start) || d.After(end) {
		t.Errorf("expected date between %v and %v, got %v", start, end, d)
	}
}

func TestRandomDateInvalidRange(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := RandomDate(start, end)
	if err == nil {
		t.Error("expected error for start > end")
	}
}

func TestGenerateUUID(t *testing.T) {
	u := GenerateUUID()
	if len(u) != 36 {
		t.Errorf("expected 36 chars, got %d", len(u))
	}
	// Check format: 8-4-4-4-12
	parts := strings.Split(u, "-")
	if len(parts) != 5 {
		t.Errorf("expected 5 parts, got %d", len(parts))
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		t.Errorf("invalid UUID format: %s", u)
	}
	// Check version is 4
	if parts[2][0] != '4' {
		t.Errorf("expected version 4, got %c", parts[2][0])
	}
}

func TestGenerateUUIDUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		u := GenerateUUID()
		if seen[u] {
			t.Errorf("duplicate UUID generated: %s", u)
		}
		seen[u] = true
	}
}
