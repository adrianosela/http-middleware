package middleman

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Firewall is a software defined, endpoint-selective firewall for HTTP servers
type Firewall struct {
	rules map[string][]net.IPNet
	log   bool
}

var (
	// ErrPathHasRule will be returned when the developer attempts to re-assign a rule to a path
	ErrPathHasRule = errors.New("path already has an associated list of trusted netblocks")
	// ErrCouldNotParseCIDR will be returned when the developer attempts to use an invalid CIDR for a rule
	ErrCouldNotParseCIDR = fmt.Errorf("could not parse CIDR")
)

const (
	firewallHeader = "[firewall]"
)

// NewFirewall is the constructor for a Firewall Middleware
func NewFirewall(rules map[string][]net.IPNet, log bool) *Firewall {
	return &Firewall{
		rules: rules,
		log:   log,
	}
}

// AddPathRule maps a list of trusted netblocks to a given path
func (fw *Firewall) AddPathRule(path string, networks []string) error {
	if _, exists := fw.rules[path]; exists {
		return ErrPathHasRule
	}
	// parse network CIDRs
	var trusted []net.IPNet
	for _, network := range networks {
		_, trustedNetblock, err := net.ParseCIDR(network)
		if err != nil {
			return ErrCouldNotParseCIDR
		}
		trusted = append(trusted, *trustedNetblock)
	}
	// add trusted netblocks to path
	fw.rules[path] = trusted
	return nil
}

// Wrap the firewall around an HTTP handler function
func (fw *Firewall) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srcIP := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
		rule, hasRule := fw.rules[r.URL.Path]
		if authorized := (hasRule && ipIsTrusted(rule, srcIP)); !authorized {
			if fw.log {
				fmt.Println(fmt.Sprintf("%s blocked request from %s for %s", firewallHeader, srcIP.String(), r.URL.Path))
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// ipIsTrusted checks whether an IP address is part of a list of trusted netblocks
func ipIsTrusted(trusted []net.IPNet, src net.IP) bool {
	if src == nil {
		return false
	}
	for _, netblock := range trusted {
		if netblock.Contains(src) {
			return true
		}
	}
	return false
}
