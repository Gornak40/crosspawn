package ejudge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/Gornak40/crosspawn/config"
)

var (
	ErrBadStatusCode = errors.New("bad status code")
	ErrBadResult     = errors.New("bad result")
)

type EjClient struct {
	cfg    *config.EjConfig
	client *http.Client
}

func NewEjClient(cfg *config.EjConfig) *EjClient {
	jar, _ := cookiejar.New(nil)

	return &EjClient{
		cfg: cfg,
		client: &http.Client{
			Jar: jar,
		},
	}
}

type ejAnswer struct {
	OK    bool `json:"ok"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	Result json.RawMessage `json:"result"`
}

func (ej *EjClient) shootEjAPIRaw(ctx context.Context, httpMethod, method string, params url.Values) ([]byte, error) {
	link, err := url.JoinPath(ej.cfg.URL, method)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if httpMethod == http.MethodPost {
		body = strings.NewReader(params.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, httpMethod, link, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer AQAA"+ej.cfg.APIKey)

	switch httpMethod {
	case http.MethodGet:
		req.URL.RawQuery = params.Encode()
	case http.MethodPost:
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := ej.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (ej *EjClient) shootEjAPIGet(ctx context.Context, method string, params url.Values) (*ejAnswer, error) {
	data, err := ej.shootEjAPIRaw(ctx, http.MethodGet, method, params)
	if err != nil {
		return nil, err
	}

	return parseEjAPIAnswer(data)
}

func parseEjAPIAnswer(data []byte) (*ejAnswer, error) {
	var answer ejAnswer
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&answer); err != nil {
		return nil, err
	}
	if !answer.OK {
		return nil, fmt.Errorf("%w: %s", ErrBadResult, answer.Error.Message)
	}

	return &answer, nil
}
