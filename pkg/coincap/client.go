package coincap

import (
	"fmt"
	"io"
	"log"
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

func InitConfig() *Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfgPath := filepath.Join(dir, "config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(cfgFile, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func setHeaders(apiKey string) http.Header {
	headers := http.Header{
		"User-Agent":    {"Mozilla/5.0 (iPhone; CPU iPhone OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Mobile/15E148 Safari/604.1"},
		"Accept":        {"*/*"},
		"Connection":    {"Keep-alive"},
		"Authorization": {fmt.Sprintf("BEARER %s", apiKey)},
	}

	return headers
}

func GetAssetPrice(client CoincapClient, asset models.Asset, cfg Config) ([]byte, error) {
	url := cfg.APIurl + "/" + asset.Name

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
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
