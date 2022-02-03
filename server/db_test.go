package server

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"testing"

	boltdb "github.com/ric-v/golang-key-value-db-browser/bolt-db"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

// serve serves http request using provided fasthttp handler
func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {

	// new in memory listener for testing apis
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	// start the server
	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	// create a new client
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	// make the request
	return client.Do(req)
}

func multiFormFileGen(fileName string) (*multipart.Writer, *bytes.Buffer, error) {

	var body = new(bytes.Buffer)
	var writer = multipart.NewWriter(body)

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return writer, body, err
	}

	formFile, err := writer.CreateFormFile("boltdb", fileName)
	if err != nil {
		fmt.Println(err)
		return writer, body, err
	}

	_, err = formFile.Write(fileContent)
	if err != nil {
		fmt.Println(err)
		return writer, body, err
	}

	return writer, body, err
}

func Test_storeDBFile(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		wantStatusCode int
		wantErr        bool
	}{
		{
			name:           "valid file",
			fileName:       "test.db",
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid file",
			fileName:       "wrong.db",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	testDB, _ := boltdb.New("test.db")
	defer testDB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			writer, body, err := multiFormFileGen(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
			writer.Close()

			req, err := http.NewRequest("POST", "http://api/v1/db/upload", body)
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Content-Type", writer.FormDataContentType())

			// set query params
			q := req.URL.Query()
			q.Add("dbtype", "boltdb")
			req.URL.RawQuery = q.Encode()

			// create a new client
			res, _ := serve(uploadFile, req)
			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("got http code = %v, want http code %v", res.StatusCode, tt.wantStatusCode)
			}
		})
	}
	os.RemoveAll("test.db")
	os.RemoveAll("temp")
}

func Test_createNewDBFile(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create new db file",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newFile(tt.args.ctx)
		})
	}
}
