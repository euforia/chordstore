package chordstore

import "fmt"

func mergeErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return fmt.Errorf("%s\n%s", err1, err2)
	}
}
