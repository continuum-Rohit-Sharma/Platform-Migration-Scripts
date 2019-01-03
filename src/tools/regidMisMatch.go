package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"fmt"
	"os"

	"github.com/gocql/gocql"
)

type assetCollection struct {
	EndPointList []EndPointList `json:"endPointList"`
	TotalPoint   int            `json:"totalCount"`
}

//EndPointList is the structuire
type EndPointList struct {
	RegID        string `json:"regID,omitempty"`
	FriendlyName string `json:"friendlyName"`
	MachineName  string `json:"machineName"`
	SiteName     string `json:"siteName"`
	//OperatingSystem       string `json:"operatingSystem"`
	Availability int `json:"availability"`
	//IPAddress             string `json:"ipAddress"`
	RegType string `json:"regType"`
	//LmiStatus             int    `json:"lmiStatus"`
	//ResType               string `json:"resType"`
	SiteID int `json:"siteId"`
	//SmartDisk             int    `json:"smartDisk"`
	//Amt                   int    `json:"amt"`
	//MbSyncstatus          int    `json:"mbSyncstatus"`
	//EncryptedResourceName string `json:"encryptedResourceName"`
	//EncryptedSiteName     string `json:"encryptedSiteName"`
	// LmiHostID             string `json:"lmiHostId"`
	//RequestRegID string `json:"requestRegId"`
	EndpointID string `json:"endpointId"`
}
type Config struct {
	NewWebAPIURL            string
	OldWebAPIURL            string
	IsSite                  bool
	PartnerID               string
	SiteIDs                 string
	NumOfWorkers            int
	CassandraHostIPForAsset string
	CSVPathForPartners      string
	ResourceType            string
}

//Session session variable
var Session *gocql.Session
var clusterIP string
var legacyWebAPIURL, legacyWebAPIURLDesktop, newWebAPIURL string
var totalNoPartnersFromCSV, totalNoPartners int
var PartnerEndpoints map[string]map[string]fnEndpoint

//var SiteIDs string

const (
	queryToUpdateFriendlyName     = "UPDATE platform_asset_db.partner_asset_collection set friendly_name=? where partner_id=? AND endpoint_id=?"
	queryToGetEndpointBypartnerID = "SELECT endpoint_id, reg_id, friendly_name FROM  platform_asset_db.partner_asset_collection WHERE partner_id='%s'"
)

//LoadConfiguration is a method to load configuration File
func LoadConfiguration(filePath string) (Config, error) {
	mbMessage := &Config{}
	file, err := os.Open(filePath)
	if err != nil {
		return *mbMessage, err
	}
	defer file.Close()

	deser := json.NewDecoder(file)
	err = deser.Decode(mbMessage)
	if err != nil {
		return *mbMessage, err
	}
	return *mbMessage, nil
}

