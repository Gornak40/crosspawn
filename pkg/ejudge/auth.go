package ejudge

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AuthHeader struct {
	ContestID uint
	Login     string
	Password  string
}

func (ej *EjClient) verifyAuth(ctx context.Context, method string, auth AuthHeader) error {
	link, err := url.JoinPath(ej.cfg.URL, method)
	if err != nil {
		return err
	}

	params := url.Values{
		"contest_id": {strconv.Itoa(int(auth.ContestID))},
		"login":      {auth.Login},
		"password":   {auth.Password},
		"action":     {"2"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, link, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := ej.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}

	if resp.Request.URL.Query().Get("SID") == "" {
		return ErrInvalidCredentials
	}

	return nil
}

func (ej *EjClient) AuthUser(auth AuthHeader) error {
	return ej.verifyAuth(context.TODO(), "new-client", auth)
}
