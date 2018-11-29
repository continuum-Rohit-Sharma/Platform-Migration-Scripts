package checksum

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

//sha1Impl is the struct for Sha1 validator
type sha1Impl struct{}

//Calculate is the method to get the Sha1 validator checksum
func (s sha1Impl) Calculate(reader io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//Validate is to verify and validate checksum
func (s sha1Impl) Validate(reader io.Reader, checksum string) (bool, error) {
	sum, err := s.Calculate(reader)
	if err != nil {
		return false, fmt.Errorf("Checksum cannot be Calculated from downloaded component, Err : %v", err)
	}
	log := logging.GetLoggerFactory().Get()
	if strings.ToUpper(sum) != strings.ToUpper(checksum) {
		log.Logf(logging.DEBUG, "Invalid CheckSum :- Calculated checksum : %s & downloaded file checksum : %s", sum, checksum)
		return false, errors.New(ErrChecksumInvalid)
	}
	return true, nil
}
