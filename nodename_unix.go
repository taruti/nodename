package nodename

import (
	"errors"
	"net"
	"os"
	"strings"
)

// Get the name of the current machine as a short hostname, domain name, full name (host.domain) and an error.
func Get() (host string, domain string, full string, err error) {
	host, err = os.Hostname()
	if err != nil {
		return
	}
	host = removeTrailingDot(host)
	host, domain = split2(host, '.')
	if domain != "" {
		full = host + "." + domain
	} else {
		full, err = resolveNetFullname(host)
		full = removeTrailingDot(full)
		host, domain = split2(full, '.')
	}
	return
}

func removeTrailingDot(s string) string {
	if len(s) == 0 {
		return ""
	}
	if s[len(s)-1] == '.' {
		return s[:len(s)-1]
	}
	return s
}

func resolveNetFullname(host string) (string, error) {
	as, e := net.LookupHost(host)
	if e != nil {
		return "", e
	}
	hndot := host + "."
	for _, a := range as {
		hs, e := net.LookupAddr(a)
		if e != nil {
			continue
		}
		for _, h := range hs {
			if strings.HasPrefix(h, hndot) {
				return h, nil
			}
		}
	}
	return "", errors.New("network resolution of host domain failed")
}

func split2(s string, char byte) (string, string) {
	i := strings.IndexRune(s, rune(char))
	if i < 0 {
		return s, ""
	}
	return s[0:i], s[i+1:]
}
