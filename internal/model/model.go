package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
	"time"
)

// --- Random String ---

const (
	charsetLower   = "abcdefghijklmnopqrstuvwxyz"
	charsetUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigits  = "0123456789"
	charsetSymbols = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	charsetHex     = "0123456789abcdef"
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetAlnum   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// RandomString generates a random string of the given length using the given charset.
func RandomString(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be > 0 | hint: provide a positive integer for the length parameter")
	}
	if charset == "" {
		return "", fmt.Errorf("charset must not be empty | hint: use one of: lower, upper, digits, symbols, hex, alpha, alnum, or provide a custom charset")
	}
	if length > 10000 {
		return "", fmt.Errorf("length too large (max 10000) | hint: reduce the length parameter")
	}

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// ResolveCharset maps a charset name to the actual character set.
func ResolveCharset(name string) (string, error) {
	switch name {
	case "lower":
		return charsetLower, nil
	case "upper":
		return charsetUpper, nil
	case "digits":
		return charsetDigits, nil
	case "symbols":
		return charsetSymbols, nil
	case "hex":
		return charsetHex, nil
	case "alpha":
		return charsetAlpha, nil
	case "alnum":
		return charsetAlnum, nil
	case "":
		return charsetAlnum, nil // default
	default:
		// If it's not a known name, treat it as a custom charset
		return name, nil
	}
}

// --- Random Number ---

// RandomInt generates a random integer in [min, max].
func RandomInt(min, max int64) (int64, error) {
	if min > max {
		return 0, fmt.Errorf("min (%d) > max (%d) | hint: ensure min is less than or equal to max", min, max)
	}
	rangeVal := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(rangeVal))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}
	return min + n.Int64(), nil
}

// RandomFloat generates a random float64 in [0.0, 1.0).
func RandomFloat() (float64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<53))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random float: %w", err)
	}
	return float64(n.Int64()) / float64(1<<53), nil
}

// RandomFloatRange generates a random float64 in [min, max).
func RandomFloatRange(min, max float64) (float64, error) {
	if min > max {
		return 0, fmt.Errorf("min (%.2f) > max (%.2f) | hint: ensure min is less than or equal to max", min, max)
	}
	f, err := RandomFloat()
	if err != nil {
		return 0, err
	}
	return min + f*(max-min), nil
}

// --- Random Password ---

// PasswordOptions controls password generation.
type PasswordOptions struct {
	Length  int
	Lower   bool
	Upper   bool
	Digits  bool
	Symbols bool
}

// DefaultPasswordOptions returns sensible defaults for strong passwords.
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:  16,
		Lower:   true,
		Upper:   true,
		Digits:  true,
		Symbols: true,
	}
}

// GeneratePassword creates a random password with the given options.
func GeneratePassword(opts PasswordOptions) (string, error) {
	if opts.Length < 4 {
		return "", fmt.Errorf("password length must be >= 4 | hint: increase the length parameter")
	}
	if opts.Length > 256 {
		return "", fmt.Errorf("password length too large (max 256) | hint: reduce the length parameter")
	}

	var charset strings.Builder
	if opts.Lower {
		charset.WriteString(charsetLower)
	}
	if opts.Upper {
		charset.WriteString(charsetUpper)
	}
	if opts.Digits {
		charset.WriteString(charsetDigits)
	}
	if opts.Symbols {
		charset.WriteString(charsetSymbols)
	}

	if charset.Len() == 0 {
		return "", fmt.Errorf("at least one character class must be enabled | hint: enable lower, upper, digits, or symbols")
	}

	cs := charset.String()
	pw, err := RandomString(opts.Length, cs)
	if err != nil {
		return "", err
	}

	// Ensure at least one char from each enabled class
	result := []byte(pw)
	ensureCharClass(&result, opts.Lower, charsetLower, 0)
	ensureCharClass(&result, opts.Upper, charsetUpper, 1)
	ensureCharClass(&result, opts.Digits, charsetDigits, 2)
	ensureCharClass(&result, opts.Symbols, charsetSymbols, 3)

	// Shuffle
	mrand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result), nil
}

func ensureCharClass(pw *[]byte, enabled bool, charset string, pos int) {
	if !enabled || len(*pw) <= pos {
		return
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	(*pw)[pos] = charset[n.Int64()]
}

// --- Random Hex Token ---

// RandomHex generates a random hex string of the given byte length.
func RandomHex(byteLen int) (string, error) {
	if byteLen <= 0 {
		return "", fmt.Errorf("bytes must be > 0 | hint: provide a positive integer for the byte length")
	}
	if byteLen > 4096 {
		return "", fmt.Errorf("byte length too large (max 4096) | hint: reduce the bytes parameter")
	}
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random hex: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// --- Random Bytes (hex encoded) ---

// RandomBytes generates random bytes and returns them hex encoded.
func RandomBytes(byteLen int) (string, error) {
	if byteLen <= 0 {
		return "", fmt.Errorf("bytes must be > 0 | hint: provide a positive integer for the byte length")
	}
	if byteLen > 4096 {
		return "", fmt.Errorf("byte length too large (max 4096) | hint: reduce the bytes parameter")
	}
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

// --- Lorem Ipsum ---

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
	"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
	"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla",
	"pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident",
	"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id",
	"est", "laborum",
}

// LoremWords generates n random lorem ipsum words.
func LoremWords(n int) string {
	if n <= 0 {
		n = 1
	}
	if n > 1000 {
		n = 1000
	}
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = loremWords[mrand.Intn(len(loremWords))]
	}
	return strings.Join(words, " ")
}

// LoremSentence generates one lorem ipsum sentence (8-20 words).
func LoremSentence() string {
	n := 8 + mrand.Intn(13)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = loremWords[mrand.Intn(len(loremWords))]
	}
	sentence := strings.Join(words, " ")
	if len(sentence) > 0 {
		sentence = strings.ToUpper(sentence[:1]) + sentence[1:]
	}
	return sentence + "."
}

