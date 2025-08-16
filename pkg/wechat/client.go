package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client 微信API客户端
type Client struct {
	httpClient *http.Client
	baseURL    string
}

var (
	globalClient *Client
	once         sync.Once
)

// NewClient 创建微信API客户端
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// InitClient 初始化全局微信客户端
func InitClient(baseURL string) {
	once.Do(func() {
		globalClient = NewClient(baseURL)
	})
}

// GetClient 获取全局微信客户端
func GetClient() *Client {
	if globalClient == nil {
		InitClient("https://api.weixin.qq.com")
	}
	return globalClient
}

// AccessTokenRequest 获取access_token请求参数
type AccessTokenRequest struct {
	AppID     string `json:"appid"`
	Secret    string `json:"secret"`
	GrantType string `json:"grant_type"`
}

// AccessTokenResponse 获取access_token响应
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
}

// GetAccessToken 获取access_token
func (c *Client) GetAccessToken(ctx context.Context, appID, secret string) (*AccessTokenResponse, error) {
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("secret", secret)
	params.Set("grant_type", "client_credential")

	url := fmt.Sprintf("%s/cgi-bin/token?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result AccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}
