package coincap

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
	"gopkg.in/yaml.v3"
)

type CoincapClient struct {
	httpClient http.Client
}

type Config struct {
	APIkey    string `yaml:"api_key"`
	APIurl    string `yaml:"api_url"`
	BotAPIkey string `yaml:"bot_api_key"`
}

func NewCoincapClient() *CoincapClient {
	return &CoincapClient{
		httpClient: http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
			},
			Timeout: 10 * time.Second,
		},
	}
}

func InitConfig() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfgPath := filepath.Join(dir, "config.yaml")

	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(cfgFile, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setHeaders(apiKey string) http.Header {
	headers := http.Header{
		"Accept":        {"*/*"},
		"Connection":    {"Keep-alive"},
		"Authorization": {fmt.Sprintf("Bearer %s", apiKey)},
	}

	return headers
}

func GetAssetPrice(client CoincapClient, asset models.Asset, cfg Config) ([]byte, error) {
	url := cfg.APIurl + "/" + asset.Name

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header = setHeaders(cfg.APIkey)

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error code %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return io.ReadAll(resp.Body)
}
