package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/relentlessworks/randomkit/internal/auth"
	"github.com/relentlessworks/randomkit/internal/model"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	auth *auth.Auth
}

// New creates a new API handler.
func New(a *auth.Auth) *Handler {
	return &Handler{auth: a}
}

// Routes returns the HTTP mux with all routes registered.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", h.health)

	// Help / self-documentation
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.help)

	// Auth
	mux.HandleFunc("/auth/request", h.authRequest)
	mux.HandleFunc("/auth/verify", h.authVerify)

	// Random string
	mux.HandleFunc("/random/string", h.randomString)

	// Random number
	mux.HandleFunc("/random/int", h.randomInt)
	mux.HandleFunc("/random/float", h.randomFloat)

	// Random password
	mux.HandleFunc("/password", h.password)

	// Random hex token
	mux.HandleFunc("/random/hex", h.randomHex)

	// Random bytes
	mux.HandleFunc("/random/bytes", h.randomBytes)

	// Lorem ipsum
	mux.HandleFunc("/lorem/words", h.loremWords)
	mux.HandleFunc("/lorem/sentences", h.loremSentences)
	mux.HandleFunc("/lorem/paragraphs", h.loremParagraphs)

	// Random choice
	mux.HandleFunc("/random/choice", h.randomChoice)
	mux.HandleFunc("/random/sample", h.randomSample)
	mux.HandleFunc("/random/shuffle", h.randomShuffle)

	// Random boolean
	mux.HandleFunc("/random/bool", h.randomBool)

	// Random color
	mux.HandleFunc("/color", h.randomColor)

	// Coin flip
	mux.HandleFunc("/coin", h.coinFlip)

	// Dice roll
	mux.HandleFunc("/dice", h.diceRoll)

	// Random date
	mux.HandleFunc("/random/date", h.randomDate)

	// UUID
	mux.HandleFunc("/uuid", h.uuid)

	// Root - quick reference
	mux.HandleFunc("/", h.root)

	return h.auth.Middleware(mux)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "status=ok service=randomkit")
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, r, http.StatusNotFound, fmt.Sprintf("unknown endpoint: %s | hint: call GET /help for the full API reference", r.URL.Path))
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "randomkit — agentic-first random data generation service")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Endpoints:")
	fmt.Fprintln(w, "  GET  /random/string?length=32&charset=alnum  — random string")
	fmt.Fprintln(w, "  GET  /random/int?min=1&max=100              — random integer")
	fmt.Fprintln(w, "  GET  /random/float?min=0&max=1              — random float")
	fmt.Fprintln(w, "  GET  /password?length=20&lower=true&upper=true&digits=true&symbols=true  — random password")
	fmt.Fprintln(w, "  GET  /random/hex?bytes=32                     — random hex token")
	fmt.Fprintln(w, "  GET  /random/bytes?bytes=16                  — random bytes (hex)")
	fmt.Fprintln(w, "  GET  /lorem/words?count=50                   — lorem ipsum words")
	fmt.Fprintln(w, "  GET  /lorem/sentences?count=3               — lorem ipsum sentences")
	fmt.Fprintln(w, "  GET  /lorem/paragraphs?count=2               — lorem ipsum paragraphs")
	fmt.Fprintln(w, "  GET  /random/choice?items=red,green,blue    — pick one random item")
	fmt.Fprintln(w, "  GET  /random/sample?items=a,b,c,d&count=2    — pick k random items")
	fmt.Fprintln(w, "  GET  /random/shuffle?items=a,b,c,d          — shuffle a list")
	fmt.Fprintln(w, "  GET  /random/bool                            — random true/false")
	fmt.Fprintln(w, "  GET  /color?format=hex                       — random color (hex|rgb|hsl)")
	fmt.Fprintln(w, "  GET  /coin                                   — coin flip (heads/tails)")
	fmt.Fprintln(w, "  GET  /dice?count=2&sides=6                   — roll dice")
	fmt.Fprintln(w, "  GET  /random/date?start=2020-01-01&end=2025-12-31  — random date")
	fmt.Fprintln(w, "  GET  /uuid                                   — random UUID v4")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "GET /help for full documentation")
}

func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, helpText)
}

