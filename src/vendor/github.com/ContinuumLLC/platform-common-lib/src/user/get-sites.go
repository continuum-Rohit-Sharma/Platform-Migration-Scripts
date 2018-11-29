package user

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (s service) getSiteIDs(httpClient *http.Client, partnerID, token string) (siteIDs []int64, err error) {
	var (
		sitesData struct {
			SiteList []struct {
				ID int64 `json:"siteId"`
			} `json:"siteDetailList"`
		}
		request  *http.Request
		response *http.Response
		url      = fmt.Sprintf("%s/partner/%s/sites", s.sitesELB, partnerID)
	)
	siteIDs = make([]int64, 0)

	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	response, err = httpClient.Do(request)
	if err != nil {
		return
	}

	defer func() {
		err = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		err = errors.Errorf("got wrong http status [%d]; expected status [%d]", response.StatusCode, http.StatusOK)
		return
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		err = readErr
		return
	}

	err = json.Unmarshal(body, &sitesData)
	if err != nil {
		return
	}

	for _, site := range sitesData.SiteList {
		siteIDs = append(siteIDs, site.ID)
	}
	return
}
