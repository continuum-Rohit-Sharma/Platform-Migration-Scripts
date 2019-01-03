package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
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

//Session session variable
var Session *gocql.Session
var clusterIP string
var legacyWebAPIURL, newWebAPIURL string
var totalNoPartnersFromCSV, totalNoPartners int

const (
	queryToUpdateFriendlyName = "UPDATE platform_asset_db.partner_asset_collection set friendly_name=? where partner_id=? AND endpoint_id=?"
)

//func main1() {

//	fmt.Println("-----------------------------------------Friendly name migration reported started------")
//	startTime := time.Now()
//	fmt.Println("Migration Tool Started... Time started : ", startTime)
//	commandArgs := os.Args[1:]

//	serverHost := commandArgs[0]
//	endpointCsvName := commandArgs[1]
//	fmt.Println("---------------------------------serverHost ------", serverHost)
//	endpoints, err := readAssetTables(endpointCsvName)
//	if err != nil {
//		fmt.Println("Error Occured while getting the partners from file , Error : ", err)
//		return
//	}
//	fmt.Println("-----------------endpoints-size-------", len(endpoints))
//	url := ""
//	fmt.Println("----------------base url -------", url)

//}

type asssetMigCollection struct {
	partnerId, EID, regId, FName string
}

func create(line []string) asssetMigCollection {
	return asssetMigCollection{
		partnerId: line[3],
		EID:       line[0],
		regId:     line[1],
		FName:     line[2],
	}
}

func readAssetTables(filePath string) (map[string]map[string]asssetMigCollection, []string, error) {
	partnerAssetCollection := make(map[string]map[string]asssetMigCollection)
	partnerIDs := make([]string, 0)
	csvFile, ferr := os.Open(filePath)
	if ferr != nil {
		return partnerAssetCollection, partnerIDs, ferr
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var err error
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			err = error
			log.Fatal(error)
		}
		astCol := create(line)

		partnerId := strings.TrimSpace(astCol.partnerId)
		if partnerId == "" {
			continue
		}
		partnerIDs = append(partnerIDs, partnerId)
		regIdMap := partnerAssetCollection[partnerId]
		if regIdMap == nil {
			regIdMap = make(map[string]asssetMigCollection)
		}
		rId := strings.TrimSpace(astCol.regId)
		asssetMigCollection := regIdMap[rId]
		if asssetMigCollection.EID == "" {
			regIdMap[rId] = astCol
			partnerAssetCollection[partnerId] = regIdMap
		}
	}
	fmt.Println("------------------partnerAssetCollection----------------", partnerAssetCollection)
	return partnerAssetCollection, partnerIDs, err
}

