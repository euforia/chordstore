package chordstore

import "fmt"

/*type ChecksumWriter struct {
	hash   hash.Hash
	writer io.Writer
}

// NewChecksumWriter based on the provided hash function
func NewChecksumWriter(hash hash.Hash, w io.Writer) *ChecksumWriter {
	return &ChecksumWriter{hash: hash, writer: w}
}

// Write
func (w *ChecksumWriter) Write(p []byte) (n int, err error) {
	w.hash.Write(p)
	return w.writer.Write(p)
}

func (w *ChecksumWriter) Sum() []byte {
	return w.hash.Sum(nil)
}*/

func mergeErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return fmt.Errorf("%s\n%s", err1, err2)
	}
}