func main() {

	// commandArgs := os.Args[1:]
	// if len(commandArgs) != 6 {
	// 	fmt.Println("Usage: syncFriendlyName <csv File> <LegacyBaseURL> <RMM2.0BaseURL> <CassandraNodeIP> <# of workers> <resourceType>")
	// 	fmt.Println("Example: syncFriendlyName partnerIDs.csv https://rmmapi.dtitsupport247.net http://internal-continuum-asset-service-elb-int-1972580147.ap-south-1.elb.amazonaws.com 172.28.48.6 10 desktop/server")
	// 	return
	// }
	startTime := time.Now()

	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: syncfn config.json")
		return
	}

	fmt.Println("Tool Started... Time started : ", startTime)

	configObject, err := LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}

	fmt.Println("Migration Tool Started... Time started : ", startTime)
	excelPath := configObject.CSVPathForPartners
	excelPath = "partners.csv"
	legacyWebAPIURL = configObject.OldWebAPIURL
	legacyWebAPIURL = "https://rmmapi.itsupport247.net"
	// newWebAPIURL = configObject.NewWebAPIURL
	clusterIP = configObject.CassandraHostIPForAsset
	noOfWorkers := configObject.NumOfWorkers
	noOfWorkers = 40
	resourceType := configObject.ResourceType
	resourceType = "server"

	if legacyWebAPIURL != "" && resourceType != "" {
		url := legacyWebAPIURL + "/itswebapi/v1/partner/%s/endpoints?"
		legacyWebAPIURL = url + "resourceType=server"
		legacyWebAPIURLDesktop = url + "resourceType=desktop"
		fmt.Println("--------------------------------legacyWebAPIURL : ", legacyWebAPIURL)
		fmt.Println("--------------------------------legacyWebAPIURLDesktop : ", legacyWebAPIURLDesktop)

		//fmt.Printf("Legacy Service URL : %s\n", legacyWebAPIURL)
	}

	uniquePartnerIDs := make([]string, 0)
	uniquePartnerIDs, err = readAssetTables(excelPath)
	if err != nil {
		fmt.Println("Error Occured while readAssetTables, Error : ", err)
		return
	}
	totalNoPartnersUnique := len(uniquePartnerIDs)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan PartnerMetrics, totalNoPartnersUnique)

	for w := 1; w <= noOfWorkers; w++ {
		go processParters(w, jobs, results)
	}
	fmt.Println("--------------------------------Number of workers : ", noOfWorkers)

	for a := 0; a < len(uniquePartnerIDs); a++ {
		jobs <- uniquePartnerIDs[a]
	}
	fmt.Println("--------------------------Number of jobs ", totalNoPartnersUnique)
	close(jobs)

	fmt.Println("Waiting for results ")
	ptrRes := make([]PartnerMetrics, len(uniquePartnerIDs))
	var a int
	for a = 0; a < len(uniquePartnerIDs); a++ {
		ptrRes[a] = <-results
	}
	fmt.Printf("\n\n*** --------------------------Partner Level Metrics -------------------------------***\n\n")
	for a = 0; a < len(uniquePartnerIDs); a++ {
		//	fmt.Printf("Partner ID %s: RMM1Err: %s, RMM2Err: %s, NoEndpointsInRMM1: %v, NOE already in Sync: %d, NOE updated: %d, NOE with UpdateDBError: %d, NOE missing in RMM2: %d \n", ptrRes[a].ID, ptrRes[a].RMM1Err, ptrRes[a].RMM2Err, ptrRes[a].RMM1NoEndpointPartner, len(ptrRes[a].FNameInSyncEndpoints), len(ptrRes[a].UpdatedEndpoints), len(ptrRes[a].UpdateDBErrorEndpoints), len(ptrRes[a].RMM2NoEndpoints))
	}
	close(results)
	//Session.Close()

}

func getAllLegacyAssetsByPartnerID(partnerID string) (assetCollection, error) {
	var retErr error
	assetColl := assetCollection{}
	dekstopData, err := getLegacyPartnerData(partnerID, legacyWebAPIURLDesktop)
	if err != nil {
		retErr = err
		fmt.Println("------------dekstop errr in legacy for partner -------------------", partnerID)
	} else {
		assetColl.EndPointList = dekstopData.EndPointList
		assetColl.TotalPoint = dekstopData.TotalPoint
	}
	serverData, e := getLegacyPartnerData(partnerID, legacyWebAPIURL)
	if e != nil {
		retErr = e
		fmt.Println("------------serverData errr in legacy for partner -------------------", partnerID)
	} else {
		if assetColl.EndPointList == nil || len(assetColl.EndPointList) < 1 {
			assetColl.EndPointList = serverData.EndPointList
			assetColl.TotalPoint = serverData.TotalPoint
		} else {
			for _, val := range serverData.EndPointList {
				assetColl.EndPointList = append(assetColl.EndPointList, val)
				cnt := assetColl.TotalPoint
				cnt = cnt + 1
				assetColl.TotalPoint = cnt
			}
		}

	}
	return assetColl, retErr
}

func getLegacyPartnerData(partnerID string, url string) (assetCollection, error) {
	assetColl := assetCollection{}
	partnerID = strings.TrimSpace(partnerID)
	requestURL := fmt.Sprintf(url, partnerID)
	res, err := http.Get(requestURL)
	if err != nil {
		//fmt.Println("RMM 1.0 : Error 1.0 While getting the Response for Partner Id :", partnerID)
		return assetColl, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("RMM 1.0 : Error 2.0 while Parsing the Response body for Partner Id :", partnerID)
		return assetColl, err
	}
	bodyString := string(bodyBytes)

	err = json.Unmarshal([]byte(bodyString), &assetColl)
	if err != nil {
		//fmt.Println("RMM 1.0 : Error 3.0 Occured while unmarshalling the Resp for Partner Id :", partnerID)
		return assetColl, err
	}
	return assetColl, nil
}

type fnEndpoint struct {
	EID   string
	FName string
}

func getAllNewAssetsByPartnerID(partnerID string) (map[string]fnEndpoint, error) {

	partnerID = strings.TrimSpace(partnerID)
	mapRegIDEndpointID := PartnerEndpoints[partnerID]
	return mapRegIDEndpointID, nil
}

