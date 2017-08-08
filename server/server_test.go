package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorWriter struct{}

func (e *ErrorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("Expected error")
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	buf := testTemplate(indexTemplate, nil, t)

	if rr.Body.String() != buf.String() {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), buf.String())
	}

}

func TestBunReport(t *testing.T) {
	//test GET
	req, err := http.NewRequest("GET", "/bunReport", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bunReport)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	buf := testTemplate(reportTemplate, nil, t)
	if rr.Body.String() != buf.String() {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), buf.String())
	}

	//test POST
	req, err = http.NewRequest("POST", "/bunReport?size=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	req, err = http.NewRequest("POST", "/bunReport", nil)
	if err != nil {
		t.Fatal("Unable to handle bunReport with no size")
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
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
