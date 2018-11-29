package user

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	realmHeader             = `realm`
	userNameHeader          = `username`
	nocRealm                = `/activedirectory`
	iPlanetDirectoryPro     = `iPlanetDirectoryPro`
	iPlanetDirectoryProPlus = `iPlanetDirectoryProPlus`
	uidHeader               = `uid`
)

type Service interface {
	GetUser(r *http.Request, httpClient *http.Client, needSites bool) (User, error)
}

type service struct {
	sitesELB string
}

func NewService(sitesELB string) Service {
	return service{
		sitesELB: sitesELB,
	}
}

type User interface {
	PartnerID() string
	Name() string
	UID() string
	Token() string
	HasNOCAccess() bool
	SiteIDs() []int64
}

type user struct {
	name         string
	partnerID    string
	uid          string
	token        string
	hasNOCAccess bool
	siteIDs      []int64
}

func (s service) GetUser(r *http.Request, httpClient *http.Client, needSites bool) (User, error) {
	var (
		partnerID    = mux.Vars(r)["partnerID"]
		name         = r.Header.Get(userNameHeader)
		uid          = r.Header.Get(uidHeader)
		hasNOCAccess = false
		err          error
		siteIDs      []int64
	)

	token := r.Header.Get(iPlanetDirectoryProPlus)
	if len(token) == 0 {
		token = r.Header.Get(iPlanetDirectoryPro)
	}

	if realm := r.Header.Get(realmHeader); realm == nocRealm {
		hasNOCAccess = true
	}

	if needSites {
		if len(partnerID) == 0 || len(token) == 0 {
			return nil, errors.New("empty headers")
		}

		siteIDs, err = s.getSiteIDs(httpClient, partnerID, token)
		if err != nil {
			return nil, err
		}
	}

	return user{
		name:         name,
		partnerID:    partnerID,
		uid:          uid,
		token:        token,
		hasNOCAccess: hasNOCAccess,
		siteIDs:      siteIDs,
	}, nil
}

func (u user) PartnerID() string {
	return u.partnerID
}

func (u user) Name() string {
	return u.name
}

func (u user) UID() string {
	return u.uid
}

func (u user) Token() string {
	return u.token
}

func (u user) HasNOCAccess() bool {
	return u.hasNOCAccess
}

func (u user) SiteIDs() []int64 {
	return u.siteIDs
}
