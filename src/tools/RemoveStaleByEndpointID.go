package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/gocql/gocql"
	"github.com/xuri/excelize"
)

type assetCollection struct {
	EndPointList []EndPointList `json:"endPointList"`
	TotalPoint   int            `json:"totalCount"`
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
	PartnerID      string
	EndpointID     string `json:"endpointId"`
	RegID          string `json:"regID"`
	ClientID       string
	siteID         string
	HeartBeartFlag bool
	installedFlag  bool
	AgentID        string
}

//LegacyAssets is the structure to get the 1.0 asset datas
type LegacyAssets struct {
	Collection []LegacyAsset
}

//LegacyAsset is the structure to get the 1.0 asset data
type LegacyAsset struct {
	PartnerID int `json:"memberID"`
	SiteID    int `json:"siteID"`
	RegID     int `json:"regID"`
}

type resultValue struct {
	partnerID string
	endpoints []string
}

var (
	configObject                                  Config
	logger                                        *log.Logger
	endpoints                                     []string
	mapPartnerIDRegIDRemoved, mapPartnerIDSiteIDs map[string][]string
	//AssetSession session variable
	AssetSession *gocql.Session
	//AgentSession session variable
	AgentSession                            *gocql.Session
	clusterIP                               string
	newWebAPIURL, getAssetURL, oldAssetURL  string
	lengthOfTotalPartnerID, totalNoPartners int

	mapPartIDHbTrueFalseAndInstTrue, mapPartIDHbTrueAndInstFalse map[string][]fnEndpoint
	mapPartIDHbFalseAndInstFalse                                 map[string][]fnEndpoint
	doDelete                                                     bool
	totalCount, totalPartners, noOfWorkers                       int
)

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
	CassandraHost             string
	CSVPathForPartners        string
	PickUninstallListFromFile bool
	ThresholdDays             int
	PartnerIds                []string
	Messages                  []MailboxMessage
	CassandraHostIPForAsset   string
	CassandraHostIPForAgent   string
}

type Mailbox struct {
	Endpoints []string
	SourceURL string
	Message   MailboxMessage
}

type MailboxMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Path         string    `json:"path"`
	Message      string    `json:"message"`
}

