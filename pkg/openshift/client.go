package openshift

import (
	"fmt"
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/pkg"
	projectapi "github.com/openshift/origin/pkg/project/api/v1"
	rolebindingapi "github.com/openshift/origin/pkg/rolebinding/api/v1"
	"github.com/zonesan/clog"
	kapi "k8s.io/kubernetes/pkg/api/v1"
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

func (oc *OClient) CreateProject(r *http.Request, name string) (*projectapi.Project, error) {

	uri := "/projectrequests"

	projRequest := new(projectapi.ProjectRequest)
	{
		projRequest.DisplayName = name
		projRequest.Name = oc.user + "-org-" + genRandomName(8)
		// projRequest.Annotations = make(map[string]string)
		// projRequest.Annotations["datafoundry.io/requester"] = oc.user
	}

	proj := new(projectapi.Project)

	oc.client.OPost(uri, projRequest, proj)
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

	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
		return nil, oc.client.Err
	}
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

func (oc *OClient) RoleAdd(r *http.Request, project, name string, admin bool) (*rolebindingapi.RoleBinding, error) {

	if name == "" {
		return nil, pkg.ErrorNew(pkg.ErrCodeInvalidParam)
	}

	roles, err := oc.ListRoles(r, project)
	if err != nil {
		clog.Error(err)
		return nil, err
	}
	roleRef := "edit"
	if admin {
		roleRef = "admin"
	}
	role := findRole(roles, roleRef)
	create := false

	if role == nil { //post else put
		create = true
		role = new(rolebindingapi.RoleBinding)
		role.Name = roleRef
	}

	if exist := findUserInRoles(roles, name); exist {
		return nil, pkg.ErrorNew(pkg.ErrCodeConflict)
	}

	subject := kapi.ObjectReference{Kind: "User", Name: name}
	role.Subjects = append(role.Subjects, subject)
	role.UserNames = append(role.UserNames, name)

	uri := fmt.Sprintf("/namespaces/%v/rolebindings", project)

	if create {
		oc.client.OPost(uri, role, role)
	} else {
		uri += "/" + roleRef
		oc.client.OPut(uri, role, role)
	}

	return role, oc.client.Err
}

func (oc *OClient) RoleRemove(r *http.Request, project, name string) error {
	return nil
}

func findRole(roles *rolebindingapi.RoleBindingList, roleRef string) *rolebindingapi.RoleBinding {
	for _, role := range roles.Items {
		if role.RoleRef.Name == roleRef {
			return &role
		}
	}
	return nil
}

func findUserInRole(users []string, user string) bool {
	for _, v := range users {
		if user == v {
			return true
		}
	}
	return false
}

func findUserInRoles(roles *rolebindingapi.RoleBindingList, username string) bool {
	for _, role := range roles.Items {
		if exist := findUserInRole(role.UserNames, username); exist {
			clog.Warnf("duplicate user: %v, role: %v", username, role.RoleRef.Name)
			return exist
		}
	}
	return false
}
