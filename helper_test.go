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

package cloudhsm_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	client "github.com/sacloud/api-client-go"
	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Message string `json:"error_msg"`
	IsOk    bool   `json:"is_ok"`
}

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		IsOk:    false,
	}
}

func newTestClient(v any, s ...int) *v1.Client {
	s = append(s, http.StatusOK)
	j, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st := s[0]

		w.WriteHeader(st)
		if st == http.StatusNoContent {
			return
		}
		if _, e = w.Write(j); e != nil {
			panic(e)
		}
	})
	sv := httptest.NewServer(h)
	c, e := NewClientWithApiUrlAndClient(sv.URL, sv.Client())
	if e != nil {
		panic(e)
	}
	return c
}

func newIntegratedClient(t *testing.T, params ...client.ClientParam) *v1.Client {
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
	)(t)

	apiUrl := DefaultAPIRootURL
	if root, ok := os.LookupEnv("SAKURACLOUD_LOCAL_ENDPOINT_CLOUDHSM"); ok {
		apiUrl = root
	}
	ret, err := NewClientWithApiUrl(apiUrl, append(params, client.WithApiKeys(
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
	))...)

	require.NoError(t, err)
	return ret
}

func ref[T any](v T) *T {
	return &v
}

var TemplateDateTime = func() v1.DateTime {
	var ret v1.DateTime
	ret.SetFake()
	return ret
}

var TemplateTags = []string{"tag1", "tag2"}

var TemplateLicense = func() v1.CloudHSMSoftwareLicense {
	var ret v1.CloudHSMSoftwareLicense
	ret.SetFake()
	ret.SetTags(TemplateTags)

	return ret
}()

var TemplateCreateLicense = func() v1.CreateCloudHSMSoftwareLicense {
	var ret v1.CreateCloudHSMSoftwareLicense
	ret.SetFake()
	ret.SetTags(TemplateTags)

	return ret
}()

var TemplateWrappedCreateLicense = func() v1.WrappedCreateCloudHSMSoftwareLicense {
	var ret v1.WrappedCreateCloudHSMSoftwareLicense
	ret.SetCloudHSM(TemplateCreateLicense)

	return ret
}()

var TemplateWrappedLicense = func() v1.WrappedCloudHSMSoftwareLicense {
	var ret v1.WrappedCloudHSMSoftwareLicense
	ret.SetCloudHSM(TemplateLicense)

	return ret
}()
