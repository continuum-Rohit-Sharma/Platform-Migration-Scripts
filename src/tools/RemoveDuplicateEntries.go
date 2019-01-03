package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/gocql/gocql"
)

type assetCollection struct {
	EndPointList []EndPointList `json:"endPointList"`
	TotalPoint   int            `json:"totalCount"`
}

//Config Struct to hold Configuration for Load Generation
type Config struct {
	PrintResult               bool
	RequestWaitInSecond       int64
	ServerAddress             string
	NewWebAPIURL              string
	OldWebAPIURL              string
	IsDelete                  bool
	IsSite                    bool
	SiteIDs                   string
	NumOfWorkers              int
	CassandraHostIPForAsset   string
	CassandraHostIPForAgent   string
	CSVPathForPartners        string
	PickUninstallListFromFile bool
	ThresholdDays             int
	PartnerIds                []string
}

//EndPointList is the structuire
type EndPointList struct {
	RegID                 string `json:"regID,omitempty"`
	FriendlyName          string `json:"friendlyName"`
	MachineName           string `json:"machineName"`
	SiteName              string `json:"siteName"`
	OperatingSystem       string `json:"operatingSystem"`
	Availability          int    `json:"availability"`
	IPAddress             string `json:"ipAddress"`
	RegType               string `json:"regType"`
	LmiStatus             int    `json:"lmiStatus"`
	ResType               string `json:"resType"`
	SiteID                int    `json:"siteId"`
	SmartDisk             int    `json:"smartDisk"`
	Amt                   int    `json:"amt"`
	MbSyncstatus          int    `json:"mbSyncstatus"`
	EncryptedResourceName string `json:"encryptedResourceName"`
	EncryptedSiteName     string `json:"encryptedSiteName"`
	LmiHostID             string `json:"lmiHostId"`
	RequestRegID          string `json:"requestRegId"`
	EndpointID            string `json:"endpointId"`
}
type fnEndpoint struct {
	PartnerID     string
	EndpointID    string `json:"endpointId"`
	SystemName    string `json:"systemName"`
	BiosSerialNum string `json:"serialNumber"`
	RegID         string `json:"regID"`
}

type resultValue struct {
	partnerID string
	endpoints [][]string
}

//Session session variable
var configObject Config

//AssetSession is for gocql asset session
var AssetSession *gocql.Session

//AgentSession is for gocql asset session
var AgentSession *gocql.Session
var clusterIP string
var newWebAPIURL, getAssetURL string
var totalNoPartnersFromCSV, totalNoPartners int
var mapSystemNameFnEndpoint map[string][]fnEndpoint
var duplicateEndpoints [][]string
var mapDuplicateEndpoints map[string][][]string
var doDelete bool
var totalCount, totalPartners int

