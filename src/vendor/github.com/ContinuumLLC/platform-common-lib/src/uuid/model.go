package uuid

import (
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/google/uuid"
)

const (
	//ErrCantParseUUIDString : Error if parse fails
	ErrCantParseUUIDString = "ErrCantParseUUIDString"
)

//NewRandomUUID : Generates the new random UUID
func NewRandomUUID() (newuuid uuid.UUID, err error) {
	return uuid.New(), nil
}

//ParseUUID : Parses the given string UUID
func ParseUUID(xxxx string) (parseuuid uuid.UUID, err error) {
	parseuuid, err = uuid.Parse(xxxx)
	if err != nil {
		err = exception.New(ErrCantParseUUIDString, nil)
	}
	return
}
