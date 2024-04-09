package tests

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestGetUserBanner(t *testing.T) {

	var (
		tagId           = "7"
		featureId       = "3"
		useLastRevision = "false"
	)

	uri := "/user_banner?tag_id=" + tagId + "&feature_id=" + featureId + "&use_last_revision=" + useLastRevision
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, nil)
	if err != nil {
		t.Errorf("send request error: %v", err)
	}

	fmt.Println(resp)
}
