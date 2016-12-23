package chordstore

import (
	"fmt"
	"net"
	"strings"
)

func mergeErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return fmt.Errorf("%s\n%s", err1, err2)
	}
}

// ParsePeersList parses a comma separate list of peers
func ParsePeersList(peerList string) []string {
	out := []string{}
	for _, v := range strings.Split(peerList, ",") {
		if c := strings.TrimSpace(v); c != "" {
			out = append(out, c)
		}
	}
	return out
}

// AutoDetectIPAddress traverses interfaces eliminating, localhost, ifaces with
// no addresses and ipv6 addresses.  It returns a list by priority
func AutoDetectIPAddress() ([]string, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	out := []string{}

	for _, ifc := range ifaces {
		if strings.HasPrefix(ifc.Name, "lo") {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			// Add ipv4 addresses to list
			if ip.To4() != nil {
				out = append(out, ip.String())
			}
		}

	}
	if len(out) == 0 {
		return nil, fmt.Errorf("could not detect ip addresses")
	}

	return out, nil
}

// IsAdvertisableAddress checks if an address can be used as the advertise address.
func IsAdvertisableAddress(hp string) (bool, error) {
	pp := strings.Split(hp, ":")
	if len(pp) < 1 {
		return false, fmt.Errorf("could not parse: %s", hp)
	}

	if pp[0] == "" {
		return false, nil
	} else if pp[0] == "0.0.0.0" {
		return false, nil
	}

	if _, err := net.ResolveIPAddr("ip4", pp[0]); err != nil {
		return false, err
	}

	return true, nil
}

// getAdvertiseAddr returns the advertise address including the port based on the
// bind and advertise addresses passed in. If an adv. address is not provided, it
// tries use use the bind address if it is an actual ip or uses the first available
// and usable ip.
func getAdvertiseAddr(bindAddr, advAddr string) (string, error) {
	if advAddr == "" {
		if k, _ := IsAdvertisableAddress(bindAddr); k {
			return bindAddr, nil
		}

		addrs, err := AutoDetectIPAddress()
		if err != nil {
			return "", err
		}
		hp := strings.Split(bindAddr, ":")
		return fmt.Sprintf("%s:%s", addrs[0], hp[len(hp)-1]), nil
	}

	k, err := IsAdvertisableAddress(advAddr)
	if err == nil {
		if k {
			return advAddr, nil
		}
	}
	return "", err
}
