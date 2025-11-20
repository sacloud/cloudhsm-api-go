// Copyright 2025- The sacloud/cloudhsm-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudhsm

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	client "github.com/sacloud/api-client-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	saht "github.com/sacloud/go-http"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1b/api/cloud/1.1/"
)

var (
	// UserAgent APIリクエスト時のユーザーエージェント
	UserAgent = fmt.Sprintf(
		"cloudhsm-api-go/%s (%s/%s; +https://github.com/sacloud/cloudhsm-api-go) %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		client.DefaultUserAgent,
	)

	RequestCustomizers = []saht.RequestCustomizer{
		func(req *http.Request) error {
			req.Header.Set("X-Sakura-Bigint-As-Int", "1")
			return nil
		},
	}
)

type EmptySecuritySource struct{}

func (this EmptySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{}, nil
}

func NewClient(params ...client.ClientParam) (*v1.Client, error) {
	return NewClientWithApiUrl(DefaultAPIRootURL, params...)
}

func NewClientWithApiUrl(apiUrl string, params ...client.ClientParam) (*v1.Client, error) {
	return NewClientWithApiUrlAndClient(apiUrl, nil, params...)
}

func NewClientWithApiUrlAndClient(apiUrl string, apiClient *http.Client, params ...client.ClientParam) (*v1.Client, error) {
	var cli, opts client.ClientParam
	if apiClient == nil {
		cli = func(i *client.ClientParams) {}
	} else {
		cli = client.WithHTTPClient(apiClient)
	}
	ua := client.WithUserAgent(UserAgent)
	opts = func(p *client.ClientParams) {
		if p.Options == nil {
			p.Options = &client.Options{}
		}
		if p.Options.RequestCustomizers == nil {
			p.Options.RequestCustomizers = []saht.RequestCustomizer{}
		}
		p.Options.RequestCustomizers = append(p.Options.RequestCustomizers, RequestCustomizers...)
	}
	c, err := client.NewClient(apiUrl, append(params, ua, cli, opts)...)
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	d, err := v1.NewClient(c.ServerURL(), EmptySecuritySource{}, v1.WithClient(c.NewHttpRequestDoer()))
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	return d, nil
}

func WithZone(z string) client.ClientParam {
	// 現在対応している既知のゾーン
	switch z {
	case "is1b":
		return func(p *client.ClientParams) {
			p.APIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1b/api/cloud/1.1/"
		}
	case "tk1a":
		return func(p *client.ClientParams) {
			p.APIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1/"
		}
	default:
		// 未知(あるいは未サポート)のゾーン
		// エラーを返す方法がないのでpanic
		// funcの中でpanicすることもできるが、エラー検知は早い方がいいだろう
		panic("unsupported zone: " + z)
	}
}