const (
	queryToDeleteEndpoints       = `DELETE FROM platform_asset_db.partner_asset_collection where partner_id= '%s' AND endpoint_id in (%s)`
	cSelectMultipleEndpointQuery = `select partner_id, endpoint_id,  DcDateTimeUTC from platform_agent_db.agent_heartbeat where partner_id ='%s' and endpoint_id in (%s)`
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
	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: RemoveDuplicates config.json")
		return
	}
	startTime := time.Now()
	fmt.Println("Tool Started... Time started : ", startTime)

	configObject, err := LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}
	excelPath := configObject.CSVPathForPartners
	newWebAPIURL = configObject.NewWebAPIURL
	assetClusterIP := configObject.CassandraHostIPForAsset
	agentClusterIP := configObject.CassandraHostIPForAgent
	noOfWorkers := configObject.NumOfWorkers
	doDelete = configObject.IsDelete

	if newWebAPIURL != "" {
		getAssetURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints/%s"
		newWebAPIURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints?field=bios&field=system"
		//fmt.Printf("New System Service URL %s\n", newWebAPIURL)
	}
	cassError := getCassandraSession(assetClusterIP, agentClusterIP)
	if cassError != nil {
		fmt.Println("Error Occured while setting up the Tool with Cassandra for Asset, Error : ", cassError)
		return
	}

	partnerIDs, err := ReadPartnersFromExcel(excelPath, "sheetName")
	if err != nil {
		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
		return
	}
	totalNoPartnersFromCSV = len(partnerIDs)
	//Getting Unique PartnerIds
	uniquePartnerIDs, err := removeDuplicatePartnerIds(partnerIDs)
	if err != nil {
		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
		return
	}
	totalNoPartnersUnique := len(uniquePartnerIDs)
	fmt.Println("# of partners :", totalNoPartnersUnique)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan resultValue, totalNoPartnersUnique)
	mapDuplicateEndpoints = make(map[string][][]string)
	for w := 1; w <= noOfWorkers; w++ {
		go processParters(w, jobs, results)
	}
	fmt.Println("# of workers : ", noOfWorkers)
	for a := 0; a < len(uniquePartnerIDs); a++ {
		jobs <- uniquePartnerIDs[a]
	}
	fmt.Println("# of jobs ", totalNoPartnersUnique)
	close(jobs)
	fmt.Println("Waiting for results ")
	var a int

	for a = 0; a < len(uniquePartnerIDs); a++ {
		duplicateResultValue := <-results
		partnerID := duplicateResultValue.partnerID
		endpoints := duplicateResultValue.endpoints
		if len(endpoints) > 0 {
			mapDuplicateEndpoints[partnerID] = endpoints
			fmt.Printf("\nDuplicate Results found for partner Id %s, values are : %v\n", partnerID, mapDuplicateEndpoints[partnerID])
		}
	}
	endTime := time.Since(startTime)
	fmt.Println("Done processing results. Total recieved : ", a)
	fmt.Println("Identification completed in : ", endTime)
	fmt.Printf("\n\n*** Duplicate endpoints Metrics ***\n")
	for partnerID, endpointIDs := range mapDuplicateEndpoints {
		fmt.Printf("\n Partner ID %s, || Endpoint Ids : %v  \n\n", partnerID, endpointIDs)
	}
	fmt.Println("Is Delete option selected : ", doDelete)

	identifyStaleEntriesAndDelete(mapDuplicateEndpoints, doDelete)
	close(results)
	AssetSession.Close()
	AgentSession.Close()
}
func getAllNewAssetsByPartnerID(partnerID string) ([]fnEndpoint, error) {
	//mapRegIDEndpointID := make(map[string]string)
	assetCollections := []fnEndpoint{}
	partnerID = strings.TrimSpace(partnerID)
	requestURL := fmt.Sprintf(newWebAPIURL, partnerID)
	res, err := http.Get(requestURL)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 1.0 occured while getting the response for partner Id :", partnerID)
		return assetCollections, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 2.0 occured while Reading the response for partner Id :", partnerID)
		return assetCollections, err
	}
	bodyString := string(bodyBytes)
	assetColl := []asset.AssetCollection{}
	err = json.Unmarshal([]byte(bodyString), &assetColl)
	if err != nil {
		//fmt.Println("RMM 2.0 : Error 3.0 occured while  unmarshalling the response for partner Id :", partnerID)
	}
	for _, asset := range assetColl {
		//regID := asset.RegID
		tmp := fnEndpoint{
			EndpointID:    asset.EndpointID,
			BiosSerialNum: asset.Bios.SerialNumber,
			PartnerID:     partnerID,
			RegID:         asset.RegID,
			SystemName:    asset.System.SystemName,
		}
		assetCollections = append(assetCollections, tmp)
	}
	//fmt.Println(assetCollections)
	return assetCollections, nil
}
func processParters(id int, jobs <-chan string, results chan<- resultValue) {
	mapSystemNameFnEndpoint := make(map[string][]fnEndpoint)

	for partnerID := range jobs {
		result := resultValue{partnerID: partnerID}
		mapDupEndpoints := make(map[string][][]string)
		//fmt.Printf("/nWorker %d took partner Id : %s for processing /n", id, partnerID)

		//assetCollections, err := test(partnerID)
		//fmt.Println(assetCollections)
		assetCollections, err := getAllNewAssetsByPartnerID(partnerID)
		if err == nil {
			if len(assetCollections) > 0 {
				//duplicateEndpoints := make([][]string, len(assetCollections))
				// mapSystemNameFnEndpoint contains key = system name and value = fnendpoint(endpointid,biosserialnum,partnerid,regId,systemname)
				mapSystemNameFnEndpoint = filterOnMachineName(assetCollections)

				// mapSystemNameFnEndpointRemoveRegId := make(map[string][]fnEndpoint)
				mapSystemNameFnEndpoint, mapDupEndpoints = filterOnRegID(mapSystemNameFnEndpoint)
				//fmt.Println("mapSystemNameFnEndpointRemoveRegId : ",mapSystemNameFnEndpointRemoveRegId)
				//fmt.Println("before mapDupEndpoints : ",mapDupEndpoints)
				mapDupEndpoints = filterOnBIOSSerialNum(mapSystemNameFnEndpoint, mapDupEndpoints)
				//fmt.Println("after mapDupEndpoints : ",mapDupEndpoints)
			}
			if len(mapDupEndpoints[partnerID]) > 0 {
				result = resultValue{
					partnerID: partnerID,
					endpoints: mapDupEndpoints[partnerID],
				}
			}
			//fmt.Printf("/nfor Partner %s, result is %v /n", partnerID, result)
			results <- result
		} else {
			fmt.Println("Error Occured while getting the endpoints for partner Id :", partnerID)
		}
	}
}

