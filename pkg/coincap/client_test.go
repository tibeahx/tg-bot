package coincap

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

func TestGetAssetPrice(t *testing.T) {
	tests := []struct {
		name          string
		asset         models.Asset
		cfg           Config
		status        int
		responseBody  string
		expectedError string
	}{
		{
			name:  "Successful request",
			asset: models.Asset{Name: "bitcoin"},
			cfg: Config{
				Coincap: CoincapConfig{
					APIurl: "https://example.com",
					APIkey: "API_KEY",
				},
			},
			status:        http.StatusOK,
			responseBody:  `{"data": {"priceUsd": "50000"}}`,
			expectedError: "",
		},
		{
			name:  "Not Found",
			asset: models.Asset{Name: "xyu"},
			cfg: Config{
				Coincap: CoincapConfig{
					APIurl: "https://example.com",
					APIkey: "API_KEY",
				},
			},
			status:        http.StatusNotFound,
			responseBody:  `{"data": {"priceUsd": "50000"}}`,
			expectedError: "",
		},
		{
			name:  "Empty Asset Name",
			asset: models.Asset{Name: ""},
			cfg: Config{
				Coincap: CoincapConfig{
					APIurl: "https://example.com",
					APIkey: "API_KEY",
				},
			},
			status:        http.StatusNotFound,
			responseBody:  `{"data": {"priceUsd": "50000"}}`,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/bitcoin" {
					rw.WriteHeader(tt.status)
					rw.Write([]byte(tt.responseBody))
				}
			}))
			defer server.Close()

			client := NewCoincapClient()
			cfg := tt.cfg

			resp, err := GetAssetPrice(*client, tt.asset, cfg)
			if err != nil {
				if tt.expectedError == "" {
					t.Errorf("unexpected error: %v", err)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if tt.expectedError != "" {
					t.Errorf("expected error containing %q, got nil", tt.expectedError)
				}

				if string(resp) != tt.responseBody {
					t.Errorf("unexpected response body: got %q, want %q", string(resp), tt.responseBody)
				}
			}
		})
	}
}
