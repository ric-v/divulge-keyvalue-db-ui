package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// check if required directories exist
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		log.Println("creating ./temp directory")
		os.Mkdir("temp", 0755)
	}
}

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

	formFile, err := writer.CreateFormFile("file", fileName)
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

func Test_uploadFile(t *testing.T) {
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
			wantErr:        false,
		},
	}

	testDB, _ := database.NewDB("test.db", database.BUNT_DB)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			writer, body, _ := multiFormFileGen(tt.fileName)
			writer.Close()

			req, err := http.NewRequest("POST", "http://api/v1/upload", body)
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
			res, err := serve(uploadFile, req)
			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("got http code = %v, want http code %v", res.StatusCode, tt.wantStatusCode)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, want error %v", err, tt.wantErr)
				return
			}
			fmt.Println(res, err)
			var apiResp apiResponse
			json.NewDecoder(res.Body).Decode(&apiResp)
			db, ok := session.Load(apiResp.DBKey)
			if ok {
				db.(Session).DB.CloseDB()
				session.Delete(apiResp.DBKey)
			}
		})
	}
	testDB.CloseDB()
	os.RemoveAll("test.db")
}

func Test_newFile(t *testing.T) {
	type args struct {
		dbtype string
		file   string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "bolt db 200",
			args: args{
				dbtype: database.BOLT_DB,
				file:   "test.db",
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "bolt db 500",
			args: args{
				dbtype: database.BOLT_DB,
				file:   "/unknown/folder/test.db",
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        false,
		},
		{
			name: "bunt db 200",
			args: args{
				dbtype: database.BUNT_DB,
				file:   "test.db",
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "bunt db 500",
			args: args{
				dbtype: database.BUNT_DB,
				file:   "/unknown/folder/test.db",
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest("POST", "http://api/v1/new", nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			// set query params
			q := req.URL.Query()
			q.Add("dbtype", tt.args.dbtype)
			q.Add("file", tt.args.file)
			req.URL.RawQuery = q.Encode()

			// create a new client
			res, err := serve(newFile, req)
			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("got http code = %v, want http code %v", res.StatusCode, tt.wantStatusCode)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, want error %v", err, tt.wantErr)
				return
			}
			fmt.Println(res, err)
			var apiResp apiResponse
			json.NewDecoder(res.Body).Decode(&apiResp)
			db, ok := session.Load(apiResp.DBKey)
			if ok {
				db.(Session).DB.CloseDB()
				session.Delete(apiResp.DBKey)
			}
		})
	}
}

func Test_listKeyValue(t *testing.T) {
	type args struct {
		accesskey string
		dbtype    string
		bucket    string
		file      string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "bolt db",
			args: args{
				accesskey: "test",
				dbtype:    database.BOLT_DB,
				bucket:    "test-bucket",
				file:      "test.db",
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        false,
		},
		{
			name: "bolt db - invalid accesskey",
			args: args{
				accesskey: "test-invalid",
				dbtype:    database.BOLT_DB,
				bucket:    "test-bucket",
				file:      "test.db",
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        false,
		},
		{
			name: "bunt db",
			args: args{
				accesskey: "test",
				dbtype:    database.BUNT_DB,
				file:      "test.db",
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "bunt db - invalid accesskey",
			args: args{
				accesskey: "test-invalid",
				dbtype:    database.BUNT_DB,
				file:      "test.db",
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testDB, _ := database.NewDB("test.db", tt.args.dbtype)
			testDB.Add("test-key", "test-value")
			defer testDB.CloseDB()

			session.Store("test", Session{tt.args.accesskey, tt.args.file, tt.args.file, testDB})

			req, err := http.NewRequest("POST", "http://api/v1/db", nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			// set query params
			q := req.URL.Query()
			q.Add("accesskey", tt.args.accesskey)
			q.Add("dbtype", tt.args.dbtype)
			q.Add("bucket", tt.args.bucket)
			q.Add("file", tt.args.file)
			req.URL.RawQuery = q.Encode()

			// create a new client
			res, err := serve(listKeyValue, req)
			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("got http code = %v, want http code %v", res.StatusCode, tt.wantStatusCode)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, want error %v", err, tt.wantErr)
				return
			}
			fmt.Println(res, err)
			var apiResp apiResponse
			json.NewDecoder(res.Body).Decode(&apiResp)
			db, ok := session.Load(apiResp.DBKey)
			if ok {
				db.(Session).DB.CloseDB()
				session.Delete(apiResp.DBKey)
			}
		})
	}
	os.Remove("test.db")
}
