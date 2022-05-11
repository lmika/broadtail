package reqbind_test

import (
	"github.com/google/uuid"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBind(t *testing.T) {
	t.Run("should bind from query request parameters to a struct", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://www.example.com/?foo=hello&bar=world", nil)
		s := struct {
			Foo string `req:"foo"`
			Bar string `req:"bar"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, "hello", s.Foo)
		assert.Equal(t, "world", s.Bar)
	})

	t.Run("should bind to from request parameters to a struct", func(t *testing.T) {
		formParams := make(url.Values)
		formParams.Set("baz", "Bladibla")
		formParams.Set("fed", "Flimflam")

		req := httptest.NewRequest("POST", "https://www.example.com/", strings.NewReader(formParams.Encode()))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := struct {
			Foo string `req:"baz"`
			Bar string `req:"fed"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, "Bladibla", s.Foo)
		assert.Equal(t, "Flimflam", s.Bar)
	})

	t.Run("primitive types should bind to the last value", func(t *testing.T) {
		formParams := make(url.Values)
		formParams.Set("baz", "Bladibla")
		formParams.Add("baz", "Really use this one")
		formParams.Set("checkbox", "off")
		formParams.Add("checkbox", "on")

		req := httptest.NewRequest("POST", "https://www.example.com/", strings.NewReader(formParams.Encode()))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := struct {
			Foo string `req:"baz"`
			Bar bool   `req:"checkbox"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, "Really use this one", s.Foo)
		assert.True(t, s.Bar)
	})

	t.Run("should bind to ints", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://www.example.com/?foo=123&bar=-456", nil)
		s := struct {
			Foo int `req:"foo"`
			Bar int `req:"bar"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, 123, s.Foo)
		assert.Equal(t, -456, s.Bar)
	})

	t.Run("should bind to TextUnmarshaler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://www.example.com/?id=e559c6f0-ae06-4fc4-9d08-10a5994a8d11", nil)
		s := struct {
			ID uuid.UUID `req:"id"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, uuid.MustParse("e559c6f0-ae06-4fc4-9d08-10a5994a8d11"), s.ID)
	})

	t.Run("should bind to bools", func(t *testing.T) {
		scenarios := []struct {
			Param    string
			Expected bool
		}{
			{"t", true},
			{"true", true},
			{"on", true},
			{"1", true},
			{"f", false},
			{"false", false},
			{"off", false},
			{"0", false},
		}
		for _, scenario := range scenarios {
			t.Run(scenario.Param, func(t *testing.T) {
				formParams := make(url.Values)
				formParams.Set("active", scenario.Param)

				req := httptest.NewRequest("POST", "https://www.example.com/", strings.NewReader(formParams.Encode()))
				req.Header.Set("Content-type", "application/x-www-form-urlencoded")

				s := struct {
					Active bool `req:"active"`
				}{Active: !scenario.Expected}

				err := reqbind.Bind(&s, req)
				assert.NoError(t, err, nil)
				assert.Equal(t, scenario.Expected, s.Active)
			})
		}
	})

	t.Run("should zero bool if field is not present", func(t *testing.T) {
		formParams := make(url.Values)

		req := httptest.NewRequest("POST", "https://www.example.com/", strings.NewReader(formParams.Encode()))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := struct {
			Active bool `req:"active,zero"`
		}{Active: true}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.False(t, s.Active)
	})

	t.Run("should bind to fields from of nested structures", func(t *testing.T) {
		formParams := make(url.Values)
		formParams.Set("rules.title", "Rule title")
		formParams.Set("rules.category", "favourites")
		formParams.Set("action.active", "on")
		formParams.Set("action.remarks", "This is the remarks")

		req := httptest.NewRequest("POST", "https://www.example.com/", strings.NewReader(formParams.Encode()))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := struct {
			Rules struct {
				Name     string `req:"title"`
				Category string `req:"category"`
			} `req:"rules"`
			Action struct {
				Active  bool   `req:"active"`
				Remarks string `req:"remarks"`
			} `req:"action"`
		}{}

		err := reqbind.Bind(&s, req)
		assert.NoError(t, err, nil)
		assert.Equal(t, "Rule title", s.Rules.Name)
		assert.Equal(t, "favourites", s.Rules.Category)
		assert.Equal(t, true, s.Action.Active)
		assert.Equal(t, "This is the remarks", s.Action.Remarks)
	})
}
