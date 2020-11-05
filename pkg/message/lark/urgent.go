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
	urgentTypeApp = "app"
)

const (
	urgentPath = "/open-apis/message/v4/urgent/"
)

type urgentResponse struct {
	Code           int      `json:"code"`
	InvalidOpenIds []string `json:"invalid_open_ids"`
}

type urgentRequest struct {
	MessageID  string   `json:"message_id"`
	UrgentType string   `json:"urgent_type"`
	OpenIds    []string `json:"open_ids"`
}

func sendUrgent(ctx context.Context, appID string, request *urgentRequest) (*urgentResponse, error) {
	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(urgentPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("urgent reponse status code: %d, not 200", statusCode)
	}

	resp := &urgentResponse{}
	if err := json.Unmarshal(rspBytes, resp); err != nil {
		return nil, fmt.Errorf("unmarshal urgent response error: %s", err.Error())
	}

	return resp, nil
}
