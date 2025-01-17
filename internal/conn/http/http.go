package http

import (
	"context"
	"io"
	"net/http"
)

type MetaHttpConn interface {
	GetDiagnose(ctx context.Context) (string, error)
}

type MetaHttpConnection struct {
	endpoint string
}

func NewMetaHttpConnection(endpoint string) MetaHttpConn {
	ep := endpoint
	if ep[len(ep)-1] == '/' {
		ep = ep[:len(ep)-1]
	}

	return &MetaHttpConnection{
		endpoint: ep,
	}
}

func (m *MetaHttpConnection) GetDiagnose(ctx context.Context) (string, error) {
	return get(ctx, m.endpoint+"/api/monitor/diagnose/")
}

func get(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
