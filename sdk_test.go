package tgbot

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func alwaysOkTgServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/notify")
		assert.Equal(t, req.Method, "POST")
		rw.WriteHeader(200)
		rw.Write([]byte(`{}`))
	}))
}

func alwaysErrTgServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/notify")
		assert.Equal(t, req.Method, "POST")
		rw.WriteHeader(400)
		rw.Write([]byte(`{"message": "unknown error"}`))
	}))
}

func TestSDK(t *testing.T) {
	server := alwaysOkTgServer(t)
	defer server.Close()

	sdk := NewSDK(&http.Client{}, server.URL)

	assert.NoError(t, sdk.Send(context.Background(), time.Now(), "error", "some msg"))
}

func TestSDKHandleRemoteError(t *testing.T) {
	server := alwaysErrTgServer(t)
	defer server.Close()

	sdk := NewSDK(&http.Client{}, server.URL)

	assert.Error(t, sdk.Send(context.Background(), time.Now(), "error", "some msg"))
}