func main() {

	commandArgs := os.Args[1:]
	if len(commandArgs) != 6 {
		fmt.Println("Usage: syncFriendlyName <csv File> <LegacyBaseURL> <RMM2.0BaseURL> <CassandraNodeIP> <# of workers> <resourceType>")
		fmt.Println("Example: syncFriendlyName partnerIDs.csv https://rmmapi.dtitsupport247.net http://internal-continuum-asset-service-elb-int-1972580147.ap-south-1.elb.amazonaws.com 172.28.48.6 10 desktop/server")
		return
	}
	startTime := time.Now()

	fmt.Println("Migration Tool Started... Time started : ", startTime)
	excelPath := commandArgs[0]
	legacyWebAPIURL = commandArgs[1]
	newWebAPIURL = commandArgs[2]
	clusterIP = commandArgs[3]
	noOfWorkers, errorParsing := strconv.Atoi(commandArgs[4])
	if errorParsing != nil {
		fmt.Println("Error.. while parsing the no. of workers.")
		return
	}
	resourceType := commandArgs[5]

	if legacyWebAPIURL != "" && resourceType != "" {
		legacyWebAPIURL = legacyWebAPIURL + "/itswebapi/v1/partner/%s/endpoints?resourceType=" + resourceType
		//fmt.Printf("Legacy Service URL : %s\n", legacyWebAPIURL)
	}

	if newWebAPIURL != "" {
		newWebAPIURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints?field=friendlyName"
		//fmt.Printf("New System Service URL %s\n", newWebAPIURL)
	}
	cassError := getCassandraSession(clusterIP)
	if cassError != nil {
		fmt.Println("Error Occured while setting up the Tool with Cassandra, Error : ", cassError)
		return
	}
	//	partnerIDs, err := ReadPartnersFromExcel(excelPath, "sheetName")
	//	if err != nil {
	//		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
	//		return
	//	}
	//	totalNoPartnersFromCSV = len(partnerIDs)
	//	//Getting Unique PartnerIds
	//	uniquePartnerIDs, err := removeDuplicatePartnerIds(partnerIDs)
	//	if err != nil {
	//		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
	//		return
	//	}

	//	//var a1 LoggingStruct[]
	//	totalNoPartnersUnique := len(uniquePartnerIDs)

	partnerAssetDBData, partnerIDs, asstErr := readAssetTables(excelPath)
	if asstErr != nil || len(partnerIDs) < 1 {
		fmt.Println("Error Occured while readAssetTables from Excel, partnerIDs :  , Error : ", partnerIDs, asstErr)
		return
	}
	totalNoPartnersUnique := len(partnerIDs)
	fmt.Println("partnerIDs : ", totalNoPartnersUnique)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan PartnerMetrics, totalNoPartnersUnique)

	for w := 1; w <= noOfWorkers; w++ {
		go processParters(w, jobs, results, partnerAssetDBData)
	}
	fmt.Println("# of workers : ", noOfWorkers)

	for a := 0; a < totalNoPartnersUnique; a++ {
		jobs <- partnerIDs[a]
	}
	fmt.Println("# of jobs ", totalNoPartnersUnique)
	close(jobs)

	fmt.Println("Waiting for results ")
	ptrRes := make([]PartnerMetrics, len(partnerIDs))
	var a int
	for a = 0; a < len(partnerIDs); a++ {
		ptrRes[a] = <-results
	}
	endTime := time.Since(startTime)
	fmt.Println("Done processing results. Total recieved : ", a)
	fmt.Println("Migration completed in : ", endTime)
	fmt.Printf("\n\n*** Partner Level Metrics ***\n\n")
	for a = 0; a < len(partnerIDs); a++ {
		fmt.Printf("Partner ID %s: RMM1Err: %s, RMM2Err: %s, NoEndpointsInRMM1: %v, NOE already in Sync: %d, NOE updated: %d, NOE with UpdateDBError: %d, NOE missing in RMM2: %d \n", ptrRes[a].ID, ptrRes[a].RMM1Err, ptrRes[a].RMM2Err, ptrRes[a].RMM1NoEndpointPartner, len(ptrRes[a].FNameInSyncEndpoints), len(ptrRes[a].UpdatedEndpoints), len(ptrRes[a].UpdateDBErrorEndpoints), len(ptrRes[a].RMM2NoEndpoints))
	}
	close(results)
	Session.Close()

}
func getAllLegacyAssetsByPartnerID(partnerID string) (assetCollection, error) {
	assetColl := assetCollection{}
	partnerID = strings.TrimSpace(partnerID)
	requestURL := fmt.Sprintf(legacyWebAPIURL, partnerID)
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
	mapRegIDEndpointID := make(map[string]fnEndpoint)
	partnerID = strings.TrimSpace(partnerID)
	requestURL := fmt.Sprintf(newWebAPIURL, partnerID)
	res, err := http.Get(requestURL)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 1.0 occured while getting the response for partner Id :", partnerID)
		return mapRegIDEndpointID, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 2.0 occured while Reading the response for partner Id :", partnerID)
		return mapRegIDEndpointID, err
	}

	bodyString := string(bodyBytes)
	assetColl := []asset.AssetCollection{}
	err = json.Unmarshal([]byte(bodyString), &assetColl)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 3.0 occured while  unmarshalling the response for partner Id :", partnerID)
		return mapRegIDEndpointID, err
	}

	for _, asset := range assetColl {
		regID := asset.RegID
		tmp := fnEndpoint{
			EID:   asset.EndpointID,
			FName: asset.FriendlyName,
		}
		if regID != "" {
			mapRegIDEndpointID[regID] = tmp
		}

	}

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

