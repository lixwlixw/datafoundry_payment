package openshift

import (
	"fmt"
	"net/http"

	projectapi "github.com/openshift/origin/pkg/project/api/v1"
	rolebindingapi "github.com/openshift/origin/pkg/rolebinding/api/v1"
	"github.com/zonesan/clog"
)

type OClient struct {
	client *OpenshiftREST
	user   string
}

func NewOClient(host, token, username string) *OClient {

	clog.Debugf("%v:(%v)@%v", username, token, host)

	client := NewOpenshiftREST(NewOpenshiftTokenClient(host, token))
	return &OClient{client: client, user: username}
}

func NewAdminOClient(adminClient *OpenshiftClient) *OClient {
	adminRESTClient := NewOpenshiftREST(adminClient)
	client := &OClient{client: adminRESTClient}
	return client
}

func (oc *OClient) CreateProject(r *http.Request, name string) (*projectapi.ProjectRequest, error) {

	uri := "/projectrequests"

	proj := new(projectapi.ProjectRequest)
	{
		proj.DisplayName = name
		proj.Name = oc.user + "-org-" + genRandomName(8)
		proj.Annotations = make(map[string]string)
		proj.Annotations["datafoundry.io/requester"] = oc.user
	}

	oc.client.OPost(uri, proj, proj)
	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
		return nil, oc.client.Err
	}

	return proj, nil
}

func (oc *OClient) ListRoles(r *http.Request, project string) (*rolebindingapi.RoleBindingList, error) {
	uri := fmt.Sprintf("/namespaces/%v/rolebindings", project)

	roles := new(rolebindingapi.RoleBindingList)

	oc.client.OGet(uri, roles)
	//clog.Debug(roles)

	rolesResult := new(rolebindingapi.RoleBindingList)

	for _, role := range roles.Items {
		for _, subject := range role.Subjects {
			if subject.Kind == "User" {
				if role.RoleRef.Name == "view" || role.RoleRef.Name == "admin" ||
					role.RoleRef.Name == "edit" {
					//clog.Debugf("%#v", role)
					rolesResult.Items = append(rolesResult.Items, role)
					break
				}
			}
		}
	}
	return rolesResult, nil
}
