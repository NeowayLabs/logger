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

	if firstNamespace.Level != 3 || secondNamespace.Level != 3 {
		t.Fatal("Level should be Info, but got", firstNamespace.Level, secondNamespace.Level)
	}

	var jsonStr = []byte(`{"level":"debug"}`)
	url := "http://testeurl.com/logger/all"

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.RequestURI = url
	w := httptest.NewRecorder()

	HTTPFunc(w, req)

	if firstNamespace.Level != 4 || secondNamespace.Level != 4 {
		t.Fatal("Level should be", string(jsonStr), ", but got", firstNamespace.Level, secondNamespace.Level)
	}
}
