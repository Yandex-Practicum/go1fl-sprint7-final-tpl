package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"strconv"

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
        count int   // передаваемое значение count
        want  int   // ожидаемое количество кафе в ответе
    }{
        {0, 0},
		{1, 1},
		{2, 2},
		{100, min(len(cafeList["moscow"]), 100)},
    }
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(v.count), nil)

		handler.ServeHTTP(response, req)

		responseArr := []string{}
		if response.Body.String() != "" {
			responseArr = strings.Split(response.Body.String(), ",")
		}

		assert.Equal(t, v.want, len(responseArr))
	} 
}

func TestCafeSearch (t *testing.T) {
    handler := http.HandlerFunc(mainHandle)
	
	requests := []struct {
        search    string   // передаваемое значение search 
        wantCount int      // ожидаемое количество кафе в ответе
    }{
        {"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
    }

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)

		handler.ServeHTTP(response, req)

		responseArr := []string{}
		if response.Body.String() != "" {
			responseArr = strings.Split(response.Body.String(), ",")
			for _, s := range responseArr {
				assert.Contains(t, strings.ToLower(s), strings.ToLower(v.search))
			}
		}

		assert.Equal(t, v.wantCount, len(responseArr))
	}
}