package dislog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/BOOMfinity/golog/v3/gcore"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

func TestLongLogsAreSentWithoutTruncation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		message string
		stack   bool
		entries int
	}{
		{"1800_with_stack", strings.Repeat("x", 1800), true, 1},
		{"1801_with_stack", strings.Repeat("x", 1801), true, 1},
		{"1900_with_stack", strings.Repeat("x", 1900), true, 1},
		{"1900_without_stack", strings.Repeat("x", 1900), false, 1},
		{"multiple_messages", strings.Repeat("x", 12000), true, 1},
		{"unicode_boundaries", strings.Repeat("x", 3999) + strings.Repeat("😀zażółć", 1500), true, 1},
		{"combined_text_limit", strings.Repeat("x", 1900), false, 5},
		{"combined_component_limit", "short", false, 25},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var payloadMu sync.Mutex
			var payloads []testWebhookPayload
			client := webhook.New(1, "unused-test-token", webhook.WithRestClientConfigOpts(
				rest.WithHTTPClient(&http.Client{Transport: testRoundTripper(func(request *http.Request) (*http.Response, error) {
					var payload testWebhookPayload
					if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
						return nil, err
					}
					payloadMu.Lock()
					payloads = append(payloads, payload)
					payloadMu.Unlock()
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"id":"1"}`)),
					}, nil
				})}),
			))
			mu.Lock()
			previousHooks := append([]*hook(nil), hooks...)
			mu.Unlock()
			send := Init(client)
			t.Cleanup(func() {
				mu.Lock()
				hooks = previousHooks
				mu.Unlock()
			})

			var expected strings.Builder
			logger := gcore.Init("test", func(gcore.Context) {}).Clone("scope").WithAutoStack(tc.stack).
				WithHook(func(entry gcore.Context) {
					expected.WriteString(entry.Message() + "\n\n`request(42)`")
					if tc.stack {
						expected.WriteString("\n\n```" + entry.StackTrace() + "```")
					}
					expected.WriteString("\n\n-# test@scope")
				}, send)
			for i := 0; i < tc.entries; i++ {
				logger.Error().Attr("request", 42).Msg(tc.message)
			}
			FlushAll()

			payloadMu.Lock()
			defer payloadMu.Unlock()
			var actual strings.Builder
			for i, payload := range payloads {
				text, components := collectTestComponents(t, payload.Components)
				if len(text) > 4000 || components > 40 {
					t.Errorf("payload %d exceeds Discord limits: %d bytes, %d components", i, len(text), components)
				}
				actual.WriteString(text)
			}
			if actual.String() != expected.String() {
				t.Fatalf("log content lost or changed: got %d bytes, want %d", actual.Len(), expected.Len())
			}
			if expected.Len() > 4000 && len(payloads) < 2 {
				t.Fatal("oversized log was not split into multiple messages")
			}
		})
	}
}

type testWebhookPayload struct {
	Components []testComponent `json:"components"`
}

type testComponent struct {
	Type       int             `json:"type"`
	Content    string          `json:"content"`
	Components []testComponent `json:"components"`
}

func collectTestComponents(t *testing.T, components []testComponent) (string, int) {
	t.Helper()
	var text strings.Builder
	count := len(components)
	for _, component := range components {
		if !utf8.ValidString(component.Content) || strings.ContainsRune(component.Content, utf8.RuneError) {
			t.Error("UTF-8 character was broken across message parts")
		}
		text.WriteString(component.Content)
		nestedText, nestedCount := collectTestComponents(t, component.Components)
		text.WriteString(nestedText)
		count += nestedCount
	}
	return text.String(), count
}

type testRoundTripper func(*http.Request) (*http.Response, error)

func (fn testRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.Method != http.MethodPost {
		return nil, fmt.Errorf("unexpected method: %s", request.Method)
	}
	return fn(request)
}