func processParters(id int, jobs <-chan string, results chan<- PartnerMetrics, assetDataMap map[string]map[string]asssetMigCollection) {
	var counter, rmmoneerr, rmmtwoerr, endpointNotFoundCounter, dbErr, unknwonptr, goodones, inSync int
	var rmm1tt, rmm2tt, dbupdatett float64
	var rmm1cntr, rmm2cntr, dbupdatecntr int
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
		//mapRegIDEndpointID, err := getAllNewAssetsByPartnerID(partnerID)
		mapRegIDAsset := assetDataMap[partnerID]
		rmm2cntr = rmm2cntr + 1
		rmm2tt = rmm2tt + time.Since(rmm2ProcessingTime).Seconds()
		if err != nil {
			//fmt.Println("RMM 2.0 : Error total while getting the endpoints from Juno for the partner Id : total end", partnerID)
			ls.RMM2Err = err.Error()
			results <- ls
			rmmtwoerr = rmmtwoerr + 1
			continue
		}
		fmt.Println("Good one Partner Id = ", partnerID)
		goodones = goodones + 1
		if len(legacyAssetCollections.EndPointList) <= 0 {
			ls.RMM1NoEndpointPartner = true
			unknwonptr = unknwonptr + 1
		}
		for _, EndPointList := range legacyAssetCollections.EndPointList {
			regID := EndPointList.RegID
			endpoint := mapRegIDAsset[regID]
			friendlyName := EndPointList.FriendlyName

			//endpoint := mapRegIDEndpointID[regID]

			if endpoint.EID != "" {
				if endpoint.FName == friendlyName {
					inSync = inSync + 1
					ls.FNameInSyncEndpoints = append(ls.FNameInSyncEndpoints, endpoint.EID)
					continue
				}
				dbUpdateProcessingTime := time.Now()
				errs := UpdateFriendLyNameByEndpointID(endpoint.EID, partnerID, friendlyName)
				dbupdatecntr = dbupdatecntr + 1
				dbupdatett = dbupdatett + time.Since(dbUpdateProcessingTime).Seconds()
				if errs != nil {
					//fmt.Println("Error while getting the endpoints from Juno for the partner Id : ", partnerID)
					dbErr = dbErr + 1
					ls.UpdateDBErrorEndpoints = append(ls.UpdateDBErrorEndpoints, endpoint.EID)

				} else {
					//fmt.Printf("Friendly Name updated. PartnerID : %s , Reg Id : %s , Endpoint ID : %s  and Friendly Name : %s\n", partnerID, regID, endpointID, friendlyName)
					counter = counter + 1
					ls.UpdatedEndpoints = append(ls.UpdatedEndpoints, endpoint.EID)
				}

			} else {
				//fmt.Printf("Endpoint id  not found for PartnerID : %s , reg Id : %s and friendly Name : %s\n", partnerID, regID, friendlyName)
				endpointNotFoundCounter = endpointNotFoundCounter + 1
				ls.RMM2NoEndpoints = append(ls.RMM2NoEndpoints, endpoint.EID)

			}
		}

		results <- ls
	}
	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM1.0 errors: %d, RMM2.0 errors: %d, GoodCandidates: %d, Partners with no endpoints in RMM 1.0: %d, Common Partners with endpoints: %d\n", id, rmmoneerr, rmmtwoerr, goodones, unknwonptr, (goodones - unknwonptr))
	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM1.0 GetLegacyEndpointlist total time :  %f , instances: %d,  avg per partner: %f\n", id, rmm1tt, rmm1cntr, rmm1tt/float64(rmm1cntr))
	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM2.0 GetJunoEndpointlist   total time :  %f , instances: %d,  avg per partner: %f\n", id, rmm2tt, rmm2cntr, rmm2tt/float64(rmm2cntr))

	fmt.Printf("Goroutine id# %d, PerThreadEndpointMetrics : Already in Sync: %d, DB Error : %d, EndpointNotfound : %d,  processed : %d \n", id, inSync, dbErr, endpointNotFoundCounter, counter)
	fmt.Printf("Goroutine id# %d, PerThreadEndpointMetrics : DBupdate total time :  %f , instances: %d,  avg per partner: %f\n", id, dbupdatett, dbupdatecntr, dbupdatett/float64(dbupdatecntr))
}

//UpdateFriendLyNameByEndpointID is to update the friendlyName by endpoint_id and partner_Ids
func UpdateFriendLyNameByEndpointID(endpointID, partnerID, friendlyName string) error {
	err := Session.Query(queryToUpdateFriendlyName,
		friendlyName,
		partnerID,
		endpointID).Exec()
	return err
}
func getCassandraSession(clusterIP string) error {
	var err error
	cluster := gocql.NewCluster(clusterIP)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "platform_asset_db"
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Printf("Error occured : %v", err)
		return err
	}
	return nil
}

//ReadPartnersFromExcel ...
func ReadPartnersFromExcel(excelPath string, sheetName string) ([]string, error) {

	var partners [][]string
	var partnerIDs []string
	file, err := os.Open(excelPath)
	if err != nil {
		fmt.Printf("Error while Reading the file : %s and error is : %v\n", excelPath, err)
		fmt.Println(err)
		os.Exit(1)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	partners, err = reader.ReadAll()
	length := len(partners)

	if length > 0 {
		partnerIDs = make([]string, length)
	}

	for cnt := 0; cnt < length; cnt++ {
		partnerIDs[cnt] = partners[cnt][0]
	}
	if len(partnerIDs) > 0 {
		fmt.Println("Partners with no errors: ", len(partnerIDs))
	} else {
		fmt.Println("No partners found in the excrel so returning nil")
	}

	return partnerIDs, nil
}

func removeDuplicatePartnerIds(partnerIDs []string) ([]string, error) {
	var err error
	if len(partnerIDs) < 1 {
		return partnerIDs, err
	}

	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range partnerIDs {
		if encountered[partnerIDs[v]] == true {
		} else {
			encountered[partnerIDs[v]] = true
			partnerID := strings.TrimSpace(partnerIDs[v])
			result = append(result, partnerID)
		}
	}
	return result, nil
}