const (
	queryToGetPartners               = `select partner_id from platform_asset_db.partner_asset_collection`
	queryToDeleteEndpoints           = `DELETE FROM platform_asset_db.partner_asset_collection where partner_id= '%s' AND endpoint_id in (%s)`
	queryToDeleteAgentVersionDetails = `DELETE FROM platform_asset_db.agent_version_details where partner_id= '%s' AND endpoint_id in (%s)`
	cSelectMultipleEndpointQuery     = `select partner_id, endpoint_id,  DcDateTimeUTC from platform_agent_db.agent_heartbeat where partner_id ='%s' and endpoint_id in (%s)`
	queryToGetHeartBeatByPartnerID   = `select partner_id, endpoint_id, dcdatetimeutc from platform_agent_db.agent_heartbeat where partner_id = '%s'`
	queryToPopulateInstalledFlag     = `select endpoint_id, agent_id, client_id, installed, legacy_reg_id, partner_id, site_id from platform_agent_db.partnerendpointmap where endpoint_id in (%s)`
	queryToInsertAssetInActive       = `INSERT INTO platform_asset_db.asset_partner_endpoint_details_inactive(partner_id,client_id,site_id,endpoint_id,reg_id,dc_ts_utc,agent_ts_utc,created_by,name,type,remote_address,resource_type,endpoint_type,baseboard,bios,drives,physical_memory,networks,os,processors,raidcontroller,system, installed_softwares, keyboards, mouse, monitors, physical_drives,users,services,shares) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	cUpdatePEMInstallFlgQuery        = `update platform_agent_db.partnerendpointmap set installed = %t where endpoint_id in (%s)`
)

func main() {
	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: RemoveStaleEntriesByExcel config.json")
		return
	}
	startTime := time.Now()
	fmt.Println("Tool Started... Time started : ", startTime)

	logfile, err := os.OpenFile("UtilityLog.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	logger = log.New(logfile, "", logFlags)
	configObject, err = LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}
	excelPath := configObject.CSVPathForPartners
	newWebAPIURL = configObject.NewWebAPIURL
	//oldAssetURL = configObject.OldWebAPIURL
	assetClusterIP := configObject.CassandraHostIPForAsset
	agentClusterIP := configObject.CassandraHostIPForAgent
	noOfWorkers = configObject.NumOfWorkers
	doDelete = configObject.IsDelete

	isSite := configObject.IsSite

	cassError := getCassandraSession(assetClusterIP, agentClusterIP)
	if cassError != nil {
		fmt.Println("Error Occured while setting up the Tool with Cassandra, Error : ", cassError)
		return
	}

	//https://rmmapi.dtitsupport247.net/itswebapi/v1/GetEndpointsUninstalled/NumofDays/MemberID

	mapPartnerIDListEndpoints, err := ReadPartnersFromExcel(excelPath, "sheetName")
	if err != nil {
		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
		return
	}

	if newWebAPIURL != "" {
		if isSite {
			newWebAPIURL = newWebAPIURL + "/asset/v1/partner/%s/sites/%s/summary"
		} else {
			//getAssetURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints/%s"
			newWebAPIURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints?field=bios&field=system"
		}

		// 	//fmt.Printf("New System Service URL %s\n", newWebAPIURL)
	}

	// //To get the List of all juno2.0 partner from CSV

	fmt.Println("\n****** Identified Entries START **********")
	fmt.Println("Partner ID | Total no. entries | Endpoint IDs")
	for partID, endpoints := range mapPartnerIDListEndpoints {
		fmt.Printf("\n  %s | %d | %v \n", partID, len(endpoints), endpoints)
	}
	fmt.Println("\n****** Identified Entries END  **********")
	if doDelete {
		processDelete(mapPartnerIDListEndpoints)
	}

	AssetSession.Close()
	AgentSession.Close()
}

func readEndpointsByAPI(partnerIds []string, thresholdDays int) (map[string][]string, []string, map[string][]string, error) {
	fmt.Printf("\n \n Came readEndpointsByAPI for partner Ids : (%v) \n \n", partnerIds)
	mapPartnerIDRegID := make(map[string][]string)
	mapPartnerIDSiteIDs := make(map[string][]string)
	partnerIDs := make([]string, 0)
	//legacyAssets := &LegacyAssets{}
	if len(partnerIds) > 0 {
		var err error
		for _, partnerID := range partnerIds {
			partnerID = strings.TrimSpace(partnerID)
			if thresholdDays <= 0 {
				thresholdDays = 7
			}

			mapPartnerIDRegID, mapPartnerIDSiteIDs, err = requestLegacyDataForPartnerID(partnerID, thresholdDays, mapPartnerIDRegID)
			if err != nil {
				fmt.Printf("\n \n Error Occurred while getting the Uninstalled Data from 1.0 for PartnerID :%s \n \n", partnerID)
				return mapPartnerIDRegID, partnerIDs, mapPartnerIDSiteIDs, err
			}
		}
	} else {
		partnerID := "-1"
		var err error
		mapPartnerIDRegID, mapPartnerIDSiteIDs, err = requestLegacyDataForPartnerID(partnerID, thresholdDays, mapPartnerIDRegID)
		if err != nil {
			fmt.Printf(" \n \n Error Occurred while getting the Uninstalled Data from 1.0 for PartnerID :%s", partnerID)
			return mapPartnerIDRegID, partnerIDs, mapPartnerIDSiteIDs, err
		}
	}
	RegIDS := make([][]string, 0)
	for partnerID, RegID := range mapPartnerIDRegID {
		partnerIDs = append(partnerIDs, partnerID)
		RegIDS = append(RegIDS, RegID)
	}
	fmt.Println("Partner IDS obtained :", partnerIDs)
	fmt.Println("Uninstalled Data Obtained :", mapPartnerIDRegID)
	return mapPartnerIDRegID, partnerIDs, mapPartnerIDSiteIDs, nil
}

func requestLegacyDataForPartnerID(partnerID string, thresholdDays int, mapPartnerIDRegID map[string][]string) (map[string][]string, map[string][]string, error) {
	legacyAssets := make([]LegacyAsset, 0)
	mapPartnerIDSiteIDs := make(map[string][]string)
	requestURL := fmt.Sprintf(oldAssetURL, thresholdDays, partnerID)
	//requestURL := fmt.Sprintf(requestURL, "")
	fmt.Println("About to get the Asset by URL : ", requestURL)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("RMM 1.0 : Error occurred while getting the response for partner Id :", partnerID)
		return mapPartnerIDRegID, mapPartnerIDSiteIDs, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("RMM 1.0 : Error occurred while Reading the response for partner Id :", partnerID)
		return mapPartnerIDRegID, mapPartnerIDSiteIDs, err
	}
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &legacyAssets)
	if err != nil {
		fmt.Printf("\n \n RMM 2.0 : Error occured while  unmarshalling the response for partner Id : %s, And error : %v", partnerID, err)
		return mapPartnerIDRegID, mapPartnerIDSiteIDs, err
	}
	for _, LegacyAsset := range legacyAssets {
		regID := strconv.Itoa(LegacyAsset.RegID)
		partnerID := strconv.Itoa(LegacyAsset.PartnerID)
		siteID := strconv.Itoa(LegacyAsset.SiteID)
		mapPartnerIDRegID[partnerID] = append(mapPartnerIDRegID[partnerID], regID)
		mapPartnerIDSiteIDs[partnerID] = append(mapPartnerIDSiteIDs[partnerID], siteID)
	}
	return mapPartnerIDRegID, mapPartnerIDSiteIDs, nil
}

func processDeletionAndPostMailBoxMessages(mapPartnerIDListEndpoints map[string][]string, mapPartIDHbTrueFalseAndInstTrue, mapPartIDHbTrueAndInstFalse, mapPartIDHbFalseAndInstFalse map[string][]fnEndpoint, doDelete bool) {
	mapEndptIdsToInstallFalse := make(map[string][]string)
	mapEndptIdsToInstallTrue := make(map[string][]string)
	for partnerID, endpoints := range mapPartIDHbTrueFalseAndInstTrue {
		endpointIDs := mapPartnerIDListEndpoints[partnerID]
		if len(endpointIDs) > 0 {
			mapEndptIdsToInstallFalse = returnMapToUpdateFlagsIfMatches(partnerID, endpoints, endpointIDs)
		}
	}

	for partnerID, endpoints := range mapPartIDHbTrueAndInstFalse {
		endpointIDs := mapPartnerIDListEndpoints[partnerID]
		if len(endpointIDs) > 0 {
			mapEndptIdsToInstallTrue = returnMapToUpdateFlagsIfNotMatches(partnerID, endpoints, endpointIDs)
		}
	}

	if doDelete {
		preProcessMailBox(mapPartnerIDListEndpoints)
	}
	for partnerID, endpoints := range mapPartIDHbFalseAndInstFalse {
		for _, fnEndpoint := range endpoints {
			mapPartnerIDListEndpoints[partnerID] = append(mapPartnerIDListEndpoints[partnerID], fnEndpoint.EndpointID)
		}
	}

	// fmt.Println("For updating the PEM flags to false following are identified :", mapEndptIdsToInstallFalse)
	// fmt.Println("For updating the PEM flags to true following are identified :", mapEndptIdsToInstallTrue)

	if doDelete {
		processDelete(mapPartnerIDListEndpoints)
		updatePEMInstalledFlag(mapEndptIdsToInstallFalse, mapEndptIdsToInstallTrue)
	}
}

func updatePEMInstalledFlag(mapEndptIdsToInstallFalse, mapEndptIdsToInstallTrue map[string][]string) {
	// fmt.Println("Came to update the PEM")

	for _, endpoints := range mapEndptIdsToInstallFalse {
		updatePEM(false, endpoints)
	}
	for _, endpoints := range mapEndptIdsToInstallTrue {
		updatePEM(true, endpoints)
	}
}

func updatePEM(updateFlgValue bool, endpointIDs []string) {
	//fmt.Printf("\nCame to Update the endpoints : %v, and flag %t", endpointIDs, updateFlgValue)
	sEndpointIDs := strings.Join(endpointIDs, "','")
	sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
	query := fmt.Sprintf(cUpdatePEMInstallFlgQuery, updateFlgValue, sEndpointIDs)
	err := AgentSession.Query(query).Exec()
	if err != nil {
		fmt.Println("Error occured while Running the Following Update Query :", query)
	}
}

func processDelete(mapPartnerIDListEndpoints map[string][]string) {
	fmt.Println("Deleting the Stale entries")
	for partnerID, endpointIDs := range mapPartnerIDListEndpoints {
		lengthOfStale := len(endpointIDs)
		jobs := make(chan string, lengthOfStale)
		results := make(chan asset.AssetCollection, lengthOfStale)
		for w := 1; w <= lengthOfStale; w++ {
			go getEndpointByEndpointID(w, jobs, results, partnerID)
		}
		//fmt.Println("# of workers : ", lengthOfStale)
		for a := 0; a < lengthOfStale; a++ {
			jobs <- endpointIDs[a]
		}
		//fmt.Println("# of jobs ", lengthOfStale)
		close(jobs)
		//fmt.Println("Waiting for results ")
		//ptrRes := make([]PartnerMetrics, len(uniquePartnerIDs))
		var a int

		for a = 0; a < lengthOfStale; a++ {
			assetCollection := <-results
			if assetCollection.EndpointID != "" {
				err := postAssetInactive(assetCollection)
				if err != nil {
					fmt.Println("Error Occured while inserting the Asset in inactive, following is the asset : ", assetCollection)
				}
			}
		}
		close(results)
		err := deleteEndpoints(endpointIDs, partnerID)
		if err != nil {
			fmt.Println("Error Occured while deleting the stale entries from Asset Collection for Partner Id : ", partnerID)
		}
		err = deleteEndpointsFromVersion(endpointIDs, partnerID)
		if err != nil {
			fmt.Println("Error Occured while deleting the stale entries from Agent Version Details for Partner Id : ", partnerID)
		}

	}
}

func postAssetInactive(asset asset.AssetCollection) error {
	fmt.Printf("\n Asset data inserted in inactive table : %v", asset)
	return AssetSession.Query(queryToInsertAssetInActive,
		asset.PartnerID,
		asset.ClientID,
		asset.SiteID,
		asset.EndpointID,
		asset.RegID,
		time.Now(),
		asset.CreateTimeUTC,
		asset.CreatedBy,
		asset.Name,
		asset.Type,
		asset.RemoteAddress,
		asset.ResourceType,
		asset.EndpointType,
		asset.BaseBoard,
		asset.Bios,
		asset.Drives,
		asset.Memory,
		asset.Networks,
		asset.Os,
		asset.Processors,
		asset.RaidController,
		asset.System,
		asset.InstalledSoftwares,
		asset.Keyboards,
		asset.Mouse,
		asset.Monitors,
		asset.PhysicalDrives,
		asset.Users,
		asset.Services,
		asset.Shares).Exec()
}

func preProcessMailBox(mapPartIDEndpointIDsForMailBox map[string][]string) {
	partnerIDs := make([]string, 0)
	for partnerID, _ := range mapPartIDEndpointIDsForMailBox {
		partnerIDs = append(partnerIDs, partnerID)
	}
	totalNoPartnersUnique := len(partnerIDs)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan []string, totalNoPartnersUnique)
	for w := 1; w <= noOfWorkers; w++ {
		resp := make([]string, 0)
		go sendMessages(mapPartIDEndpointIDsForMailBox, jobs, results, configObject, resp)
	}

	for a := 0; a < totalNoPartnersUnique; a++ {
		jobs <- partnerIDs[a]
	}
	close(jobs)

	var a int
	//resp := make([]string, 0)
	for a = 0; a < totalNoPartnersUnique; a++ {
		<-results

		//fmt.Printf("Mail Box Messages Sent : %v", resp)
	}
	close(results)

}

func sendMessages(mapPartIDEndpointIDsForMailBox map[string][]string, jobs <-chan string, results chan<- []string, message Config, resp []string) {
	for partnerID := range jobs {
		endpoints := mapPartIDEndpointIDsForMailBox[partnerID]
		for _, m := range message.Messages {
			msg := Mailbox{
				Endpoints: endpoints,
				Message:   m,
			}
			encoded, err := json.Marshal(msg)
			if err != nil {
				results <- resp
				continue
			}
			payload := strings.NewReader(string(encoded))
			req, err := http.NewRequest("POST", message.ServerAddress, payload)
			if err != nil {
				fmt.Println(time.Now(), " Error while creating request ", err, " for endpoints ", endpoints)
				results <- resp
				continue
			}
			req.Header.Add("content-type", "application/json")
			req.Header.Add("cache-control", "no-cache")
			fmt.Println("Message : ", msg)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(time.Now(), " Error while geting response ", err, " for endpoints ", endpoints)
				results <- resp
				continue
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(time.Now(), " Error while Reading message ", err, " for endpoints ", endpoints)
				results <- resp
				continue
			}
			resp = append(resp, string(body))
			results <- resp
		}
	}

	// results <- resp
}

func returnMapToUpdateFlagsIfMatches(partnerID string, endpoints []fnEndpoint, endpointIDs []string) map[string][]string {
	mapToReturn := make(map[string][]string)
	for _, fnEndpoint := range endpoints {
		endpointID := fnEndpoint.EndpointID
		if contains(endpointIDs, endpointID) {
			mapToReturn[partnerID] = append(mapToReturn[partnerID], endpointID)
		}
	}
	return mapToReturn
}
func returnMapToUpdateFlagsIfNotMatches(partnerID string, endpoints []fnEndpoint, endpointIDs []string) map[string][]string {
	mapToReturn := make(map[string][]string)
	for _, fnEndpoint := range endpoints {
		endpointID := fnEndpoint.EndpointID
		if !contains(endpointIDs, endpointID) {
			mapToReturn[partnerID] = append(mapToReturn[partnerID], endpointID)
		}

	}
	return mapToReturn
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getTheHeartbeatAndInstalledMap(uniquePartnerIDs []string) (map[string][]fnEndpoint, map[string][]fnEndpoint, map[string][]fnEndpoint) {
	totalNoPartnersUnique := len(uniquePartnerIDs)

	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan []fnEndpoint, totalNoPartnersUnique)

	for w := 1; w <= noOfWorkers; w++ {
		go getHeartBeatForPartner(jobs, results)
	}

	for a := 0; a < totalNoPartnersUnique; a++ {
		jobs <- uniquePartnerIDs[a]
	}

	close(jobs)

	//ptrRes := make([]PartnerMetrics, len(uniquePartnerIDs))
	var a int
	endpoints := make([]fnEndpoint, 0)
	for a = 0; a < len(uniquePartnerIDs); a++ {
		endpoints = <-results
		fmt.Println()
		//fmt.Printf("Total no. of endpoint recieved are : %d for partnerId : %s \n", len(endpoints), uniquePartnerIDs[a])
	}
	close(results)
	mapPartIDHbTrueFalseAndInstTrue := make(map[string][]fnEndpoint)
	mapPartIDHbTrueAndInstFalse := make(map[string][]fnEndpoint)
	mapPartIDHbFalseAndInstFalse := make(map[string][]fnEndpoint)
	//mapPartIDHbFalseAndInstTrue := make(map[string][]fnEndpoint)
	for _, endpoint := range endpoints {
		HeartBeartFlag := endpoint.HeartBeartFlag
		installedFlag := endpoint.installedFlag
		partnerID := endpoint.PartnerID
		if (HeartBeartFlag == true && installedFlag == true) || (HeartBeartFlag == false && installedFlag == true) {
			mapPartIDHbTrueFalseAndInstTrue[partnerID] = append(mapPartIDHbTrueFalseAndInstTrue[partnerID], endpoint)
		} else if HeartBeartFlag == true && installedFlag == false {
			mapPartIDHbTrueAndInstFalse[partnerID] = append(mapPartIDHbTrueAndInstFalse[partnerID], endpoint)
		} else if HeartBeartFlag == false && installedFlag == false {
			mapPartIDHbFalseAndInstFalse[partnerID] = append(mapPartIDHbFalseAndInstFalse[partnerID], endpoint)
		}

	}
	// fmt.Println("mapPartIDHbTrueFalseAndInstTrue : ", mapPartIDHbTrueFalseAndInstTrue)
	// fmt.Println("mapPartIDHbTrueAndInstFalse : ", mapPartIDHbTrueAndInstFalse)
	// fmt.Println("mapPartIDHbFalseAndInstFalse : ", mapPartIDHbFalseAndInstFalse)
	// fmt.Println("Done processing results. Total recieved : ", a)
	return mapPartIDHbTrueFalseAndInstTrue, mapPartIDHbTrueAndInstFalse, mapPartIDHbFalseAndInstFalse
}

func processForMailBoxAndLegacy(uniquePartnerIDs []string) map[string][]string {
	//fmt.Println("Started the Processing for Mail box and legacy")
	mapPartnerIDListEndpoints := make(map[string][]string)
	totalNoPartnersUnique := len(uniquePartnerIDs)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan resultValue, totalNoPartnersUnique)

	for w := 1; w <= noOfWorkers; w++ {
		go processEndpoints(jobs, results)
	}
	//fmt.Println("# of workers : ", noOfWorkers)
	for a := 0; a < len(uniquePartnerIDs); a++ {
		jobs <- uniquePartnerIDs[a]
	}
	// /fmt.Println("# of jobs ", totalNoPartnersUnique)
	close(jobs)
	//fmt.Println("Waiting for results ")
	var a int

	for a = 0; a < len(uniquePartnerIDs); a++ {
		// endpoints := make([]string, 0)
		//endpointsValue := resultValue{}
		endpointsValue := <-results
		endpoints := endpointsValue.endpoints
		partnerID := endpointsValue.partnerID
		//fmt.Printf("partnerId : %s & endpoints : %v", partnerID, endpoints)
		if len(endpoints) > 0 {
			//fmt.Printf("For partnerId : %s following endpoints are uninstalled in 1.0 : %v", partnerID, endpoints)
			mapPartnerIDListEndpoints[partnerID] = endpoints
			//fmt.Printf("\nTotal no. of endpoint recieved are : %d for partnerId : %s and endpoints are : %v\n", len(endpoints), uniquePartnerIDs[a], endpoints)
		} else {
			fmt.Println("No Entries detected to be removed for Partner Id : ", partnerID)
		}
	}
	close(results)
	return mapPartnerIDListEndpoints
}

func processEndpoints(jobs <-chan string, results chan<- resultValue) {

	// mapRegIDEndpointID := make(map[string]string)
	for partnerID := range jobs {
		endpoints := make([]string, 0)
		resultValue := resultValue{
			partnerID: partnerID,
		}
		mapRegIDEndpointID, err := getAllNewAssetsByPartnerID(partnerID)
		if err != nil {
			fmt.Printf("** For Partner ID %s, Error Occurred while fetching the juno 2.0 endpoints Error is %v \n \n", partnerID, err)
			results <- resultValue
			continue
		}
		if len(mapRegIDEndpointID) < 1 {
			fmt.Printf("\n\n ** For Partner ID %s, Endpoints not found", partnerID)
			results <- resultValue
			continue
		}
		regIDList := mapPartnerIDRegIDRemoved[partnerID]
		if len(regIDList) < 1 {
			results <- resultValue
			continue
		}
		for _, regID := range regIDList {
			endpointID := mapRegIDEndpointID[regID]
			if endpointID != "" {
				endpoints = append(endpoints, endpointID)
			}
		}
		resultValue.endpoints = endpoints
		results <- resultValue
	}
}

func getAllNewAssetsByPartnerID(partnerID string) (map[string]string, error) {
	mapRegIDEndpointID := make(map[string]string)
	partnerID = strings.TrimSpace(partnerID)
	siteIDs := mapPartnerIDSiteIDs[partnerID]
	var sSiteID string
	for _, siteID := range siteIDs {
		if sSiteID == "" {
			sSiteID = siteID
		} else {
			sSiteID = sSiteID + "," + siteID
		}
	}

	requestURL := fmt.Sprintf(newWebAPIURL, partnerID, sSiteID)
	//fmt.Println("URL : ", requestURL)
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
		endpoinID := asset.EndpointID
		if regID != "" {
			mapRegIDEndpointID[regID] = endpoinID
		}

	}

	return mapRegIDEndpointID, nil
}
func getHeartBeatForPartner(jobs <-chan string, results chan<- []fnEndpoint) {
	for partnerID := range jobs {
		fnEndpoints := getHeartBeatFromCassForPartnerID(partnerID)
		fnEndpointsWithInstalledFlag := populateInstalledFlagForEndpoints(fnEndpoints)
		results <- fnEndpointsWithInstalledFlag
	}
}

func populateInstalledFlagForEndpoints(endpoints []fnEndpoint) []fnEndpoint {
	var endpoint_id, agent_id, client_id, legacy_reg_id, partner_id, site_id string
	var installed bool
	fnEndpoints := make([]fnEndpoint, 0)
	var endpoint fnEndpoint
	mapEndpointHeartBeatFlg := make(map[string]bool)

	//var endpointIDs []string
	endpointIDs := make([]string, 0)
	for _, fnEndpoint := range endpoints {
		endpointIDs = append(endpointIDs, fnEndpoint.EndpointID)
		mapEndpointHeartBeatFlg[fnEndpoint.EndpointID] = fnEndpoint.HeartBeartFlag
	}
	for _, fnEndpoint := range endpoints {
		endpointIDs = append(endpointIDs, fnEndpoint.EndpointID)
	}
	sEndpointIDs := strings.Join(endpointIDs, "','")
	sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)

	query := fmt.Sprintf(queryToPopulateInstalledFlag, sEndpointIDs)
	iter := AgentSession.Query(query).Iter()
	for iter.Scan(&endpoint_id, &agent_id, &client_id, &installed, &legacy_reg_id, &partner_id, &site_id) {
		endpoint = fnEndpoint{
			PartnerID:     partner_id,
			EndpointID:    endpoint_id,
			AgentID:       agent_id,
			ClientID:      client_id,
			installedFlag: installed,
			RegID:         legacy_reg_id,
			siteID:        site_id,
		}
		endpoint.HeartBeartFlag = mapEndpointHeartBeatFlg[endpoint.EndpointID]
		fnEndpoints = append(fnEndpoints, endpoint)
	}
	return fnEndpoints
}

func getHeartBeatFromCassForPartnerID(partnerID string) []fnEndpoint {
	var fnEndpoints []fnEndpoint
	var endpoint fnEndpoint

	var endpoint_id, partner_id string
	var dcDatetimeutc time.Time
	query := fmt.Sprintf(queryToGetHeartBeatByPartnerID, partnerID)
	iter := AgentSession.Query(query).Iter()
	for iter.Scan(&partner_id, &endpoint_id, &dcDatetimeutc) {
		HeartBeatFlag := compareWithCurrDate(dcDatetimeutc)
		endpoint = fnEndpoint{
			PartnerID:      partner_id,
			EndpointID:     endpoint_id,
			HeartBeartFlag: HeartBeatFlag,
		}

		fnEndpoints = append(fnEndpoints, endpoint)
	}

	return fnEndpoints
}

func compareWithCurrDate(dcDatetimeutc time.Time) bool {
	year, month, day := dcDatetimeutc.Date()
	currYear, currMonth, currDay := time.Now().Date()
	if year == currYear && month == currMonth && day >= currDay-1 {
		return true
	}
	return false
}

//GetPartnerIDs  is the function to get the partnerId
func GetPartnerIDs() ([]string, error) {
	var partner_id string
	var partnerIDs []string
	iter := AssetSession.Query(queryToGetPartners).Iter()
	for iter.Scan(&partner_id) {
		partnerIDs = append(partnerIDs, partner_id)
	}
	if err := iter.Close(); err != nil {
		fmt.Println("Error occured while getting the partner : ", err)
		return partnerIDs, err
	}
	return partnerIDs, nil
}

func getCassandraSession(assetClusterIP, agentClusterIP string) error {
	var err error
	assetCluster := gocql.NewCluster(assetClusterIP)
	assetCluster.Consistency = gocql.Quorum
	assetCluster.Keyspace = "platform_asset_db"
	AssetSession, err = assetCluster.CreateSession()
	if err != nil {
		fmt.Printf("Error occured in cassandra setup for Asset : %v", err)
		return err
	}

	agentCluster := gocql.NewCluster(agentClusterIP)
	agentCluster.Consistency = gocql.Quorum
	agentCluster.Keyspace = "platform_agent_db"
	AgentSession, err = agentCluster.CreateSession()
	if err != nil {
		fmt.Printf("Error occured in cassandra setup for Agent : %v", err)
		return err
	}
	return nil
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

func getEndpointByEndpointID(id int, jobs <-chan string, results chan<- asset.AssetCollection, partnerID string) {
	for endpointID := range jobs {
		//fmt.Println(getAssetURL)
		assetColl := asset.AssetCollection{}
		partnerID = strings.TrimSpace(partnerID)
		endpointID = strings.TrimSpace(endpointID)
		requestURL := fmt.Sprintf(getAssetURL, partnerID, endpointID)
		res, err := http.Get(requestURL)
		if err != nil {
			//fmt.Println("RMM 2.0 : Error 1.0 occured while getting the response for partner Id :", partnerID)
			results <- assetColl
			continue
		}
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			//fmt.Println("RMM 2.0 : Error 2.0 occured while Reading the response for partner Id :", partnerID)
			results <- assetColl
			continue
		}
		bodyString := string(bodyBytes)

		err = json.Unmarshal([]byte(bodyString), &assetColl)
		if err != nil {
			//fmt.Println("RMM 2.0 : Error 3.0 occured while  unmarshalling the response for partner Id :", partnerID)
		}
		results <- assetColl
	}

}
func deleteEndpoints(endpointIDs []string, partnerID string) error {
	//fmt.Println("Deleting the Assets for Endpoint ids : ", endpointIDs)
	sEndpointIDs := strings.Join(endpointIDs, "','")
	sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
	query := fmt.Sprintf(queryToDeleteEndpoints, partnerID, sEndpointIDs)
	err := AssetSession.Query(query).Exec()
	return err
}

func deleteEndpointsFromVersion(endpointIDs []string, partnerID string) error {
	//fmt.Println("Deleting the Assets for Endpoint ids : ", endpointIDs)
	sEndpointIDs := strings.Join(endpointIDs, "','")
	sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
	query := fmt.Sprintf(queryToDeleteAgentVersionDetails, partnerID, sEndpointIDs)
	err := AssetSession.Query(query).Exec()
	return err
}

//ReadPartnersFromExcel ............
func ReadPartnersFromExcel(excelPath string, sheetName string) (map[string][]string, error) {
	partnerIDEndpointIDs := make(map[string][]string)
	xlsx, err := excelize.OpenFile(excelPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("Sheet1")

	for count, row := range rows {
		if count == 0 {
			continue
		}
		partnerIDEndpointIDs[row[0]] = append(partnerIDEndpointIDs[row[0]], row[1])
	}
	return partnerIDEndpointIDs, nil
}

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

func readEndpoints(sheetPath string) (map[string][]string, error) {
	mapPartnerIDRegIDRemoved = make(map[string][]string)
	xlFile, err := excelize.OpenFile(sheetPath)
	if err != nil {
		logger.Println("Error opening excel sheet")
		return mapPartnerIDRegIDRemoved, err
	}
	rows := xlFile.GetRows("Sheet1")
	for _, row := range rows {
		if retVal, ok := mapPartnerIDRegIDRemoved[row[0]]; ok {
			retVal = append(retVal, row[1])
			mapPartnerIDRegIDRemoved[row[0]] = retVal
		} else {
			mapPartnerIDRegIDRemoved[row[0]] = retVal
		}
	}
	logger.Println(mapPartnerIDRegIDRemoved)
	return mapPartnerIDRegIDRemoved, err
}

const logFlags int = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile
