package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
		fmt.Println(response.Body.String())
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe",                     // city не указан
		"/cafe?city=unknown",        // неизвестный город
		"/cafe?city=moscow&count=x", // некорректный count
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusBadRequest, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, tc := range tests {
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", tc.count)

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())
		var items []string
		if body != "" {
			items = strings.Split(body, ",")
		}

		assert.Equal(t, tc.want, len(items), "count=%d", tc.count)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, tc := range tests {
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", tc.search)

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())
		var items []string
		if body != "" {
			items = strings.Split(body, ",")
		}

		assert.Equal(t, tc.wantCount, len(items), "search=%s", tc.search)

		for _, name := range items {
			assert.Contains(t,
				strings.ToLower(name),
				strings.ToLower(tc.search),
			)
		}
	}
}
