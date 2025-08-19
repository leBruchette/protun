package client

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type ProtunClient struct {
	ProxyEndpoints []string
	HttpClient     *http.Client
}

type RequestConfig struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func NewProtunClient() *ProtunClient {
	proxyEndpointsEnv := os.Getenv("PROXY_ENDPOINTS")
	var proxyEndpoints []string

	if proxyEndpointsEnv != "" {
		proxyEndpoints = strings.Split(proxyEndpointsEnv, ",")
	} else {
		// Default endpoints for docker-compose services
		// FIXME
		// when running this outside of a container, we don't have access to http://<service-name>:<port> from docker's internal DNS
		// cheaply assigning to known constants from docker-compose.private.yml which is .gitignore'd
		proxyEndpoints = []string{
			"http://0.0.0.0:80",
			"http://0.0.0.0:81",
			"http://0.0.0.0:82",
		}
	}

	return &ProtunClient{
		ProxyEndpoints: proxyEndpoints,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (pc *ProtunClient) MakeRequest(config RequestConfig) (*http.Response, error) {
	proxyEndpoint := pc.getRandomProxy()
	// Construct proxy URL: http://proxy-service/proxy/target-url
	proxyURL := fmt.Sprintf("%s/proxy/%s", proxyEndpoint, strings.TrimPrefix(config.URL, "https://"))

	log.Printf("Using proxy: %s for request to: %s", proxyEndpoint, config.URL)

	var body io.Reader
	if config.Body != "" {
		body = strings.NewReader(config.Body)
	}

	req, err := http.NewRequest(config.Method, proxyURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	return pc.HttpClient.Do(req)
}

func (pc *ProtunClient) getRandomProxy() string {
	if len(pc.ProxyEndpoints) == 0 {
		log.Fatal("No proxy endpoints configured")
	}
	return pc.ProxyEndpoints[rand.Intn(len(pc.ProxyEndpoints))]
}
