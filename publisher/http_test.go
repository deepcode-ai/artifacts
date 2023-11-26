package publisher

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPPublisher(t *testing.T) {
	type args struct {
		opts       *HTTPOpts
		httpClient *http.Client
	}
	tests := []struct {
		name string
		args args
		want *HTTPPublisher
	}{
		{
			name: "with default opts",
			args: args{
				opts: &HTTPOpts{
					URL:     "http://localhost:8080",
					Timeout: 20 * time.Second,
				},
				httpClient: nil,
			},
			want: &HTTPPublisher{
				URL:         "http://localhost:8080",
				MaskingMode: MaskingModeNone,
				HTTPClient: &http.Client{
					Timeout:   20 * time.Second,
					Transport: &http.Transport{},
				},
			},
		},
		{
			name: "with custom opts",
			args: args{
				opts: &HTTPOpts{
					URL:         "http://localhost:8080",
					MaskingMode: MaskingModeEncrypted,
					Timeout:     20 * time.Second,
				},
			},
			want: &HTTPPublisher{
				URL:         "http://localhost:8080",
				MaskingMode: MaskingModeEncrypted,
				HTTPClient: &http.Client{
					Timeout:   20 * time.Second,
					Transport: &http.Transport{},
				},
			},
		},
		{
			name: "with no timeout set",
			args: args{
				opts: &HTTPOpts{
					URL:         "http://localhost:8080",
					MaskingMode: MaskingModeEncrypted,
				},
			},
			want: &HTTPPublisher{
				URL:         "http://localhost:8080",
				MaskingMode: MaskingModeEncrypted,
				HTTPClient: &http.Client{
					Timeout:   DefaultTimeout,
					Transport: &http.Transport{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NewHTTPPublisher(tt.args.opts).(*HTTPPublisher)
			if !ok {
				t.Errorf("NewHTTPPublisher() expected object to be *HTTPPublisher got = %v, want %v", got, tt.want)
			}
			if got.URL != tt.want.URL {
				t.Errorf("NewHTTPPublisher() URL = %v, want %v", got.URL, tt.want.URL)
			}
			if got.MaskingMode != tt.want.MaskingMode {
				t.Errorf("NewHTTPPublisher() MaskingMode = %v, want %v", got.MaskingMode, tt.want.MaskingMode)
			}
			if got.HTTPClient.Timeout != tt.want.HTTPClient.Timeout {
				t.Errorf("NewHTTPPublisher() HTTPClient.Timeout = %v, want %v", got.HTTPClient.Timeout, tt.want.HTTPClient.Timeout)
			}
		})
	}
}

func TestHTTPPublisher_Publish(t *testing.T) {
	type fields struct {
		MaskingMode string
		URL         string
		Token       string
		HTTPClient  *http.Client
	}
	type args struct {
		ctx     context.Context
		payload Payload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		server  *httptest.Server
	}{
		{
			name: "HTTPPublisher with server success",
			fields: fields{
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "HTTPPublisher with token",
			fields: fields{
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
				Token:       "test-token",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token := r.Header.Get("Authorization")
				parts := strings.Split(token, " ")
				if len(parts) != 2 {
					t.Error("HTTPPublisher with token: token is not in the correct format")
				}
				if parts[1] != "test-token" {
					t.Error("HTTPPublisher with token: token is not correct")
				}
				w.WriteHeader(http.StatusOK)
			})),
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "HTTPPublisher with server error",
			fields: fields{
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			wantErr: true,
		},
		{
			name: "HTTPPublisher with server error",
			fields: fields{
				URL:         "http://xxx",
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			wantErr: true,
		},
		{
			name: "compress error",
			fields: fields{
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					err: fmt.Errorf("compress error"),
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			wantErr: true,
		},
		{
			name: "request error",
			fields: fields{
				MaskingMode: MaskingModeNone,
				HTTPClient:  http.DefaultClient,
				URL:         string([]byte{0x7f}),
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.URL == "" {
				tt.fields.URL = tt.server.URL
			}
			h := &HTTPPublisher{
				URL:         tt.fields.URL,
				MaskingMode: tt.fields.MaskingMode,
				Token:       tt.fields.Token,
				HTTPClient:  tt.fields.HTTPClient,
			}
			if err := h.Publish(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("HTTPPublisher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
