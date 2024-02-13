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
	"github.com/smokinjoints/crypto-price-bot/pkg/models"
)

type CoincapClient struct {
	httpClient http.Client
}

type CoincapConfig struct {
	APIurl    string
	BotAPIkey string
	APIkey    string
}

type Config struct {
	Coincap CoincapConfig
}

var once sync.Once

func NewCoincapClient() *CoincapClient {
	var client *CoincapClient

	once.Do(func() {
		client = &CoincapClient{
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
	})

	return client
}

func ReadConfig() *Config {
	var config *Config

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
	ReadConfig()
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
