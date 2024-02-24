package coincap

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tibeahx/tg-bot/models"
)

type CoincapClient struct {
	httpClient http.Client
	config     Config
}

type CoincapConfig struct {
	APIurl    string
	BotAPIkey string
	APIkey    string
}

type Config struct {
	Coincap CoincapConfig
}

func NewCoincapClient(cfg Config) *CoincapClient {
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

func ReadConfig() *Config {
	var config *Config

	var once sync.Once
	once.Do(func() {
		config = &Config{
			Coincap: CoincapConfig{
				APIurl:    os.Getenv("API_URL"),
				APIkey:    os.Getenv("API_KEY"),
				BotAPIkey: os.Getenv("BOT_API_KEY"),
			},
		}
	})

	return config
}

func setHeaders(apiKey string) http.Header {
	headers := http.Header{
		"Accept":        {"*/*"},
		"Connection":    {"Keep-alive"},
		"Authorization": {fmt.Sprintf("Bearer %s", apiKey)},
	}

	return headers
}

func GetAssetPrice(client CoincapClient, asset models.Asset) ([]byte, error) {
	url := ReadConfig().Coincap.APIurl + "/" + asset.Name

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(ReadConfig().Coincap.APIkey)

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
