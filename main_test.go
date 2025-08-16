package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		count int
		want  int
		city  string
	}{
		{count: 0, want: 0, city: "moscow"},
		{count: 1, want: 1, city: "moscow"},
		{count: 2, want: 2, city: "moscow"},
		{count: 100, want: min(100, len(cafeList["moscow"])), city: "moscow"},
		{count: 3, want: 3, city: "tula"},
		{count: 5, want: min(5, len(cafeList["tula"])), city: "tula"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&count="+strconv.Itoa(v.count), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		actualCount := 0
		str := strings.TrimSpace(response.Body.String())
		if str != "" {
			actualCount = len(strings.Split(str, ","))
		}
		assert.Equal(t, v.want, actualCount)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
		city      string
	}{
		{search: "фасоль", wantCount: 0, city: "moscow"},
		{search: "кофе", wantCount: 2, city: "moscow"},
		{search: "вилка", wantCount: 1, city: "moscow"},
		{search: "мир", wantCount: 1, city: "tula"},
		{search: "кофе", wantCount: 0, city: "tula"},
		{search: "Завтрак", wantCount: 1, city: "tula"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&search="+v.search, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		actualCount := 0
		str := strings.TrimSpace(response.Body.String())
		if str != "" {
			results := strings.Split(str, ",")
			actualCount = len(results)
			for _, entry := range results {
				ok := strings.Contains(strings.ToLower(entry), strings.ToLower(v.search))
				assert.True(t, ok)
			}
		}
		assert.Equal(t, v.wantCount, actualCount)
	}
}
