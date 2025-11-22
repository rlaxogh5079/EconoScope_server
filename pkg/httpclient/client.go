package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rlaxogh5079/EconoScope/config"
	"github.com/rlaxogh5079/EconoScope/pkg/logger"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers http.Header
}

// 옵션 패턴
type Option func(*HTTPClient)

func WithTimeout(d time.Duration) Option {
	return func(h *HTTPClient) {
		h.client.Timeout = d
	}
}

func WithHeader(key, value string) Option {
	return func(h *HTTPClient) {
		h.headers.Set(key, value)
	}
}

func WithHeaders(hs http.Header) Option {
	return func(h *HTTPClient) {
		for k, vv := range hs {
			for _, v := range vv {
				h.headers.Add(k, v)
			}
		}
	}
}

// config 기반으로 외부 API 클라이언트 만들 때 사용 가능
func New(baseURL string, opts ...Option) *HTTPClient {
	c := &HTTPClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
		headers: make(http.Header),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GET + JSON decode
func (h *HTTPClient) GetJSON(ctx context.Context, path string, query string, out interface{}) error {
	url := h.baseURL + path
	if query != "" {
		url = fmt.Sprintf("%s?%s", url, query)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	for k, vv := range h.headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	logger.Log.WithField("url", url).Debug("HTTP GET")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Log.WithFields(map[string]interface{}{
			"url":    url,
			"status": resp.StatusCode,
			"body":   string(body),
		}).Warn("external API non-2xx")
		return fmt.Errorf("external api error: status=%d", resp.StatusCode)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return err
	}

	return nil
}

// POST JSON + JSON decode
func (h *HTTPClient) PostJSON(ctx context.Context, path string, reqBody interface{}, out interface{}) error {
	url := h.baseURL + path

	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	for k, vv := range h.headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	logger.Log.WithField("url", url).Debug("HTTP POST")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Log.WithFields(map[string]interface{}{
			"url":    url,
			"status": resp.StatusCode,
			"body":   string(body),
		}).Warn("external API non-2xx")
		return fmt.Errorf("external api error: status=%d", resp.StatusCode)
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return err
		}
	}

	return nil
}

func NewNewsAPIClient() *HTTPClient {
	cfg := config.AppConfig
	client := New(
		cfg.ExternalAPI.NewsAPI.BaseURL,
		WithTimeout(5*time.Second),
	)
	client.headers.Set("X-Api-Key", cfg.ExternalAPI.NewsAPI.APIKey)
	return client
}