package publisher

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	PublisherTypeHTTP string = "http"
)

const (
	MaskingModeNone      = "none"
	MaskingModeSimple    = "simple"
	MaskingModeEncrypted = "encrypted"

	DefaultTimeout = 20 * time.Second
)

type HTTPPublisher struct {
	URL         string
	MaskingMode string
	Token       string
	HTTPClient  *http.Client
}

type HTTPOpts struct {
	URL         string
	MaskingMode string
	Token       string
	Timeout     time.Duration
}

func NewHTTPPublisher(opts *HTTPOpts) Publisher {
	if opts.MaskingMode == "" {
		opts.MaskingMode = MaskingModeNone
	}
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}

	return &HTTPPublisher{
		URL:         opts.URL,
		MaskingMode: opts.MaskingMode,
		Token:       opts.Token,
		HTTPClient: &http.Client{
			Timeout: opts.Timeout,
		},
	}
}

// Publish sends the payload to the configured URL.
// TODO: Add backoff and retry logic.
func (h *HTTPPublisher) Publish(ctx context.Context, payload Payload) error {
	body, err := payload.Bytes()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	if h.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	}

	res, err := h.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