//PartnerMetrics is the map to store the logging info
type PartnerMetrics struct {
	ID                     string
	RMM1Err                string
	RMM2Err                string
	RMM1NoEndpointPartner  bool
	FNameInSyncEndpoints   []string
	RMM2NoEndpoints        []string
	UpdatedEndpoints       []string
	UpdateDBErrorEndpoints []string
}

func processParters(id int, jobs <-chan string, results chan<- PartnerMetrics) {
	var rmmoneerr, rmmtwoerr, unknwonptr, goodones int
	var rmm1tt, rmm2tt float64
	var rmm1cntr, rmm2cntr int
	for partnerID := range jobs {
		ls := PartnerMetrics{
			ID: partnerID,
		}
		rmm1ProcessingTime := time.Now()
		legacyAssetCollections, err := getAllLegacyAssetsByPartnerID(partnerID)
		rmm1cntr = rmm1cntr + 1
		rmm1tt = rmm1tt + time.Since(rmm1ProcessingTime).Seconds()
		if err != nil {
			//fmt.Println("RMM 1.0 : Error total  while getting the endpoints from Legacy for the partner Id : ", partnerID)
			ls.RMM1Err = err.Error()
			results <- ls
			rmmoneerr = rmmoneerr + 1
			continue
		}
		rmm2ProcessingTime := time.Now()
		// mapRegIDEndpointID, err := getAllNewAssetsByPartnerID(partnerID)

		partnerID = strings.TrimSpace(partnerID)
		mapRegIDEndpointID := PartnerEndpoints[partnerID]

		rmm2cntr = rmm2cntr + 1
		rmm2tt = rmm2tt + time.Since(rmm2ProcessingTime).Seconds()
		if err != nil {
			//fmt.Println("RMM 2.0 : Error total while getting the endpoints from Juno for the partner Id : total end", partnerID)
			ls.RMM2Err = err.Error()
			results <- ls
			rmmtwoerr = rmmtwoerr + 1
			continue
		}
		//fmt.Println("Good one Partner Id = ", partnerID)
		goodones = goodones + 1
		if len(legacyAssetCollections.EndPointList) <= 0 {
			ls.RMM1NoEndpointPartner = true
			unknwonptr = unknwonptr + 1
		}
		for regIdKey, value := range mapRegIDEndpointID {
			found := false
			for _, EndPointList := range legacyAssetCollections.EndPointList {
				legRegId := EndPointList.RegID
				if strings.TrimSpace(legRegId) == strings.TrimSpace(regIdKey) {
					fmt.Println("Data Issue partner Id,reg Id ", partnerID, "  ", regIdKey, "Friendly Name : ", value.FName)
					found = true
					break
				}
			}
			if !found {
				//fmt.Println("Data Issue partner Id,reg Id ", partnerID, "  ", regIdKey)
			}
		}
		results <- ls
	}
}

type asssetMigCollection struct {
	partnerID, EID, regId, FName string
}

func create(line []string) asssetMigCollection {
	return asssetMigCollection{
		EID:       line[3],
		regId:     line[2],
		partnerID: line[0],
	}
}

func readAssetTables(filePath string) ([]string, error) {
	fmt.Println("------------------readAssetTables from file path ----------------", filePath)
	partnerIDs := make([]string, 0)
	csvFile, ferr := os.Open(filePath)
	if ferr != nil {
		return partnerIDs, ferr
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var err error
	PartnerEndpoints = make(map[string]map[string]fnEndpoint)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			err = error
			log.Fatal(error)
		}
		astCol := create(line)

		partnerID := strings.TrimSpace(astCol.partnerID)
		if partnerID == "" {
			continue
		}
		regIDMap := PartnerEndpoints[partnerID]
		if regIDMap == nil {
			regIDMap = make(map[string]fnEndpoint)
		}
		rID := strings.TrimSpace(astCol.regId)

		endpoint := regIDMap[rID]
		if endpoint.EID == "" {
			endpoint := fnEndpoint{
				EID: astCol.EID,
			}
			regIDMap[rID] = endpoint
			PartnerEndpoints[partnerID] = regIDMap
		}
	}
	for key, _ := range PartnerEndpoints {
		partnerIDs = append(partnerIDs, key)
	}
	fmt.Println("-----------------readAssetTables as stored in memory partners size ----------------", len(PartnerEndpoints))
	fmt.Println("-----------------readAssetTables as partnerIDs size ----------------", len(partnerIDs))
	return partnerIDs, err
}
