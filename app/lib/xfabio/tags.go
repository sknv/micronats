package xfabio

import (
	"fmt"
)

func Tags(serviceName string) []string {
	return []string{
		fmt.Sprintf("urlprefix-/%s strip=/%s", serviceName, serviceName),
	}
}
