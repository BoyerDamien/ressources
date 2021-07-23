package media

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/BoyerDamien/gapi"

	"github.com/BoyerDamien/gapi/database/driver/sqlite"
)

type testApi struct {
	App *gapi.App
}

func (s *testApi) ReadData(resp *http.Response, result interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, result); err != nil {
		return err
	}
	return nil
}

func (s *testApi) Create(endpoint, path string, result interface{}) (*http.Response, error) {
	file, _ := os.Open(path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("media", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", endpoint, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := s.App.Test(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, nil
	}
	if err = s.ReadData(resp, result); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *testApi) Update(endpoint string, data interface{}, result interface{}) (*http.Response, error) {
	reqByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest("PUT", endpoint, bytes.NewReader(reqByte))
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.App.Test(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, nil
	}
	if err = s.ReadData(resp, result); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *testApi) Retrieve(endpoint string, result interface{}) (*http.Response, error) {
	req := httptest.NewRequest("GET", endpoint, nil)

	resp, err := s.App.Test(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, nil
	}
	if err = s.ReadData(resp, result); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *testApi) Delete(endpoint string, result interface{}) (*http.Response, error) {
	req := httptest.NewRequest("DELETE", endpoint, nil)

	resp, err := s.App.Test(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, nil
	}
	if err = s.ReadData(resp, result); err != nil {
		return nil, err
	}
	return resp, nil
}

func SetupApp(url string) *gapi.App {
	os.Remove("test.db")
	app := gapi.New(sqlite.Open("test.db"), gapi.Config{})
	api := app.Collection(url)
	api.AddRessources(&Media{})
	app.Static("/static", ".")
	return app
}

func ModelToString(m1 interface{}) string {
	r1, err := json.Marshal(m1)
	if err != nil {
		panic(err)
	}
	return string(r1)
}
