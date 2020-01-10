// Copyright 2020 Google LLC
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

package keystores

import (
	"net/url"
	"path"

	"github.com/srinandan/apigeecli/apiclient"
)

//Create
func Create(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	payload := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), "", "DELETE")
	return respBody, err
}

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}
