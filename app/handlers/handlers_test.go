// handlers_test.go
package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	data := url.Values{"url": {"http://google.com"}}
	resp, err := http.PostForm("http://localhost/enc/", data)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	}
	t.Error("\nRESP DUMP !!: + \n" + string(respDump))

}
