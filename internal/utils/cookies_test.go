package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetAndClearAuthCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	SetAuthCookie(rr, "tok", false)
	res := rr.Result()
	c := res.Cookies()
	require.NotEmpty(t, c, "expected cookie to be set")
	assert.Equal(t, "auth_token", c[0].Name)
	assert.Equal(t, "tok", c[0].Value)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(c[0])
	val, err := GetAuthCookie(req)
	require.NoError(t, err)
	assert.Equal(t, "tok", val)

	rr2 := httptest.NewRecorder()
	ClearAuthCookie(rr2, false)
	res2 := rr2.Result()
	assert.NotEmpty(t, res2.Cookies())
}
