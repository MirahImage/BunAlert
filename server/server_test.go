package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
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
		{"GET", "/bunReport", nil, testTemplate(reportTemplate, fmt.Sprintf("%x", generateToken()), t).String(), http.StatusOK},
		{"POST", "/bunReport", nil, "", http.StatusBadRequest},
		{"POST", "/bunReport?token=ffffffffffffffff0000000000000000", nil, "", http.StatusOK},
		{"POST", "/bunReport?size=1&token=ffffffffffffffff0000000000000000", nil, "", http.StatusOK},
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

	token := fmt.Sprintf("%x", generateToken())
	err = writeTemplate(rr, reportTemplate, token)
	if err != nil {
		t.Fatal("Error with template file:", err)
	}

	buf := testTemplate(reportTemplate, token, t)

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

func TestValidateToken(t *testing.T) {
	f0, _ := hex.DecodeString("ffffffffffffffff0000000000000000")
	var tokenTests = []struct {
		tokenString string
		token       []byte
		err         error
	}{
		{"", nil, fmt.Errorf("Token length invalid: was 0 want %v", md5.Size)},
		{"ztx", nil, fmt.Errorf("Expected error")},
		{"ffffffffffffffff", nil, fmt.Errorf("Token length invalid: was 16 want %v", md5.Size)},
		{"ffffffffffffffff0000000000000000", f0, nil},
	}

	for _, tt := range tokenTests {
		token, err := validateToken(tt.tokenString)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error %v", err)
		}
		if err == nil && tt.err != nil {
			t.Errorf("Expected error %v got %v", tt.err, err)
		}
		if !reflect.DeepEqual(token, tt.token) {
			t.Errorf("Expected token %v got %v", tt.token, token)
		}
	}
}

func TestGenerateToken(t *testing.T) {
	sum := generateToken()
	if sum == nil {
		t.Fatal("Generated empty token")
	}
	if len(sum) != md5.Size {
		t.Errorf("Token is of wrong length, got %v want %v",
			len(sum), md5.Size)
	}
	time.Sleep(time.Second)
	sum2 := generateToken()
	if reflect.DeepEqual(sum, sum2) {
		t.Errorf("Tokens are not unique: %x and %x", sum, sum2)
	}
}