func test(partnerID string) ([]fnEndpoint, error) {
	endpoints := make([]fnEndpoint, 0, 5)
	if partnerID == "9394" {
		endpoints = append(endpoints, fnEndpoint{"9394", "99441d2e-13a5-47fb-92d1-4e07b9d55fcd", "sys1", "sr1", "r1"})
		endpoints = append(endpoints, fnEndpoint{"9394", "e2", "sys1", "sr2", "r1"})
		endpoints = append(endpoints, fnEndpoint{"9394", "eed02908-1b50-4ccb-96b3-52c2a6fdb622", "sys1", "sr1", "r1"})
		endpoints = append(endpoints, fnEndpoint{"9394", "e4", "sys1", "sr2", "r2"})
		return endpoints, nil
	}
	if partnerID == "2" {
		endpoints = append(endpoints, fnEndpoint{"2", "e5", "sys2", "sr2", "r2"})
		endpoints = append(endpoints, fnEndpoint{"2", "e6", "sys2", "sr3", "r3"})
		return endpoints, nil
	}
	if partnerID == "3" {
		endpoints = append(endpoints, fnEndpoint{"3", "e7", "sys3", "sr22", "r2"})
		endpoints = append(endpoints, fnEndpoint{"3", "e8", "sys3", "sr22", "r2"})
		return endpoints, nil
	}
	if partnerID == "50001743" {
		endpoints = append(endpoints, fnEndpoint{"50001743", "37813701-826b-4687-95a4-7e74da44b477", "Juno-win7-65Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 0d", "53757935"})
		endpoints = append(endpoints, fnEndpoint{"50001743", "37813701-826b-4687-95a4-7e74da44b478", "Juno-win7-65Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 1d", "53757935"})
		endpoints = append(endpoints, fnEndpoint{"50001743", "37813701-826b-4687-95a4-7e74da44b479", "Juno-win7-65Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 0d", "53757937"})
		endpoints = append(endpoints, fnEndpoint{"50001743", "37813701-826b-4687-95a4-7e74da44b480", "Juno-win7-65Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 1d", "53757937"})
		return endpoints, nil
	}
	if partnerID == "50001794" {
		endpoints = append(endpoints, fnEndpoint{"50001794", "37813701-826b-4687-95a4-7e74da44b485", "Juno-win7-651Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 3d", "53757939"})
		endpoints = append(endpoints, fnEndpoint{"50001794", "7813701-826b-4687-95a4-7e74da44b4890", "Juno-win7-651Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 4d", "53757939"})
		endpoints = append(endpoints, fnEndpoint{"50001794", "37813701-826b-4687-95a4-7e74da44b487", "Juno-win7-651Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 3d", "53757940"})
		endpoints = append(endpoints, fnEndpoint{"50001794", "37813701-826b-4687-95a4-7e74da44b488", "Juno-win7-651Rohit", "VMware-42 27 a8 5d b3 4d ec 6c-35 6e af 47 66 ab 9b 4d", "53757940"})
		return endpoints, nil
	}

	return endpoints, nil

}
func filterOnBIOSSerialNum(mapSystemNameFnEndpoint map[string][]fnEndpoint, mapDupEndPt map[string][][]string) map[string][][]string {
	var partnerID string
	for _, EndPointList := range mapSystemNameFnEndpoint {
		mapBiosSerialNameEndpt := make(map[string][]string)
		if len(EndPointList) <= 1 {
			continue
		}
		for _, assetDtl := range EndPointList {
			partnerID = assetDtl.PartnerID
			biosSerialNum := assetDtl.BiosSerialNum
			endpointID := assetDtl.EndpointID
			if biosSerialNum == "" {
				fmt.Printf("\n Blank Bios serialnum for endpoint : %v\n", endpointID)
				continue
			}
			mapBiosSerialNameEndpt[biosSerialNum] = append(mapBiosSerialNameEndpt[biosSerialNum], endpointID)

		}
		for _, endpointIDs := range mapBiosSerialNameEndpt {
			if len(endpointIDs) > 1 {
				mapDupEndPt[partnerID] = append(mapDupEndPt[partnerID], endpointIDs)
			}
		}
	}
	return mapDupEndPt
}
func filterOnRegID(mapSystemNameFnEndpoint map[string][]fnEndpoint) (map[string][]fnEndpoint, map[string][][]string) {
	mapRegIDPartnerID := make(map[string]string)
	mapRegIDAndEndptIds := make(map[string][]string)
	var partnerID string
	mapSysNameToRemove := make(map[string]bool)
	mapDupEndpoints := make(map[string][][]string)
	mapRegIdSysName := make(map[string]string)
	for systemName, assetsWithSameSysName := range mapSystemNameFnEndpoint {

		lenAssetWithSysName := len(assetsWithSameSysName)

		if lenAssetWithSysName > 1 {
			for _, fnEndpoint := range assetsWithSameSysName {
				partnerID = fnEndpoint.PartnerID
				regID := fnEndpoint.RegID
				endpointID := fnEndpoint.EndpointID
				mapRegIDAndEndptIds[regID] = append(mapRegIDAndEndptIds[regID], endpointID)
				mapRegIdSysName[regID] = systemName
				mapRegIDPartnerID[regID] = partnerID
			}
		}
	}

	for regID, endpointIDs := range mapRegIDAndEndptIds {
		if regID == "" {
			fmt.Printf("\n Blank reg id for endpoint : %v\n", endpointIDs)
			continue
		}
		systemName := mapRegIdSysName[regID]
		if systemName == "" {
			fmt.Printf("\n Blank systemname for endpoint : %v\n", endpointIDs)
			continue
		}
		lengthEndpoints := len(endpointIDs)
		if lengthEndpoints > 1 {
			assetsWithSameSysName := mapSystemNameFnEndpoint[systemName]
			partnerID = mapRegIDPartnerID[regID]
			lengthAssets := len(assetsWithSameSysName)
			if lengthEndpoints == lengthAssets {
				//mapDupEndpoints is the final map which will be consider for deletion
				mapDupEndpoints[partnerID] = append(mapDupEndpoints[partnerID], endpointIDs)

			}
		}
	}

	return mapSystemNameFnEndpoint, mapDupEndpoints
}

