package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// wantsJSON checks if the client wants JSON responses.
func wantsJSON(r *http.Request) bool {
	if r.Header.Get("Accept") == "application/json" {
		return true
	}
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return false
}

// writeRecord writes a single record as plain text or JSON.
func writeRecord(w http.ResponseWriter, r *http.Request, fields map[string]string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fields)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var parts []string
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
}

// writeMultiLine writes multiple records, one per line (plain text) or as JSON array.
func writeMultiLine(w http.ResponseWriter, r *http.Request, lines []string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		arr := make([]map[string]string, len(lines))
		for i, line := range lines {
			arr[i] = parseFields(line)
		}
		json.NewEncoder(w).Encode(arr)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}

// writeRaw writes raw text content (for lorem ipsum, etc.).
func writeRaw(w http.ResponseWriter, r *http.Request, content string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"text": content})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, content)
}

// writeError writes an error response with a hint.
func writeError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]string{"error": msg})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "error: %s\n", msg)
}

func parseFields(line string) map[string]string {
	m := make(map[string]string)
	for _, part := range strings.Fields(line) {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return m
}
