package user

import (
	"net/http"
)

type mock struct {
	name, partnerID, uid, token string
	hasNOCAccess                bool
	siteIDs                     []int64
	err                         error
}

func NewMock(
	name,
	partnerID,
	uid,
	token string,
	hasNOCAccess bool,
	siteIDs []int64,
	err error,
) Service {
	return mock{
		name:         name,
		partnerID:    partnerID,
		uid:          uid,
		token:        token,
		hasNOCAccess: hasNOCAccess,
		siteIDs:      siteIDs,
		err:          err,
	}
}

func (s mock) GetUser(r *http.Request, httpClient *http.Client, needSites bool) (User, error) {
	if s.err != nil {
		return nil, s.err
	}

	return user{
		name:         s.name,
		partnerID:    s.partnerID,
		uid:          s.uid,
		token:        s.token,
		hasNOCAccess: s.hasNOCAccess,
		siteIDs:      s.siteIDs,
	}, nil
}
