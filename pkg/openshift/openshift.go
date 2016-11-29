package openshift

import (
	"fmt"

	project "github.com/openshift/origin/pkg/project/api/v1"
	rolebinding "github.com/openshift/origin/pkg/rolebinding/api/v1"
	user "github.com/openshift/origin/pkg/user/api/v1"
)

func Hello() {
	proj := project.Project{}
	usr := user.User{}
	role := rolebinding.RoleBinding{}

	fmt.Printf("%#v\n%#v\n%#v", proj, usr, role)

}
