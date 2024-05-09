package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"project_go/webbook/internal/domain"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WeChatDomain, error)
}

type WeChatAuthService struct {
	appId    string
	secretId string
	client   *http.Client
}

func (w WeChatAuthService) VerifyCode(ctx context.Context, code string) (domain.WeChatDomain, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		w.appId, w.secretId, code)
	res, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WeChatDomain{}, err
	}
	httpResp, err := w.client.Do(res)
	if err != nil {
		return domain.WeChatDomain{}, err
	}
	var resBody Result
	err = json.NewDecoder(httpResp.Body).Decode(&resBody)
	if err != nil {
		return domain.WeChatDomain{}, err
	}
	if resBody.ErrCode != 0 {
		return domain.WeChatDomain{}, fmt.Errorf("调用微信接口失败 errcode %d, errmsg %s", resBody.ErrCode, resBody.ErrMsg)
	}
	return domain.WeChatDomain{
		UnionId: resBody.UnionId,
		OpenId:  resBody.OpenId,
	}, nil

}

func NewService(appId string, secretId string) Service {
	return &WeChatAuthService{
		appId:    appId,
		secretId: secretId,
		client:   http.DefaultClient,
	}
}

const urlPattern = "open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

func (w WeChatAuthService) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(urlPattern, w.appId, redirectURL, state), nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
