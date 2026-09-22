package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setup creates a test handler and returns it.
func setup(t *testing.T) *Handler {
	return New(nil) // auth not needed for public endpoints
}

func TestRoot(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.root(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "randomkit") {
		t.Error("expected body to contain 'randomkit'")
	}
}

func TestHealth(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	h.health(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "status=ok") {
		t.Error("expected body to contain 'status=ok'")
	}
}

func TestHelp(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/help", nil)
	w := httptest.NewRecorder()
	h.help(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "randomkit") {
		t.Error("expected body to contain 'randomkit'")
	}
	if !strings.Contains(body, "AUTHENTICATION") {
		t.Error("expected body to contain 'AUTHENTICATION'")
	}
}

func TestRandomString(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/string?length=16&charset=hex", nil)
	w := httptest.NewRecorder()
	h.randomString(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "value=") {
		t.Error("expected body to contain 'value='")
	}
}

func TestRandomStringDefault(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/string", nil)
	w := httptest.NewRecorder()
	h.randomString(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRandomInt(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/int?min=1&max=10", nil)
	w := httptest.NewRecorder()
	h.randomInt(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "value=") {
		t.Error("expected body to contain 'value='")
	}
}

func TestRandomIntInvalidRange(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/int?min=10&max=1", nil)
	w := httptest.NewRecorder()
	h.randomInt(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRandomFloat(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/float?min=0&max=100", nil)
	w := httptest.NewRecorder()
	h.randomFloat(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "value=") {
		t.Error("expected body to contain 'value='")
	}
}

func TestPassword(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/password?length=20", nil)
	w := httptest.NewRecorder()
	h.password(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "password=") {
		t.Error("expected body to contain 'password='")
	}
}

func TestPasswordTooShort(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/password?length=2", nil)
	w := httptest.NewRecorder()
	h.password(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRandomHex(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/hex?bytes=16", nil)
	w := httptest.NewRecorder()
	h.randomHex(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "token=") {
		t.Error("expected body to contain 'token='")
	}
}

func TestRandomBytes(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/bytes?bytes=8", nil)
	w := httptest.NewRecorder()
	h.randomBytes(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLoremWords(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/lorem/words?count=10", nil)
	w := httptest.NewRecorder()
	h.loremWords(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	words := strings.Fields(body)
	if len(words) != 10 {
		t.Errorf("expected 10 words, got %d", len(words))
	}
}

func TestLoremSentences(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/lorem/sentences?count=2", nil)
	w := httptest.NewRecorder()
	h.loremSentences(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLoremParagraphs(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/lorem/paragraphs?count=2", nil)
	w := httptest.NewRecorder()
	h.loremParagraphs(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRandomChoice(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/choice?items=red,green,blue", nil)
	w := httptest.NewRecorder()
	h.randomChoice(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "choice=") {
		t.Error("expected body to contain 'choice='")
	}
}

func TestRandomChoiceNoItems(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/choice", nil)
	w := httptest.NewRecorder()
	h.randomChoice(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRandomSample(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/sample?items=a,b,c,d&count=2", nil)
	w := httptest.NewRecorder()
	h.randomSample(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRandomShuffle(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/shuffle?items=a,b,c,d,e", nil)
	w := httptest.NewRecorder()
	h.randomShuffle(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
}

func TestRandomBool(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/bool", nil)
	w := httptest.NewRecorder()
	h.randomBool(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "value=true") && !strings.Contains(body, "value=false") {
		t.Errorf("expected value=true or value=false, got %s", body)
	}
}

func TestRandomColor(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/color?format=hex", nil)
	w := httptest.NewRecorder()
	h.randomColor(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "color=") {
		t.Error("expected body to contain 'color='")
	}
}

func TestRandomColorRGB(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/color?format=rgb", nil)
	w := httptest.NewRecorder()
	h.randomColor(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "rgb(") {
		t.Error("expected body to contain 'rgb('")
	}
}

func TestRandomColorInvalidFormat(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/color?format=invalid", nil)
	w := httptest.NewRecorder()
	h.randomColor(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCoinFlip(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/coin", nil)
	w := httptest.NewRecorder()
	h.coinFlip(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "result=heads") && !strings.Contains(body, "result=tails") {
		t.Errorf("expected heads or tails, got %s", body)
	}
}

func TestDiceRoll(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/dice?count=2&sides=6", nil)
	w := httptest.NewRecorder()
	h.diceRoll(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "total=") {
		t.Error("expected body to contain 'total='")
	}
	if !strings.Contains(body, "rolls=") {
		t.Error("expected body to contain 'rolls='")
	}
}

func TestDiceRollInvalid(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/dice?count=0&sides=6", nil)
	w := httptest.NewRecorder()
	h.diceRoll(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRandomDate(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/date?start=2020-01-01&end=2025-12-31", nil)
	w := httptest.NewRecorder()
	h.randomDate(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "date=") {
		t.Error("expected body to contain 'date='")
	}
}

func TestRandomDateInvalidFormat(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/date?start=invalid", nil)
	w := httptest.NewRecorder()
	h.randomDate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUUID(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/uuid", nil)
	w := httptest.NewRecorder()
	h.uuid(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "uuid=") {
		t.Error("expected body to contain 'uuid='")
	}
}

func TestJSONResponse(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/uuid", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	h.uuid(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "\"uuid\"") {
		t.Error("expected JSON body to contain 'uuid'")
	}
}

func TestJSONResponseViaQueryParam(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/random/bool?format=json", nil)
	w := httptest.NewRecorder()
	h.randomBool(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "\"value\"") {
		t.Error("expected JSON body to contain 'value'")
	}
}

func TestNotFound(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	h.root(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
