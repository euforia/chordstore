package chordstore

import (
	"fmt"
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