const helpText = `randomkit — Agentic-First Random Data Generation Service
======================================================

Generate random data: strings, numbers, passwords, tokens, colors, dates,
lorem ipsum, choices, coin flips, dice rolls, and UUIDs. Plain text API,
agent-driven, single Go binary. No database needed — pure stateless computation.

AUTHENTICATION
--------------
Most endpoints are public (no auth required). Auth is only needed for
administrative endpoints. To get a token:

  1. POST /auth/request  — body: email=user@example.com
     Response: otp_sent=true (OTP code is logged to stderr in dev mode)
  2. POST /auth/verify    — body: email=user@example.com&code=123456
     Response: token=<bearer-token>
  3. Use: Authorization: Bearer <token>

ENDPOINTS
=========

GET /random/string?length=32&charset=alnum
  Generate a random string.
  Params: length (1-10000, default 32), charset (lower|upper|digits|symbols|hex|alpha|alnum|custom, default alnum)
  Response: value=<string> length=<n> charset=<name>

GET /random/int?min=1&max=100
  Generate a random integer in [min, max].
  Params: min (default 0), max (default 100)
  Response: value=<n> min=<min> max=<max>

GET /random/float?min=0&max=1
  Generate a random float in [min, max).
  Params: min (default 0.0), max (default 1.0)
  Response: value=<f> min=<min> max=<max>

GET /password?length=20&lower=true&upper=true&digits=true&symbols=true
  Generate a random password with mixed character classes.
  Params: length (4-256, default 16), lower (default true), upper (default true),
          digits (default true), symbols (default true)
  Response: password=<string> length=<n>

GET /random/hex?bytes=32
  Generate a random hex token.
  Params: bytes (1-4096, default 32)
  Response: token=<hex-string> bytes=<n>

GET /random/bytes?bytes=16
  Generate random bytes (hex encoded).
  Params: bytes (1-4096, default 16)
  Response: bytes=<hex-string> length=<n>

GET /lorem/words?count=50
  Generate lorem ipsum words.
  Params: count (1-1000, default 50)
  Response: raw text (words separated by spaces)

GET /lorem/sentences?count=3
  Generate lorem ipsum sentences.
  Params: count (1-100, default 3)
  Response: raw text (sentences separated by spaces)

GET /lorem/paragraphs?count=2
  Generate lorem ipsum paragraphs.
  Params: count (1-50, default 2)
  Response: raw text (paragraphs separated by blank lines)

GET /random/choice?items=red,green,blue
  Pick one random item from a comma-separated list.
  Params: items (required, comma-separated)
  Response: choice=<item>

GET /random/sample?items=a,b,c,d&count=2
  Pick k random items (without replacement).
  Params: items (required), count (default 1)
  Response: one item per line: item=<value>

GET /random/shuffle?items=a,b,c,d
  Shuffle a list of items.
  Params: items (required, comma-separated)
  Response: one item per line: item=<value>

GET /random/bool
  Generate a random boolean.
  Response: value=true (or false)

GET /color?format=hex
  Generate a random color.
  Params: format (hex|rgb|hsl, default hex)
  Response: color=<value> format=<fmt>

GET /coin
  Flip a coin.
  Response: result=heads (or tails)

GET /dice?count=2&sides=6
  Roll dice.
  Params: count (1-100, default 1), sides (1-1000, default 6)
  Response: total=<sum> rolls=<d1+d2+...> count=<n> sides=<s>

GET /random/date?start=2020-01-01&end=2025-12-31
  Generate a random date between start and end.
  Params: start (YYYY-MM-DD, default 2000-01-01), end (YYYY-MM-DD, default today)
  Response: date=<RFC3339> unix=<timestamp>

GET /uuid
  Generate a random UUID v4.
  Response: uuid=<uuid-string>

RESPONSE FORMATS
================
Default: plain text, one labeled line per record (key=value pairs).
JSON: send Accept: application/json header or add ?format=json query param.

ERRORS
======
4xx responses include a hint for self-correction:
  error: <message> | hint: <what to do next>

MCP
===
POST /mcp — Model Context Protocol (JSON-RPC 2.0) endpoint.
Supports: initialize, tools/list, tools/call.
Tools: random_string, random_int, random_float, password, random_hex,
       lorem_words, lorem_sentences, lorem_paragraphs, random_choice,
       random_bool, random_color, coin_flip, dice_roll, random_date, uuid.

EXAMPLES
========
  curl http://localhost:8787/random/string?length=16&charset=hex
  curl http://localhost:8787/password?length=24
  curl http://localhost:8787/dice?count=3&sides=20
  curl http://localhost:8787/lorem/paragraphs?count=3
  curl http://localhost:8787/color?format=hsl
  curl -H "Accept: application/json" http://localhost:8787/uuid
`