func filterOnMachineName(assetCollections []fnEndpoint) map[string][]fnEndpoint {
	mapSystemNameFnEndpoint := make(map[string][]fnEndpoint)

	for _, fnEndpoint := range assetCollections {
		systemName := fnEndpoint.SystemName
		mapSystemNameFnEndpoint[systemName] = append(mapSystemNameFnEndpoint[systemName], fnEndpoint)
	}

	return mapSystemNameFnEndpoint
}
func appendStructInMapOfArray(assets []fnEndpoint, insertionIndex int, key string, endpoint fnEndpoint) []fnEndpoint {
	newDupEndpoints := make([]fnEndpoint, insertionIndex)
	if insertionIndex == 0 {
		newDupEndpoints[0] = endpoint
		return newDupEndpoints
	}
	for cnt, endpoint := range assets {
		newDupEndpoints[cnt] = endpoint
	}
	newDupEndpoints[insertionIndex-1] = endpoint
	return newDupEndpoints
}

func getCassandraSession(assetClusterIP, agentClusterIP string) error {
	var err error
	assetCluster := gocql.NewCluster(assetClusterIP)
	assetCluster.Consistency = gocql.Quorum
	assetCluster.Keyspace = "platform_asset_db"
	AssetSession, err = assetCluster.CreateSession()
	if err != nil {
		fmt.Printf("Cassandra Error occured for Asset: %v", err)
		return err
	}
	agentCluster := gocql.NewCluster(agentClusterIP)
	agentCluster.Consistency = gocql.Quorum
	agentCluster.Keyspace = "platform_agent_db"
	AgentSession, err = agentCluster.CreateSession()
	if err != nil {
		fmt.Printf("Cassandra Error occured for Agent: %v", err)
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

type hbEndpoint struct {
	EndpointID    string
	PartnerID     string
	DcDateTimeUTC time.Time
}

//Heartbeat struct for capturing heartbeat
type Heartbeat struct {
	RegID            string
	DcDateTimeUTC    time.Time
	AgentDateTimeUTC time.Time
	HeartbeatCounter int64
	EndpointID       string
	installed        bool
}

func identifyStaleEntriesAndDelete(canDulicateEntries map[string][][]string, doDelete bool) {
	//fmt.Println("Came to identify Stale entries and Delete", canDulicateEntries)
	totalCount = 0
	totalPartners = 0
	fmt.Println("*********************Stale Entries **************************")
	var partner_id, endpoint_id string
	var DcDateTimeUTC time.Time
	var flag bool
	var maxtime time.Time
	var end string
	var dupListEndpoints []string

	for partnerID, endpoints := range canDulicateEntries {
		dupListEndpoints = make([]string, 0)
		for _, endpointIDs := range endpoints {
			sEndpointIDs := strings.Join(endpointIDs, "','")
			sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
			query := fmt.Sprintf(cSelectMultipleEndpointQuery, partnerID, sEndpointIDs)
			iter := AgentSession.Query(query).Iter()
			var heartbeats []Heartbeat
			flag = false
			for iter.Scan(&partner_id, &endpoint_id, &DcDateTimeUTC) {
				tmp := Heartbeat{
					EndpointID:    endpoint_id,
					DcDateTimeUTC: DcDateTimeUTC,
				}
				if !flag {
					maxtime = DcDateTimeUTC
					end = endpoint_id
					flag = true
				} else {
					if DcDateTimeUTC.After(maxtime) {
						dupListEndpoints = append(dupListEndpoints, end)
						maxtime = DcDateTimeUTC
						end = endpoint_id
					} else {
						dupListEndpoints = append(dupListEndpoints, endpoint_id)
					}
				}
				heartbeats = append(heartbeats, tmp)
			}
			if err := iter.Close(); err != nil {
				log.Fatal(err)
			}
		}
		lengthOfStale := len(dupListEndpoints)
		if lengthOfStale > 0 {
			fmt.Printf("\nPartner Id %s, Total %d Entries were identified for Deletion and they are : %v \n", partnerID, lengthOfStale, dupListEndpoints)
			totalCount = totalCount + lengthOfStale
			totalPartners = totalPartners + 1
		}

		if doDelete {
			if lengthOfStale < 1 {
				fmt.Printf("Nothing to Delete for partner Id hence returning")
				continue
			}
			jobs := make(chan string, lengthOfStale)
			results := make(chan asset.AssetCollection, lengthOfStale)
			for w := 1; w <= lengthOfStale; w++ {
				go getEndpointByEndpointID(w, jobs, results, partnerID)
			}
			for a := 0; a < lengthOfStale; a++ {
				jobs <- dupListEndpoints[a]
			}
			close(jobs)
			var a int
			for a = 0; a < lengthOfStale; a++ {
				asset := <-results
				fmt.Printf("\nFor Partner Id %s and Endpoint Id %s, JSON BackUP is : %v\n", partnerID, asset.EndpointID, asset)
			}
			close(results)
			err := deleteDupEndpoints(dupListEndpoints, partnerID)
			if err != nil {
				fmt.Println("Error Occured while deleting the duplicate entries for Partner Id : ", partnerID)
			}
		}
	}
	fmt.Printf("\n\n Total # Partners Having Duplicates are : %d and Total duplicate endpoint Entries are : %d \n", totalPartners, totalCount)
}

func getEndpointByEndpointID(id int, jobs <-chan string, results chan<- asset.AssetCollection, partnerID string) {
	for endpointID := range jobs {
		assetColl := asset.AssetCollection{}
		partnerID = strings.TrimSpace(partnerID)
		endpointID = strings.TrimSpace(endpointID)
		requestURL := fmt.Sprintf(getAssetURL, partnerID, endpointID)
		res, err := http.Get(requestURL)
		if err != nil {
			results <- assetColl
		}
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			results <- assetColl
		}
		bodyString := string(bodyBytes)

		err = json.Unmarshal([]byte(bodyString), &assetColl)
		if err != nil {
			fmt.Println("RMM 2.0 : Error 3.0 occured while  unmarshalling the response for partner Id :", partnerID)
		}
		results <- assetColl
	}

}
func deleteDupEndpoints(endpointIDs []string, partnerID string) error {
	fmt.Println("Deleting the Assets for Endpoint ids : ")
	sEndpointIDs := strings.Join(endpointIDs, "','")
	sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
	query := fmt.Sprintf(queryToDeleteEndpoints, partnerID, sEndpointIDs)
	err := AssetSession.Query(query).Exec()
	return err
}
