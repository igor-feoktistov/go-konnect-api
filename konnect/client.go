package konnect

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	libraryVersion = "1.0.0"
	userAgent      = "go-konnect/" + libraryVersion
)

type Client struct {
	Client             *http.Client
	BaseURL            *url.URL
	UserAgent          string
	Options		   *ClientOptions
	ResponseTimeout	   time.Duration
}

type ClientOptions struct {
	Token     string
	SSLVerify bool
	Debug     bool
	Timeout   time.Duration
}

type ErrorResponse struct {
	Status   int    `json:"status,omitempty"`
	Title    string `json:"title,omitempty"`
	Instance string `json:"instance,omitempty"`
	Message  string `json:"message,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

type RestResponse struct {
	ErrorResponse ErrorResponse
	HttpResponse *http.Response
}

func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		SSLVerify: true,
		Debug:     false,
		Timeout:   60 * time.Second,
	}
}

func NewClient(endpoint string, options *ClientOptions) *Client {
	if options == nil {
		options = DefaultOptions()
	}
	httpClient := &http.Client {
		Timeout: options.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !options.SSLVerify,
			},
		},
	}
	if !strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint + "/"
	}
	baseURL, _ := url.Parse(endpoint)
	c := &Client{
		Client:          httpClient,
		BaseURL:         baseURL,
		UserAgent:       userAgent,
		Options:         options,
		ResponseTimeout: options.Timeout,
	}
	return c
}

func (c *Client) NewRequest(method string, apiPath string, parameters []string, body interface{}) (req *http.Request, err error) {
	var payload io.Reader
	var extendedPath string
	if len(parameters) > 0 {
		extendedPath = fmt.Sprintf("%s?%s", apiPath, strings.Join(parameters, "&"))
	} else {
		extendedPath = apiPath
	}
	u, _ := c.BaseURL.Parse(extendedPath)
	if body != nil {
		buf, err := json.MarshalIndent(body, "", "  ")
		if err != nil {
			return nil, err
		}
		if c.Options.Debug {
			log.Printf("[DEBUG] request JSON:\n%v\n\n", string(buf))
		}
		payload = bytes.NewBuffer(buf)
	}
	req, err = http.NewRequest(method, u.String(), payload)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Options.Token != "" {
		req.Header.Set("Authorization", "Bearer " + c.Options.Token)
	}
	if c.Options.Debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("[DEBUG] request dump:\n%q\n\n", dump)
	}
	return
}

func (c *Client) Do(req *http.Request, v interface{}) (resp *RestResponse, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), c.ResponseTimeout)
	defer cncl()
	resp, err = checkResp(c.Client.Do(req.WithContext(ctx)))
	if err != nil {
		return
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.HttpResponse.Body)
	if err != nil {
		return
	}
	resp.HttpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	if c.Options.Debug {
		log.Printf("[DEBUG] response JSON:\n%v\n\n", string(b))
	}
	if v != nil {
		defer resp.HttpResponse.Body.Close()
		err = json.NewDecoder(resp.HttpResponse.Body).Decode(v)
	}
	return
}

func checkResp(resp *http.Response, err error) (*RestResponse, error) {
	if err != nil {
		return &RestResponse{HttpResponse: resp}, err
	}
	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return &RestResponse{HttpResponse: resp}, err
	default:
		restResp, httpErr := newHTTPError(resp)
		return restResp, httpErr
	}
}

func newHTTPError(resp *http.Response) (restResp *RestResponse, err error) {
	errResponse := ErrorResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&errResponse); err == nil {
		defer resp.Body.Close()
		if len(errResponse.Detail) > 0{
			err = fmt.Errorf("Error: HTTP code=%d, HTTP status=\"%s\", REST status=\"%d\", REST title=\"%s\", REST detail=\"%s\"", resp.StatusCode, http.StatusText(resp.StatusCode), errResponse.Status, errResponse.Title, errResponse.Detail)
		} else {
			if len(errResponse.Message) > 0 {
				err = fmt.Errorf("Error: HTTP code=%d, HTTP status=\"%s\", REST message=\"%s\"", resp.StatusCode, http.StatusText(resp.StatusCode), errResponse.Message)
			} else {
				err = fmt.Errorf("Error: HTTP code=%d, HTTP status=\"%s\"", resp.StatusCode, http.StatusText(resp.StatusCode))
			}
		}
	} else {
		err = fmt.Errorf("Error: HTTP code=%d, HTTP status=\"%s\"", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	restResp = &RestResponse{
		ErrorResponse: errResponse,
		HttpResponse: resp,
	}
	return
}