// --- Auth Handlers ---

func (h *Handler) authRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use POST to request an OTP")
		return
	}
	email := r.FormValue("email")
	if email == "" {
		writeError(w, r, http.StatusBadRequest, "email is required | hint: provide email in the POST body")
		return
	}
	code, err := h.auth.RequestOTP(email)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	// In dev mode, include the OTP in the response for convenience
	writeRecord(w, r, map[string]string{
		"otp_sent": "true",
		"email":    email,
		"code":     code, // only in dev mode; remove in production
	})
}

func (h *Handler) authVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use POST to verify an OTP")
		return
	}
	email := r.FormValue("email")
	code := r.FormValue("code")
	if email == "" || code == "" {
		writeError(w, r, http.StatusBadRequest, "email and code are required | hint: provide both email and code in the POST body")
		return
	}
	token, err := h.auth.VerifyOTP(email, code)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"token": token,
		"email": email,
	})
}

// --- Random Data Handlers ---

func (h *Handler) randomString(w http.ResponseWriter, r *http.Request) {
	length := getIntParam(r, "length", 32)
	charsetName := r.URL.Query().Get("charset")
	charset, err := model.ResolveCharset(charsetName)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	result, err := model.RandomString(length, charset)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"value":   result,
		"length":  strconv.Itoa(length),
		"charset": charsetName,
	})
}

func (h *Handler) randomInt(w http.ResponseWriter, r *http.Request) {
	min := getInt64Param(r, "min", 0)
	max := getInt64Param(r, "max", 100)
	result, err := model.RandomInt(min, max)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"value": strconv.FormatInt(result, 10),
		"min":   strconv.FormatInt(min, 10),
		"max":   strconv.FormatInt(max, 10),
	})
}

func (h *Handler) randomFloat(w http.ResponseWriter, r *http.Request) {
	min := getFloatParam(r, "min", 0.0)
	max := getFloatParam(r, "max", 1.0)
	result, err := model.RandomFloatRange(min, max)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"value": fmt.Sprintf("%.6f", result),
		"min":   fmt.Sprintf("%g", min),
		"max":   fmt.Sprintf("%g", max),
	})
}

func (h *Handler) password(w http.ResponseWriter, r *http.Request) {
	opts := model.DefaultPasswordOptions()
	opts.Length = getIntParam(r, "length", 16)
	opts.Lower = getBoolParam(r, "lower", true)
	opts.Upper = getBoolParam(r, "upper", true)
	opts.Digits = getBoolParam(r, "digits", true)
	opts.Symbols = getBoolParam(r, "symbols", true)

	result, err := model.GeneratePassword(opts)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"password": result,
		"length":   strconv.Itoa(opts.Length),
	})
}

func (h *Handler) randomHex(w http.ResponseWriter, r *http.Request) {
	byteLen := getIntParam(r, "bytes", 32)
	result, err := model.RandomHex(byteLen)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"token": result,
		"bytes": strconv.Itoa(byteLen),
	})
}

func (h *Handler) randomBytes(w http.ResponseWriter, r *http.Request) {
	byteLen := getIntParam(r, "bytes", 16)
	result, err := model.RandomBytes(byteLen)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"bytes":  result,
		"length": strconv.Itoa(byteLen),
	})
}

func (h *Handler) loremWords(w http.ResponseWriter, r *http.Request) {
	count := getIntParam(r, "count", 50)
	result := model.LoremWords(count)
	writeRaw(w, r, result)
}

func (h *Handler) loremSentences(w http.ResponseWriter, r *http.Request) {
	count := getIntParam(r, "count", 3)
	result := model.LoremSentences(count)
	writeRaw(w, r, result)
}

func (h *Handler) loremParagraphs(w http.ResponseWriter, r *http.Request) {
	count := getIntParam(r, "count", 2)
	result := model.LoremParagraphs(count)
	writeRaw(w, r, result)
}

