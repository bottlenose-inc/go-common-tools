package testhttp

import (
    "net/http"
    "net/http/httptest"
    "net/url"
)

type TestHTTPResponse struct {
	Status int
	Body   []byte
}

type MockHTTP struct {
	Server    *httptest.Server
    Client    http.Client

	Responses map[string]TestHTTPResponse
}

func InitMockHTTP() *MockHTTP {
    var mock MockHTTP

    mock.Responses = make(map[string]TestHTTPResponse)
    mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rUrl := r.URL
		response, found := mock.Responses[rUrl.String()]

		if found {
			w.WriteHeader(response.Status)
			w.Header().Set("Content-Type", "application/json")
            w.Write(response.Body)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(""))
		}
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(mock.Server.URL)
		},
	}

	mock.Client = http.Client{Transport: transport}

    return &mock
}

func (mock *MockHTTP) AddTestData(testUrl string, code int, body []byte) {
	var resp TestHTTPResponse
	resp.Status = code
	resp.Body = body
	mock.Responses[testUrl] = resp
}

func (mock *MockHTTP) DeleteTestData(testUrl string) {
    delete(mock.Responses, testUrl)
}

func (mock *MockHTTP) Close() {
    mock.Server.Close()
}
