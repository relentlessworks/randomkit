package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/relentlessworks/randomkit/internal/model"
)

// Handler handles MCP (Model Context Protocol) JSON-RPC 2.0 requests.
type Handler struct{}

// New creates a new MCP handler.
func New() *Handler {
	return &Handler{}
}

// ServeHTTP handles MCP JSON-RPC 2.0 requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONRPCError(w, nil, -32700, "Parse error")
		return
	}

	method, _ := req["method"].(string)
	id := req["id"]
	params, _ := req["params"].(map[string]interface{})

	switch method {
	case "initialize":
		writeJSONRPCResult(w, id, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "randomkit",
				"version": "0.1.0",
			},
		})

	case "tools/list":
		writeJSONRPCResult(w, id, map[string]interface{}{
			"tools": toolsList(),
		})

	case "tools/call":
		result, err := h.handleToolCall(params)
		if err != nil {
			writeJSONRPCError(w, id, -32603, err.Error())
			return
		}
		writeJSONRPCResult(w, id, result)

	default:
		writeJSONRPCError(w, id, -32601, fmt.Sprintf("Method not found: %s", method))
	}
}

func toolsList() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "random_string",
			"description": "Generate a random string with configurable length and charset.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"length":  map[string]interface{}{"type": "integer", "default": 32, "description": "Length of the string (1-10000)"},
					"charset": map[string]interface{}{"type": "string", "default": "alnum", "description": "Character set: lower, upper, digits, symbols, hex, alpha, alnum, or custom"},
				},
			},
		},
		{
			"name":        "random_int",
			"description": "Generate a random integer in [min, max].",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"min": map[string]interface{}{"type": "integer", "default": 0},
					"max": map[string]interface{}{"type": "integer", "default": 100},
				},
			},
		},
		{
			"name":        "random_float",
			"description": "Generate a random float in [min, max).",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"min": map[string]interface{}{"type": "number", "default": 0.0},
					"max": map[string]interface{}{"type": "number", "default": 1.0},
				},
			},
		},
		{
			"name":        "password",
			"description": "Generate a random password with mixed character classes.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"length":  map[string]interface{}{"type": "integer", "default": 16},
					"lower":   map[string]interface{}{"type": "boolean", "default": true},
					"upper":   map[string]interface{}{"type": "boolean", "default": true},
					"digits":  map[string]interface{}{"type": "boolean", "default": true},
					"symbols": map[string]interface{}{"type": "boolean", "default": true},
				},
			},
		},
		{
			"name":        "random_hex",
			"description": "Generate a random hex token.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"bytes": map[string]interface{}{"type": "integer", "default": 32, "description": "Number of random bytes (hex-encoded)"},
				},
			},
		},
		{
			"name":        "lorem_words",
			"description": "Generate lorem ipsum words.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer", "default": 50},
				},
			},
		},
		{
			"name":        "lorem_sentences",
			"description": "Generate lorem ipsum sentences.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer", "default": 3},
				},
			},
		},
		{
			"name":        "lorem_paragraphs",
			"description": "Generate lorem ipsum paragraphs.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer", "default": 2},
				},
			},
		},
		{
			"name":        "random_choice",
			"description": "Pick one random item from a list.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"items": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Items to choose from"},
				},
				"required": []string{"items"},
			},
		},
		{
			"name":        "random_bool",
			"description": "Generate a random boolean (true/false).",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "random_color",
			"description": "Generate a random color.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"format": map[string]interface{}{"type": "string", "default": "hex", "description": "hex, rgb, or hsl"},
				},
			},
		},
		{
			"name":        "coin_flip",
			"description": "Flip a coin (heads or tails).",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "dice_roll",
			"description": "Roll dice.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer", "default": 1, "description": "Number of dice (1-100)"},
					"sides": map[string]interface{}{"type": "integer", "default": 6, "description": "Sides per die (1-1000)"},
				},
			},
		},
		{
			"name":        "random_date",
			"description": "Generate a random date between start and end.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"start": map[string]interface{}{"type": "string", "default": "2000-01-01", "description": "Start date (YYYY-MM-DD)"},
					"end":   map[string]interface{}{"type": "string", "description": "End date (YYYY-MM-DD, default: today)"},
				},
			},
		},
		{
			"name":        "uuid",
			"description": "Generate a random UUID v4.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

func (h *Handler) handleToolCall(params map[string]interface{}) (interface{}, error) {
	name, _ := params["name"].(string)
	args, _ := params["arguments"].(map[string]interface{})

	switch name {
	case "random_string":
		length := getIntArg(args, "length", 32)
		charsetName, _ := args["charset"].(string)
		charset, err := model.ResolveCharset(charsetName)
		if err != nil {
			return nil, err
		}
		result, err := model.RandomString(length, charset)
		if err != nil {
			return nil, err
		}
		return textContent(result), nil

	case "random_int":
		min := getInt64Arg(args, "min", 0)
		max := getInt64Arg(args, "max", 100)
		result, err := model.RandomInt(min, max)
		if err != nil {
			return nil, err
		}
		return textContent(strconv.FormatInt(result, 10)), nil

	case "random_float":
		min := getFloatArg(args, "min", 0.0)
		max := getFloatArg(args, "max", 1.0)
		result, err := model.RandomFloatRange(min, max)
		if err != nil {
			return nil, err
		}
		return textContent(fmt.Sprintf("%.6f", result)), nil

	case "password":
		opts := model.DefaultPasswordOptions()
		opts.Length = getIntArg(args, "length", 16)
		opts.Lower = getBoolArg(args, "lower", true)
		opts.Upper = getBoolArg(args, "upper", true)
		opts.Digits = getBoolArg(args, "digits", true)
		opts.Symbols = getBoolArg(args, "symbols", true)
		result, err := model.GeneratePassword(opts)
		if err != nil {
			return nil, err
		}
		return textContent(result), nil

	case "random_hex":
		byteLen := getIntArg(args, "bytes", 32)
		result, err := model.RandomHex(byteLen)
		if err != nil {
			return nil, err
		}
		return textContent(result), nil

	case "lorem_words":
		count := getIntArg(args, "count", 50)
		return textContent(model.LoremWords(count)), nil

	case "lorem_sentences":
		count := getIntArg(args, "count", 3)
		return textContent(model.LoremSentences(count)), nil

	case "lorem_paragraphs":
		count := getIntArg(args, "count", 2)
		return textContent(model.LoremParagraphs(count)), nil

	case "random_choice":
		itemsRaw, ok := args["items"].([]interface{})
		if !ok || len(itemsRaw) == 0 {
			return nil, fmt.Errorf("items is required and must be a non-empty array")
		}
		items := make([]string, len(itemsRaw))
		for i, v := range itemsRaw {
			items[i], _ = v.(string)
		}
		choice, err := model.Choice(items)
		if err != nil {
			return nil, err
		}
		return textContent(choice), nil

	case "random_bool":
		return textContent(strconv.FormatBool(model.RandomBool())), nil

	case "random_color":
		format, _ := args["format"].(string)
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
			return nil, fmt.Errorf("unknown format: %s | hint: use hex, rgb, or hsl", format)
		}
		return textContent(color), nil

	case "coin_flip":
		return textContent(model.CoinFlip()), nil

	case "dice_roll":
		count := getIntArg(args, "count", 1)
		sides := getIntArg(args, "sides", 6)
		rolls, err := model.DiceRoll(count, sides)
		if err != nil {
			return nil, err
		}
		total := 0
		rollStrs := make([]string, len(rolls))
		for i, r := range rolls {
			total += r
			rollStrs[i] = strconv.Itoa(r)
		}
		return textContent(fmt.Sprintf("total=%d rolls=%s", total, strings.Join(rollStrs, "+"))), nil

	case "random_date":
		startStr, _ := args["start"].(string)
		endStr, _ := args["end"].(string)
		start, err := parseDateArg(startStr, "2000-01-01")
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %s | hint: use YYYY-MM-DD", startStr)
		}
		end, err := parseDateArg(endStr, time.Now().UTC().Format("2006-01-02"))
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %s | hint: use YYYY-MM-DD", endStr)
		}
		result, err := model.RandomDate(start, end)
		if err != nil {
			return nil, err
		}
		return textContent(result.Format("2006-01-02T15:04:05Z07:00")), nil

	case "uuid":
		return textContent(model.GenerateUUID()), nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func textContent(text string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": text,
			},
		},
	}
}

func writeJSONRPCResult(w http.ResponseWriter, id interface{}, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	})
}

func writeJSONRPCError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

func getIntArg(args map[string]interface{}, key string, def int) int {
	v, ok := args[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return def
}

func getInt64Arg(args map[string]interface{}, key string, def int64) int64 {
	v, ok := args[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int:
		return int64(n)
	case int64:
		return n
	}
	return def
}

func getFloatArg(args map[string]interface{}, key string, def float64) float64 {
	v, ok := args[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	}
	return def
}

func getBoolArg(args map[string]interface{}, key string, def bool) bool {
	v, ok := args[key]
	if !ok {
		return def
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return def
}

func parseDateArg(s, def string) (time.Time, error) {
	if s == "" {
		s = def
	}
	return time.Parse("2006-01-02", s)
}