func (h *Handler) randomChoice(w http.ResponseWriter, r *http.Request) {
	itemsStr := r.URL.Query().Get("items")
	if itemsStr == "" {
		writeError(w, r, http.StatusBadRequest, "items parameter is required | hint: provide comma-separated items, e.g. ?items=red,green,blue")
		return
	}
	items := strings.Split(itemsStr, ",")
	choice, err := model.Choice(items)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"choice": choice,
		"count":  strconv.Itoa(len(items)),
	})
}

func (h *Handler) randomSample(w http.ResponseWriter, r *http.Request) {
	itemsStr := r.URL.Query().Get("items")
	if itemsStr == "" {
		writeError(w, r, http.StatusBadRequest, "items parameter is required | hint: provide comma-separated items, e.g. ?items=a,b,c,d")
		return
	}
	items := strings.Split(itemsStr, ",")
	count := getIntParam(r, "count", 1)
	sample, err := model.Sample(items, count)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	lines := make([]string, len(sample))
	for i, s := range sample {
		lines[i] = fmt.Sprintf("item=%s", s)
	}
	writeMultiLine(w, r, lines)
}

func (h *Handler) randomShuffle(w http.ResponseWriter, r *http.Request) {
	itemsStr := r.URL.Query().Get("items")
	if itemsStr == "" {
		writeError(w, r, http.StatusBadRequest, "items parameter is required | hint: provide comma-separated items, e.g. ?items=a,b,c,d")
		return
	}
	items := strings.Split(itemsStr, ",")
	shuffled := model.Shuffle(items)
	lines := make([]string, len(shuffled))
	for i, s := range shuffled {
		lines[i] = fmt.Sprintf("item=%s", s)
	}
	writeMultiLine(w, r, lines)
}

func (h *Handler) randomBool(w http.ResponseWriter, r *http.Request) {
	result := model.RandomBool()
	writeRecord(w, r, map[string]string{
		"value": strconv.FormatBool(result),
	})
}

func (h *Handler) randomColor(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "hex"
	}
	var color string
	switch format {
	case "hex":
		color = model.RandomColor()
	case "rgb":
		color = model.RandomColorRGB()
	case "hsl":
		color = model.RandomColorHSL()
	default:
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("unknown format: %s | hint: use one of: hex, rgb, hsl", format))
		return
	}
	writeRecord(w, r, map[string]string{
		"color":  color,
		"format": format,
	})
}

func (h *Handler) coinFlip(w http.ResponseWriter, r *http.Request) {
	result := model.CoinFlip()
	writeRecord(w, r, map[string]string{
		"result": result,
	})
}

func (h *Handler) diceRoll(w http.ResponseWriter, r *http.Request) {
	count := getIntParam(r, "count", 1)
	sides := getIntParam(r, "sides", 6)
	rolls, err := model.DiceRoll(count, sides)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	total := 0
	rollStrs := make([]string, len(rolls))
	for i, r := range rolls {
		total += r
		rollStrs[i] = strconv.Itoa(r)
	}
	writeRecord(w, r, map[string]string{
		"total": strconv.Itoa(total),
		"rolls": strings.Join(rollStrs, "+"),
		"count": strconv.Itoa(count),
		"sides": strconv.Itoa(sides),
	})
}

func (h *Handler) randomDate(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := parseDate(startStr, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid start date: %s | hint: use YYYY-MM-DD format", startStr))
		return
	}
	end, err := parseDate(endStr, time.Now().UTC())
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid end date: %s | hint: use YYYY-MM-DD format", endStr))
		return
	}

	result, err := model.RandomDate(start, end)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	writeRecord(w, r, map[string]string{
		"date": result.Format(time.RFC3339),
		"unix": strconv.FormatInt(result.Unix(), 10),
	})
}

func (h *Handler) uuid(w http.ResponseWriter, r *http.Request) {
	result := model.GenerateUUID()
	writeRecord(w, r, map[string]string{
		"uuid": result,
	})
}

// --- Helpers ---

func getIntParam(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getInt64Param(r *http.Request, key string, def int64) int64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return n
}

func getFloatParam(r *http.Request, key string, def float64) float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func getBoolParam(r *http.Request, key string, def bool) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func parseDate(s string, def time.Time) (time.Time, error) {
	if s == "" {
		return def, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
