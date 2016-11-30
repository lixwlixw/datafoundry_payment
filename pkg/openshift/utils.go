package openshift

import (
	"fmt"
	"strings"
)

func httpsAddr(addr string) string {
	if strings.HasSuffix(addr, "/") {
		addr = strings.TrimRight(addr, "/")
	}

	if !strings.HasPrefix(addr, "https://") {
		return fmt.Sprintf("https://%s", addr)
	}

	return addr
}
