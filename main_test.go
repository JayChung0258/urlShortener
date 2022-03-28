package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert"
)

// test the sql func , assuming http request is OK
func TestGetURLs(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/urls", nil)
	w := httptest.NewRecorder()

	//Hack to try to fake gorilla/mux vars
	// vars := map[string]string{
	// 	"id": "1",
	// }

	// CHANGE THIS LINE!!!
	// r = mux.SetURLVars(r, vars)

	getURLs(w, r)

	fmt.Println(w.Code)
	assert.Equal(t, http.StatusOK, w.Code)
}
