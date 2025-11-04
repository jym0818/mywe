package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jym0818/mywe/internal/domain"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type WechatService interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.Wechat, error)
}
type wechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewwechatService(appId string, appSecret string) WechatService {
	return &wechatService{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (svc *wechatService) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, svc.appId, redirectURI, state), nil

}
func (svc *wechatService) VerifyCode(ctx context.Context, code string, state string) (domain.Wechat, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, svc.appId, redirectURI, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.Wechat{}, err
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return domain.Wechat{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.Wechat{}, err
	}
	if res.ErrCode != 0 {
		return domain.Wechat{}, fmt.Errorf("微信返回错误信息:%s", res.ErrMsg)
	}
	return domain.Wechat{
		OpenID:  res.OpenId,
		UnionID: res.UnionId,
	}, nil
}

// 根据腾讯文档的返回数据定义的结构体
type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	UnionId string `json:"unionid"`
	OpenId  string `json:"openid"`
}
