package http

import (
	"context"
	"io"
	"net/http"
)

type MetaHttpManagerInterface interface {
	GetDiagnose(ctx context.Context, endpoint string) (string, error)
}

type MetaHttpManager struct {
}

func NewMetaHttpManager() MetaHttpManagerInterface {
	return &MetaHttpManager{}
}

func (m *MetaHttpManager) GetDiagnose(ctx context.Context, endpoint string) (string, error) {
	return get(ctx, endpoint+"/api/monitor/diagnose/")
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
