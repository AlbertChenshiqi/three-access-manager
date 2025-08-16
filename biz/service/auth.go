package service

import (
	"context"
	"fmt"
	"third-login/config"
	"time"

	"third-login/pkg/redis"
	"third-login/pkg/wechat"
)

// AuthService 认证服务
type AuthService struct {
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// TokenResult Token获取结果
type TokenResult struct {
	AccessToken string `json:"access_token"`
}

// GetToken 统一Token获取接口
func (s *AuthService) GetToken(ctx context.Context, platform, appID string) (string, error) {
	// 获取应用配置
	appConfig, err := config.GetAppConfig(platform, appID)
	if err != nil {
		return "", fmt.Errorf("failed to get app config: %w", err)
	}

	// 检查是否已有有效的access_token
	existingToken, err := redis.GetClient().GetAccessToken(ctx, platform, appID)
	if err == nil && existingToken != "" {
		// 重新获取
		return existingToken, nil
	}

	// 根据平台类型处理登录
	switch platform {
	case "wechat_miniprogram":
		return s.handleWechatMiniProgram(ctx, appConfig, platform, appID)
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}

// handleWechatMiniProgram 处理微信小程序登录
func (s *AuthService) handleWechatMiniProgram(ctx context.Context, appConfig *config.AppConfig, platform, appID string) (string, error) {
	// 获取微信客户端
	wechatClient := wechat.GetClient()

	// 获取微信access_token（不再处理用户登录逻辑）
	tokenResp, err := wechatClient.GetAccessToken(ctx, appID, appConfig.AppSecret)
	if err != nil {
		return "", fmt.Errorf("failed to get wechat access token: %w", err)
	}

	// 计算过期时间
	expiration := time.Duration(tokenResp.ExpiresIn) * time.Second

	// 存储access_token到Redis
	if err = redis.GetClient().SetAccessToken(ctx, platform, appID, tokenResp.AccessToken, expiration); err != nil {
		return "", fmt.Errorf("failed to store access token: %w", err)
	}

	// 更新统计数据
	go s.updateStats(context.Background(), "wechat_miniprogram", appID)

	return tokenResp.AccessToken, nil
}

// updateStats 更新统计数据
func (s *AuthService) updateStats(ctx context.Context, platform, appID string) {
	today := time.Now().Format("2006-01-02")

	// 增加每日登录统计
	redis.GetClient().IncrementDailyStats(ctx, platform, appID, today)

	// 增加总登录统计
	redis.GetClient().IncrementTotalStats(ctx, platform, appID)

	// 更新应用统计
	redis.GetClient().SetAppStats(ctx, appID, "last_login", time.Now().Unix())

	// 更新平台统计
	redis.GetClient().IncrementStats(ctx, fmt.Sprintf("stats:platform:%s:total_calls", platform))
}