// LoremSentences generates n lorem ipsum sentences.
func LoremSentences(n int) string {
	if n <= 0 {
		n = 1
	}
	if n > 100 {
		n = 100
	}
	sentences := make([]string, n)
	for i := 0; i < n; i++ {
		sentences[i] = LoremSentence()
	}
	return strings.Join(sentences, " ")
}

// LoremParagraph generates one lorem ipsum paragraph (3-7 sentences).
func LoremParagraph() string {
	n := 3 + mrand.Intn(5)
	return LoremSentences(n)
}

// LoremParagraphs generates n lorem ipsum paragraphs.
func LoremParagraphs(n int) string {
	if n <= 0 {
		n = 1
	}
	if n > 50 {
		n = 50
	}
	paragraphs := make([]string, n)
	for i := 0; i < n; i++ {
		paragraphs[i] = LoremParagraph()
	}
	return strings.Join(paragraphs, "\n\n")
}

// --- Random Choice ---

// Choice picks a random element from a list.
func Choice(items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("items list is empty | hint: provide at least one item, comma-separated")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		return "", fmt.Errorf("failed to pick random choice: %w", err)
	}
	return items[n.Int64()], nil
}

// Sample picks k random elements from a list (without replacement).
func Sample(items []string, k int) ([]string, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items list is empty | hint: provide at least one item, comma-separated")
	}
	if k <= 0 {
		return nil, fmt.Errorf("count must be > 0 | hint: provide a positive integer for the count parameter")
	}
	if k > len(items) {
		k = len(items)
	}

	cp := make([]string, len(items))
	copy(cp, items)
	mrand.Shuffle(len(cp), func(i, j int) {
		cp[i], cp[j] = cp[j], cp[i]
	})
	return cp[:k], nil
}

// Shuffle shuffles a list of items.
func Shuffle(items []string) []string {
	cp := make([]string, len(items))
	copy(cp, items)
	mrand.Shuffle(len(cp), func(i, j int) {
		cp[i], cp[j] = cp[j], cp[i]
	})
	return cp
}

// --- Random Boolean ---

// RandomBool returns a random true/false.
func RandomBool() bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(2))
	return n.Int64() == 1
}

// --- Random Color ---

// RandomColor generates a random color in hex format.
func RandomColor() string {
	r := mrand.Intn(256)
	g := mrand.Intn(256)
	b := mrand.Intn(256)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// RandomColorRGB generates a random color in rgb() format.
func RandomColorRGB() string {
	r := mrand.Intn(256)
	g := mrand.Intn(256)
	b := mrand.Intn(256)
	return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
}

// RandomColorHSL generates a random color in hsl() format.
func RandomColorHSL() string {
	h := mrand.Intn(361)
	s := 50 + mrand.Intn(51) // 50-100%
	l := 25 + mrand.Intn(51) // 25-75%
	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l)
}

// --- Coin Flip ---

// CoinFlip returns heads or tails.
func CoinFlip() string {
	if RandomBool() {
		return "heads"
	}
	return "tails"
}

// --- Dice Roll ---

// DiceRoll rolls n dice with s sides each.
func DiceRoll(count, sides int) ([]int, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be > 0 | hint: provide a positive integer for the number of dice")
	}
	if count > 100 {
		return nil, fmt.Errorf("count too large (max 100) | hint: reduce the number of dice")
	}
	if sides <= 0 {
		return nil, fmt.Errorf("sides must be > 0 | hint: provide a positive integer for the number of sides")
	}
	if sides > 1000 {
		return nil, fmt.Errorf("sides too large (max 1000) | hint: reduce the number of sides")
	}

	results := make([]int, count)
	for i := 0; i < count; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(sides)))
		if err != nil {
			return nil, fmt.Errorf("failed to roll dice: %w", err)
		}
		results[i] = int(n.Int64()) + 1
	}
	return results, nil
}

// --- Random Date ---

// RandomDate generates a random date between start and end.
func RandomDate(start, end time.Time) (time.Time, error) {
	if start.After(end) {
		return time.Time{}, fmt.Errorf("start date is after end date | hint: ensure start is before end")
	}
	diff := end.Unix() - start.Unix()
	if diff <= 0 {
		return start, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to generate random date: %w", err)
	}
	return time.Unix(start.Unix()+n.Int64(), 0).UTC(), nil
}

// --- UUID v4 ---

// GenerateUUID generates a random UUID v4 string.
func GenerateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
