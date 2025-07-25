package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// ______________________________ TestCafeCount _______________________________//

func TestCafeCount(t *testing.T) {
	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	hendler := http.HandlerFunc(mainHandle)

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("Get", fmt.Sprint("/cafe?city=moscow&count="+strconv.Itoa(v.count)), nil)

		hendler.ServeHTTP(response, req)

		length := 0

		body := response.Body.String()
		res := strings.Split(body, ",")

		length = len(res)
		if body == "" {
			length = 0
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, length, v.want, fmt.Sprint("The number of cafes in the request: ", v.want, ". Number of cafes received: ", len(res)))
	}
}

// ______________________________ TestCafeSearch _______________________________//

func TestCafeSearch(t *testing.T){
	requests := []struct {
		search string
		wantCount  int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	hendler := http.HandlerFunc(mainHandle)

	for _, v := range requests {
		count := 0
		response := httptest.NewRecorder()
		req := httptest.NewRequest("Get", fmt.Sprint("/cafe?city=moscow&search=", v.search), nil)

		hendler.ServeHTTP(response, req)

		data := strings.Split(response.Body.String(), ",")
		for _, w := range data {
			if strings.Contains(strings.ToLower(w), strings.ToLower(v.search)) {
				count += 1
			}
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, count, v.wantCount, fmt.Sprint("Request: ", v.search))
	}
}