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

/**
* TestCafeWhenOk tests the cafe handler with valid requests.
* It checks that the handler returns HTTP status OK (200) for various valid requests,
* including requests with count parameter, city parameter, and search parameter.
*
* @param t *testing.T - testing object for assertions and logging
*/
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

/**
* TestCafeCount tests the cafe count functionality of the mainHandle handler.
* It sends various count values and verifies the number of cafes returned matches expectations.
*
* @param t *testing.T The testing object provided by the testing package
*/
func TestCafeCount(t *testing.T) {
    handler := http.HandlerFunc(mainHandle)

    requests := []struct {
        count int // передаваемое значение count
        want  int // ожидаемое количество кафе в ответе
    }{
        {0, 0},
        {1, 1},
        {2, 2},
        {100, len(cafeList["moscow"])}, // Москва имеет 5 кафе
    }

    for _, v := range requests {
        query := "/cafe?city=moscow&count=" + strconv.Itoa(v.count)

        response := httptest.NewRecorder()
        req := httptest.NewRequest("GET", query, nil)
        handler.ServeHTTP(response, req)

        assert.Equal(t, http.StatusOK, response.Code)

        body := response.Body.String()
        cafes := []string{}
        if body != "" {
            cafes = strings.Split(body, ",")
        }

        if v.want == 0 {
            assert.Empty(t, cafes, "0 cafes were expected")
        } else {
            assert.Len(t, cafes, v.want, "The number of cafes does not match the expected number")
        }
    }
}

func TestCafeSearch(t *testing.T) {
    handler := http.HandlerFunc(mainHandle)

    requests := []struct {
        search    string
        wantCount int
    }{
        {"фасоль", 0},
        {"кофе", 2},
        {"вилка", 1},
    }

    for _, v := range requests {
        query := "/cafe?city=moscow&search=" + v.search

        response := httptest.NewRecorder()
        req := httptest.NewRequest("GET", query, nil)
        handler.ServeHTTP(response, req)

        assert.Equal(t, http.StatusOK, response.Code)

        body := response.Body.String()
        var cafes []string
        if body != "" {
            cafes = strings.Split(body, ",")
        }

        assert.Len(t, cafes, v.wantCount, "The number of cafes found does not match")

        searchLower := strings.ToLower(v.search)
        for _, cafe := range cafes {
            assert.True(t, strings.Contains(strings.ToLower(cafe), searchLower),
                "Cafe '%s' does not contain the substring '%s'", cafe, v.search)
        }
    }
}