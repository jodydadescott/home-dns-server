package util

import (
	"fmt"
	"strings"
)

func GetARPA(ip string) (string, error) {
	ipSplit := strings.Split(ip, ".")
	if len(ipSplit) != 4 {
		return "", fmt.Errorf("IP %s is invalid", ip)
	}
	return ipSplit[3] + "." + ipSplit[2] + "." + ipSplit[1] + "." + ipSplit[0] + ".in-addr.arpa.", nil
}
