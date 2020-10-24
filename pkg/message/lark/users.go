package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
)

const (
	batchGetIDPath = "/open-apis/user/v1/batch_get_id"
)

type getOpenIDRequest struct {
	Emails []string `json:"emails"`
}

type getOpenIDResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *getOpenIDResponseData `json:"data"`
}

type getOpenIDResponseData struct {
	EmailUsers map[string][]*userInfo `json:"email_users"`
}

type userInfo struct {
	OpenID string `json:"open_id"`
	UserID string `json:"user_id"`
}

func getOpenID(ctx context.Context, appID string, email string) (string, error) {
	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return "", err
	}

	request := getOpenIDRequest{Emails: []string{email}}

	rspBytes, statusCode, err := common.DoHttpPostOApi(batchGetIDPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return "", common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("get open id reponse status code: %d, not 200", statusCode)
	}

	resp := &getOpenIDResponse{}
	if err := json.Unmarshal(rspBytes, resp); err != nil {
		return "", fmt.Errorf("unmarshal get open id response error: %s", err.Error())
	}

	if resp.Code != 0 {
		return "", fmt.Errorf("get open id response error: code=%d, msg=%s", resp.Code, resp.Message)
	}

	return resp.Data.EmailUsers[email][0].OpenID, nil
}
