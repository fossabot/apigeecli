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

package apiclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/clilog"
)

//CrmURL is the endpoint for cloud resource manager
const CrmURL = "https://cloudresourcemanager.googleapis.com/v1/projects/"

//binding for IAM Roles
type roleBinding struct {
	Role      string     `json:"role,omitempty"`
	Members   []string   `json:"members,omitempty"`
	Condition *condition `json:"condition,omitempty"`
}

//IamPolicy holds the response
type iamPolicy struct {
	Version  int           `json:"version,omitempty"`
	Etag     string        `json:"etag,omitempty"`
	Bindings []roleBinding `json:"bindings,omitempty"`
}

//SetIamPolicy holds the request to set IAM
type setIamPolicy struct {
	Policy iamPolicy `json:"policy,omitempty"`
}

//condition for Bindings
type condition struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

//CreateIAMServiceAccount create a new IAM SA with the necessary roles for Apigee
func CreateIAMServiceAccount(name string, iamRole string) (err error) {
	type KeyResponse struct {
		Name            string `json:"name,omitempty"`
		PrivateKeyType  string `json:"privateKeyType,omitempty"`
		PrivateKeyData  string `json:"privateKeyData,omitempty"`
		ValidBeforeTime string `json:"validBeforeTime,omitempty"`
		ValidAfterTime  string `json:"validAfterTime,omitempty"`
		KeyAlgorithm    string `json:"keyAlgorithm,omitempty"`
	}

	const iamURL = "https://iam.googleapis.com/v1/projects/"
	const crmBetaURL = "https://cloudresourcemanager.googleapis.com/v1beta1/projects/"
	var role string

	serviceAccountName := name + "@" + GetProjectID() + ".iam.gserviceaccount.com"

	switch iamRole {
	case "sync":
		role = "roles/apigee.synchronizerManager"
	case "analytics":
		role = "roles/apigee.analyticsAgent"
	case "metric":
		role = "roles/monitoring.metricWriter"
	case "logger":
		role = "roles/logging.logWriter"
	case "mart":
		role = ""
	case "cassandra":
		role = "roles/storage.objectAdmin"
	case "connect":
		role = "roles/apigeeconnect.Agent"
	case "all":
		role = "not-necessary-to-add-this"
	default:
		return fmt.Errorf("invalid service account role")
	}

	//Step 1: create a new service account
	u, _ := url.Parse(iamURL)
	u.Path = path.Join(u.Path, GetProjectID(), "serviceAccounts")

	iamPayload := []string{}
	iamPayload = append(iamPayload, "\"accountId\":\""+name+"\"")
	iamPayload = append(iamPayload, "\"serviceAccount\": {\"displayName\": \""+name+"\"}")

	payload := "{" + strings.Join(iamPayload, ",") + "}"

	_, err = HttpClient(false, u.String(), payload)

	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	//Step 2: create a new service account key
	u, _ = url.Parse(iamURL)
	u.Path = path.Join(u.Path, GetProjectID(), "serviceAccounts",
		serviceAccountName, "keys")

	respKeyBody, err := HttpClient(false, u.String(), "")

	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	//Step 3: read the response
	keyResponse := KeyResponse{}
	err = json.Unmarshal(respKeyBody, &keyResponse)
	if err != nil {
		return err
	}

	//Step 4: base64 decode the response to get the private key.json
	privateKey, err := base64.StdEncoding.DecodeString(keyResponse.PrivateKeyData)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	//Step 5: Write the data to a file
	file, err := os.Create(GetProjectID() + "-" + name + ".json")
	if err != nil {
		clilog.Error.Println("cannot open private key file: ", err)
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(privateKey))
	if err != nil {
		clilog.Error.Println("error writing to file: ", err)
		return err
	}

	//mart doesn't need any roles, return here.
	if iamRole == "mart" {
		return err
	}

	//Step 6: get the current IAM policies for the project
	u, _ = url.Parse(CrmURL)
	u.Path = path.Join(u.Path, GetProjectID()+":getIamPolicy")
	respBody, err := HttpClient(false, u.String(), "")

	iamPolicy := iamPolicy{}

	err = json.Unmarshal(respBody, &iamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	//Step 7: create a new policy binding for apigee
	if iamRole == "all" {
		bindings := createAllRoleBindings(serviceAccountName)
		iamPolicy.Bindings = append(iamPolicy.Bindings, bindings...)
	} else {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, "serviceAccount:"+serviceAccountName)

		iamPolicy.Bindings = append(iamPolicy.Bindings, binding)
	}

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = iamPolicy
	setIamPolicyBody, err := json.Marshal(setIamPolicy)

	//Step 8: set the iam policy
	u, _ = url.Parse(crmBetaURL)
	u.Path = path.Join(u.Path, GetProjectID()+":setIamPolicy")

	_, err = HttpClient(false, u.String(), string(setIamPolicyBody))

	return err
}

func createAllRoleBindings(name string) []roleBinding {
	var roles = [...]string{"roles/apigee.synchronizerManager", "roles/apigee.analyticsAgent",
		"roles/monitoring.metricWriter", "roles/logging.logWriter", "roles/storage.objectAdmin",
		"roles/apigeeconnect.Agent"}

	bindings := []roleBinding{}

	for _, role := range roles {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, "serviceAccount:"+name)
		bindings = append(bindings, binding)
	}

	return bindings
}

//SetIAMServiceAccount create a new IAM SA with the necessary roles for an Apigee Env
func SetIAMServiceAccount(serviceAccountName string, iamRole string) (err error) {
	var role string

	switch iamRole {
	case "sync":
		role = "roles/apigee.synchronizerManager"
	case "analytics":
		role = "roles/apigee.analyticsAgent"
	case "deploy":
		role = "roles/apigee.deployer"
	default:
		return fmt.Errorf("invalid service account role")
	}

	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":getIamPolicy")
	getIamPolicyBody, err := HttpClient(false, u.String())
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	getIamPolicy := iamPolicy{}

	err = json.Unmarshal(getIamPolicyBody, &getIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	foundRole := false
	for i, binding := range getIamPolicy.Bindings {
		if binding.Role == role {
			//found members with the role already, add the new SA to the role
			getIamPolicy.Bindings[i].Members = append(binding.Members, "serviceAccount:"+serviceAccountName)
			foundRole = true
		}
	}

	//no members with the role, add a new one
	if !foundRole {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, "serviceAccount:"+serviceAccountName)
		getIamPolicy.Bindings = append(getIamPolicy.Bindings, binding)
	}

	u, _ = url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":setIamPolicy")

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = getIamPolicy

	setIamPolicyBody, err := json.Marshal(setIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	_, err = HttpClient(false, u.String(), string(setIamPolicyBody))

	return err
}