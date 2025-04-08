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

		assert.Equal(t, http.StatusOK, response.Code)

		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

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

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	//я думаю, что такая структура будет более гибкая для тестов
	requests := []struct {
		request string
		status  int
		count   int
	}{
		{"/cafe?count=0&city=moscow", http.StatusOK, 0},
		{"/cafe?count=1&city=moscow", http.StatusOK, 1},
		{"/cafe?count=2&city=moscow", http.StatusOK, 2},
		{"/cafe?count=100&city=moscow", http.StatusOK, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, v.status, response.Code)
		if response.Body.String() == "" {
			assert.Equal(t, v.count, 0)
			continue
		}
		assert.Equal(t, v.count, len(strings.Split(response.Body.String(), ",")))

	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		request   string
		search    string
		status    int
		wantCount int
	}{
		{"/cafe?search=фасоль&city=moscow", "фасоль", http.StatusOK, 0},
		{"/cafe?search=кофе&city=moscow", "кофе", http.StatusOK, 2},
		{"/cafe?search=вилка&city=moscow", "вилка", http.StatusOK, 1},
	}

	for _, v := range requests {
		responce := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(responce, req)

		require.Equal(t, v.status, responce.Code)
		if responce.Body.String() == "" {
			assert.Equal(t, v.wantCount, 0)
			continue
		}
		cafe := strings.Split(responce.Body.String(), ",")

		count := 0
		for _, c := range cafe {
			if strings.Contains(strings.ToLower(c), strings.ToLower(v.search)) {
				count++
			}
		}
		assert.Equal(t, v.wantCount, count)

	}
}
