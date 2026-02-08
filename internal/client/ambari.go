// Package client provides the Ambari API HTTP client with connection pooling and retries
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

// AmbariClient interface for Ambari API operations
type AmbariClient interface {
	Get(ctx context.Context, path string, params map[string]string) (map[string]interface{}, error)
	Post(ctx context.Context, path string, params map[string]string, body interface{}) (map[string]interface{}, error)
	Put(ctx context.Context, path string, params map[string]string, body interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, path string, params map[string]string) (map[string]interface{}, error)
}

// Config for the Ambari client
type Config struct {
	BaseURL  string
	Username string
	Password string
	Timeout  time.Duration
	Retries  int
}

type ambariClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	retries    int
	logger     *logrus.Logger
}

// NewAmbariClient creates a new Ambari HTTP client with connection pooling
func NewAmbariClient(cfg Config, logger *logrus.Logger) AmbariClient {
	return &ambariClient{
		baseURL:  cfg.BaseURL,
		username: cfg.Username,
		password: cfg.Password,
		retries:  cfg.Retries,
		logger:   logger,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *ambariClient) Get(ctx context.Context, path string, params map[string]string) (map[string]interface{}, error) {
	return c.doRequest(ctx, "GET", path, params, nil)
}

func (c *ambariClient) Post(ctx context.Context, path string, params map[string]string, body interface{}) (map[string]interface{}, error) {
	return c.doRequest(ctx, "POST", path, params, body)
}

func (c *ambariClient) Put(ctx context.Context, path string, params map[string]string, body interface{}) (map[string]interface{}, error) {
	return c.doRequest(ctx, "PUT", path, params, body)
}

func (c *ambariClient) Delete(ctx context.Context, path string, params map[string]string) (map[string]interface{}, error) {
	return c.doRequest(ctx, "DELETE", path, params, nil)
}

func (c *ambariClient) doRequest(ctx context.Context, method, path string, params map[string]string, body interface{}) (map[string]interface{}, error) {
	var lastErr error
	for attempt := 0; attempt <= c.retries; attempt++ {
		result, err := c.execute(ctx, method, path, params, body)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if attempt < c.retries {
			backoff := time.Duration(attempt+1) * 100 * time.Millisecond
			c.logger.WithFields(logrus.Fields{"attempt": attempt + 1, "method": method, "path": path}).Warn("Retrying")
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retries+1, lastErr)
}

func (c *ambariClient) execute(ctx context.Context, method, path string, params map[string]string, body interface{}) (map[string]interface{}, error) {
	reqURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if len(params) > 0 {
		q := reqURL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-By", "mcp-ambari")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	dur := time.Since(start)
	if err != nil {
		c.logger.WithFields(logrus.Fields{"method": method, "url": reqURL.String(), "duration": dur}).Error("Request failed")
		return nil, fmt.Errorf("HTTP %s %s failed: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result map[string]interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			result = map[string]interface{}{"raw": string(respBody)}
		}
	}

	c.logger.WithFields(logrus.Fields{"method": method, "path": path, "status": resp.StatusCode, "duration": dur}).Debug("Request done")

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP %d from %s %s", resp.StatusCode, method, path)
	}
	return result, nil
}
