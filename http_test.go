package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldChangeAllNamespacesLevel(t *testing.T) {
	firstNamespace := Namespace("control")
	secondNamespace := Namespace("module")
	var jsonStr = []byte(`{"level":"debug"}`)
	url := "http://testeurl.com/logger/all"

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.RequestURI = url
	w := httptest.NewRecorder()

	HTTPFunc(w, req)

	if firstNamespace.Level != 3 || secondNamespace.Level != 3 {
		t.Fatal("Level should be", jsonStr, "But got", firstNamespace.Level, secondNamespace.Level)
	}
}
