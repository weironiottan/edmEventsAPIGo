package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeSnippetModel struct {
	lastUpdated time.Time
	found       bool
	err         error
}

func (f *fakeSnippetModel) FindAll() ([]EdmEvent, error) { return nil, nil }

func (f *fakeSnippetModel) FindAllWithFilter(string, string) ([]EdmEvent, error) { return nil, nil }

func (f *fakeSnippetModel) LastUpdated() (time.Time, bool, error) {
	return f.lastUpdated, f.found, f.err
}

func getLastUpdated(t *testing.T, model *fakeSnippetModel) *http.Response {
	t.Helper()
	app := &application{logger: log.New(io.Discard, "", 0), dbSnippets: model}
	rec := httptest.NewRecorder()
	app.GetLastUpdated(rec, httptest.NewRequest(http.MethodGet, "/find-edm/last-updated", nil))
	return rec.Result()
}

func readBody(t *testing.T, res *http.Response) string {
	t.Helper()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestGetLastUpdated(t *testing.T) {
	t.Run("returns the scrape time as RFC3339 UTC", func(t *testing.T) {
		pacific := time.FixedZone("PDT", -7*60*60)
		res := getLastUpdated(t, &fakeSnippetModel{
			lastUpdated: time.Date(2026, 10, 14, 2, 0, 0, 0, pacific),
			found:       true,
		})

		if res.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", res.StatusCode)
		}
		if got, want := readBody(t, res), `{"lastUpdated":"2026-10-14T09:00:00Z"}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
		if got := res.Header.Get("Access-Control-Allow-Origin"); got != "*" {
			t.Errorf("CORS header = %q, want *", got)
		}
		if got := res.Header.Get("Cache-Control"); got != "public, max-age=300" {
			t.Errorf("Cache-Control = %q", got)
		}
	})

	t.Run("404s before the first scrape is recorded", func(t *testing.T) {
		res := getLastUpdated(t, &fakeSnippetModel{found: false})

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", res.StatusCode)
		}
		if got, want := readBody(t, res), `{"Description":"No scrape has been recorded yet"}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
	})

	t.Run("500s when Firestore fails, without writing a second response", func(t *testing.T) {
		res := getLastUpdated(t, &fakeSnippetModel{err: errors.New("firestore unavailable")})

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", res.StatusCode)
		}
		if got, want := readBody(t, res), "The server encountered a problem and could not process your request\n"; got != want {
			t.Errorf("body = %q, want %q", got, want)
		}
	})
}
