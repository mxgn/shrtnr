// handlers_test.go
package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	data := []byte(`{"url":"http://google.com"}`)
	// t.Error(string(data))
	r := bytes.NewReader(data)
	resp, err := http.Post("http://localhost/enc/", "application/json", r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	}
	t.Error("\nRESP DUMP: + \n" + string(respDump))

}
