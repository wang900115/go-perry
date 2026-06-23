package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/navidrome/navidrome/log"
)

const cacheSizeLimit = 1000

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	cache SimpleCache[string, string]
	hc    httpDoer

	ttl time.Duration
}

type requestData struct {
	Method string
	Header http.Header
	URL    string
	Body   *string
}

func NewHTTPClient(wrapped httpDoer, ttl time.Duration) *HTTPClient {
	c := &HTTPClient{hc: wrapped, ttl: ttl}
	c.cache = NewSimpleCache[string, string](Options{
		SizeLimit:  cacheSizeLimit,
		DefaultTTL: ttl,
	})
	return c
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	key := c.serializeReq(req)
	cached := true
	start := time.Now()
	respStr, err := c.cache.GetWithLoader(key, func(key string) (string, time.Duration, error) {
		cached = false
		req, err := c.deserializeReq(key)
		if err != nil {
			log.Trace(req.Context(), "CachedHTTPClient.Do", "key", key, err)
			return "", 0, err
		}
		resp, err := c.hc.Do(req)
		if err != nil {
			log.Trace(req.Context(), "CachedHTTPClient.Do", "req", req, err)
			return "", 0, err
		}
		defer resp.Body.Close()
		return c.serializeResp(resp), c.ttl, nil
	})
	log.Trace(req.Context(), "CachedHTTPClient.Do", "key", key, "cached", cached, "duration", time.Since(start))
	if err != nil {
		return nil, err
	}
	return c.deserializeResp(respStr)
}

func (c *HTTPClient) serializeReq(req *http.Request) string {
	// Serialize the request data to create a unique cache key
	// You can customize this serialization based on your needs
	data := requestData{
		Method: req.Method,
		Header: req.Header,
		URL:    req.URL.String(),
	}
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		data.Body = new(base64.StdEncoding.EncodeToString(bodyBytes))
	}
	j, _ := json.Marshal(&data)
	return string(j)
}

func (c *HTTPClient) deserializeReq(reqStr string) (*http.Request, error) {
	var data requestData
	_ = json.Unmarshal([]byte(reqStr), &data)
	var body io.Reader
	if data.Body != nil {
		bodyStr, _ := base64.StdEncoding.DecodeString(*data.Body)
		body = strings.NewReader(string(bodyStr))
	}
	req, err := http.NewRequest(data.Method, data.URL, body)
	if err != nil {
		return nil, err
	}
	req.Header = data.Header
	return req, nil
}

func (c *HTTPClient) serializeResp(resp *http.Response) string {
	var b = &bytes.Buffer{}
	_ = resp.Write(b)
	return b.String()
}

func (c *HTTPClient) deserializeResp(respStr string) (*http.Response, error) {
	r := bufio.NewReader(strings.NewReader(respStr))
	return http.ReadResponse(r, nil)
}
