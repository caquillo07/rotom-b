package dblgo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"time"
)

const (
	baseURL                  = "https://discordbots.org/api"
	errorWhileSendingRequest = "error while sending request:"
)

// DBLApi is base api client struct
type DBLApi struct {
	AccessToken    string        // DBL access token
	RequestTimeout time.Duration // Timeout for all requests
}

// NewDBLApi returns new DBLApi struct initialized with optimal values
func NewDBLApi(accessToken string) DBLApi {
	return DBLApi{
		AccessToken:    accessToken,
		RequestTimeout: time.Second * 10,
	}
}

// PostStatsSimple sends bot guild count to the website
func (d DBLApi) PostStatsSimple(guildCount int) error {
	url := d.getRequestURL("/bots/stats")
	params := map[string]string{"server_count": fmt.Sprintf("%d", guildCount)}

	result, err := d.getBaseRequest().SetBody(params).Post(url)
	if err != nil {
		return errors.Wrap(err, errorWhileSendingRequest)
	}

	statusCode := result.StatusCode()
	if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("server returned %d http code", statusCode))
	}

	return nil
}

func (d DBLApi) getBaseRequest() *resty.Request {
	ctx, _ := context.WithTimeout(context.Background(), d.RequestTimeout)

	return resty.R().SetHeader("Authorization", d.AccessToken).SetContext(ctx)
}

// getRequestURL appends endpoint to the baseURL, and return full request URL
func (d DBLApi) getRequestURL(endpoint string) string {
	return fmt.Sprintf("%s%s", baseURL, endpoint)
}
