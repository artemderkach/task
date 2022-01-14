package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	t.Run("ping service", func(t *testing.T) {

		ts := httptest.NewServer(Router())
		defer ts.Close()

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		defer func() {
			_ = resp.Body.Close()
		}()

		result, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "pong", string(result))
	})
}
