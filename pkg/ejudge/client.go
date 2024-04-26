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

	"github.com/Gornak40/crosspawn/config"
)

var (
	ErrBadStatusCode      = errors.New("bad status code")
	ErrBadResult          = errors.New("bad result")
	ErrInvalidCredentials = errors.New("invalid credentials")
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

func (ej *EjClient) shootRaw(ctx context.Context, method string, params url.Values) ([]byte, error) {
	link, err := url.JoinPath(ej.cfg.URL, method)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer AQAA"+ej.cfg.APIKey)
	req.URL.RawQuery = params.Encode()

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

func (ej *EjClient) shootAPI(ctx context.Context, method string, params url.Values) (*ejAnswer, error) {
	data, err := ej.shootRaw(ctx, method, params)
	if err != nil {
		return nil, err
	}

	var answer ejAnswer
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&answer); err != nil {
		return nil, err
	}
	if !answer.OK {
		return nil, fmt.Errorf("%w: %s", ErrBadResult, answer.Error.Message)
	}

	return &answer, nil
}
