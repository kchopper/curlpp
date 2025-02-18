package client

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kchopper/curlpp/internal/config"
	"github.com/kchopper/curlpp/pkg/auth"
)

type Client struct {
	httpClient *http.Client
	config     *config.Config
}

// Request represents an HTTP request with additional options
type Request struct {
	URL      string
	Method   string
	Pretty   bool
	Parallel int
	Retries  int
	Headers  map[string]string
	Body     []byte
}

// Response represents an HTTP response with timing information
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Timing     *TimingInfo
}

type TimingInfo struct {
	DNSLookup     time.Duration
	TLSHandshake  time.Duration
	ServerTime    time.Duration
	TotalDuration time.Duration
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
			},
		},
		config: cfg,
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Apply authentication if configured
	if c.config.Current != "" {
		profile := c.config.Profiles[c.config.Current]
		authenticator, err := auth.GetAuthenticator(profile.Auth)
		if err != nil {
			return nil, fmt.Errorf("failed to get authenticator: %w", err)
		}

		if err := authenticator.ApplyAuth(httpReq); err != nil {
			return nil, fmt.Errorf("failed to apply authentication: %w", err)
		}
	}

	// Add custom headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
		Timing: &TimingInfo{
			TotalDuration: time.Since(start),
			// Note: For detailed timing info, we'll need to implement custom transport
		},
	}, nil
}
