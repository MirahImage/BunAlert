package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	buf := testTemplate(indexTemplate, nil, t)

	if rr.Body.String() != buf.String() {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), buf.String())
	}

}

func TestBunReport(t *testing.T) {
	var handlerTests = []struct {
		method      string
		urlString   string
		requestBody io.Reader
		body        string
		status      int
	}{
		{"GET", "/bunReport", nil, testTemplate(reportTemplate, nil, t).String(), http.StatusOK},
		{"POST", "/bunReport", nil, "", http.StatusOK},
		{"POST", "/bunReport?size=1", nil, "", http.StatusOK},
	}

	for _, tt := range handlerTests {
		req, err := http.NewRequest(tt.method, tt.urlString, tt.requestBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bunReport)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.status {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, tt.status)
		}

		if tt.body != "" && rr.Body.String() != tt.body {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), tt.body)
		}
	}
}

func TestWriteTemplate(t *testing.T) {
	rr := httptest.NewRecorder()
	err := writeTemplate(rr, "", nil)
	if err == nil {
		t.Fatal("Did not receive expected error no such file or directory:", err)
	}

	err = writeTemplate(rr, indexTemplate, nil)
	if err != nil {
		t.Fatal("Error with template file:", err)
	}

	buf := testTemplate(indexTemplate, nil, t)

	if rr.Body.String() != buf.String() {
		t.Errorf("writeTemplate returned unexpected body: got %v expected %v",
			rr.Body.String(), buf.String())
	}
}

func testTemplate(templateFile string, templateData interface{}, t *testing.T) (buf *bytes.Buffer) {
	tpl, err := template.ParseFiles(templateFile)
	if err != nil {
		t.Fatal("Error with template file:", err)
	}

	buf = new(bytes.Buffer)
	err = tpl.Execute(buf, templateData)
	if err != nil {
		t.Fatal("Error executing template:", err)
	}

	return buf
}

func TestMain(t *testing.T) {
	_, err := http.NewRequest("GET", ":-41", nil)
	if err == nil {
		t.Fatal("No error for GET from port -41 ")
	}
	_, err = http.NewRequest("GET", "/fakeURL", nil)
	if err != nil {
		t.Fatal("Error for GET from fakeURL", err)
	}
}
