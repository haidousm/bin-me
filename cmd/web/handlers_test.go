package main

import (
	"net/http"
	"testing"

	"binme.haido.us/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestBinView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/bins/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/bins/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/bins/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/bins/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/bins/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/bins/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
