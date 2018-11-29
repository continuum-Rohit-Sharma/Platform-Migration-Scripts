package checksum

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

type md5Impl struct{}

//Calculate is the method to get the MD5 validator checksum
func (c md5Impl) Calculate(reader io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//VerifyCheckSum is to verify and validate checksum
func (c md5Impl) Validate(reader io.Reader, checksum string) (bool, error) {
	log := logging.GetLoggerFactory().Get()
	calculatedVal, err := c.Calculate(reader)
	if err != nil {
		return false, fmt.Errorf("Checksum cannot be Calculated from downloaded component, Err : %v", err)
	}
	trimmedChecksum := strings.TrimSpace(checksum)
	trimmedCalculatedVal := strings.TrimSpace(calculatedVal)
	if strings.ToUpper(trimmedCalculatedVal) != strings.ToUpper(trimmedChecksum) {
		log.Logf(logging.DEBUG, "Invalid CheckSum :- Calculated checksum : %s with length : %d & Reader checksum : %s with length : %d", trimmedCalculatedVal, trimmedChecksum, len(trimmedCalculatedVal), len(trimmedChecksum))
		return false, errors.New(ErrChecksumInvalid)
	}
	return true, nil
}
