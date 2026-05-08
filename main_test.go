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

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

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

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: len(cafeList["moscow"])},
	}

	for _, tt := range requests {
		t.Run(fmt.Sprintf("count=%d", tt.count), func(t *testing.T) {
			url := fmt.Sprintf("/cafe?city=moscow&count=%d", tt.count)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			mainHandle(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			body := strings.TrimSpace(w.Body.String())

			if body == "" {
				assert.Equal(t, 0, tt.want)
				return
			}

			cafes := strings.Split(body, ",")

			assert.Equal(t, tt.want, len(cafes))
		})
	}
}

func TestCafeSearch(t *testing.T) {
	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, tt := range requests {
		t.Run(fmt.Sprintf("search=%s", tt.search), func(t *testing.T) {
			url := fmt.Sprintf("/cafe?city=moscow&search=%s", tt.search)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			mainHandle(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			body := strings.TrimSpace(w.Body.String())

			if body == "" {
				assert.Equal(t, 0, tt.wantCount)
				return
			}

			cafes := strings.Split(body, ",")

			assert.Equal(t, tt.wantCount, len(cafes))

			searchLower := strings.ToLower(tt.search)
			for _, cafe := range cafes {
				cafeTrimmed := strings.TrimSpace(cafe)
				cafeLower := strings.ToLower(cafeTrimmed)
				assert.True(t, strings.Contains(cafeLower, searchLower))
			}
		})
	}
}