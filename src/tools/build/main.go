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
	"github.com/ContinuumLLC/platform-common-lib/src/utils"
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
	SiteID       int    `json:"siteId"`
	EndpointID   string `json:"endpointId"`
	RegType      string `json:"regType"`
}

//Session session variable
var Session *gocql.Session
var assetClusterIP string
var legacyWebAPIURL, newWebAPIURL string
var totalNoPartnersFromCSV, totalNoPartners int

const (
	findEndpointSummaryQuery   = "SELECT partner_id, site_id, client_id,endpoint_id,reg_id, name,friendly_name,remote_address, type, resource_type, endpoint_type, system, os,created_by, agent_ts_utc,dc_ts_utc  FROM platform_asset_db.endpoint_summary WHERE partner_id = ? AND site_id = ? AND endpoint_id = ?"
	findEndpointDetailsQuery   = "SELECT partner_id,endpoint_id,baseboard,bios,drives,memory,physical_memory,networks,processors,raidcontroller,installed_softwares, keyboards, mouse, monitors, physical_drives,users,shares,software_licenses FROM platform_asset_db.endpoint_details  WHERE partner_id = ? AND site_id =? AND endpoint_id = ?"
	lookupEndpointQuery        = "SELECT client_id,site_id,reg_id FROM platform_asset_db.endpoint_mapping WHERE partner_id=? and endpoint_id=?"
	insertEndpointSummaryQuery = "INSERT INTO platform_asset_db.endpoint_summary(partner_id, site_id, client_id,endpoint_id,reg_id, name,friendly_name,remote_address, type, resource_type, endpoint_type, system, os,created_by, agent_ts_utc,dc_ts_utc) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	insertEndpointMappingQuery = "INSERT INTO platform_asset_db.endpoint_mapping(partner_id,endpoint_id,site_id, client_id,reg_id,dc_ts_utc,agent_ts_utc) VALUES(?,?,?,?,?,?,?)"
	//insertEndpointServiceDetailsQuery = "INSERT INTO platform_asset_db.endpoint_service(partner_id,endpoint_id,service_name,details,agent_ts_utc,dc_ts_utc) VALUES(?,?,?,?,?,?)"
	insertEndpointDetailsQuery     = "INSERT INTO platform_asset_db.endpoint_details (partner_id,site_id,endpoint_id,baseboard,bios,drives,physical_memory,networks,processors,raidcontroller,installed_softwares, keyboards, mouse, monitors, physical_drives,users,shares,software_licenses,agent_ts_utc,dc_ts_utc) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	updateEndpointSiteMappingQuery = "UPDATE platform_asset_db.endpoint_mapping set site_id = ?,client_id=? WHERE partner_id=? and endpoint_id=?"
	deleteEndpointDetailsQuery     = "DELETE FROM platform_asset_db.endpoint_details WHERE partner_id = ? AND site_id = ? AND endpoint_id = ?"
	deleteEndpointSummaryQuery     = "DELETE FROM platform_asset_db.endpoint_summary WHERE partner_id = ? AND site_id = ? AND endpoint_id = ?"
	deleteEndpointMappingQuery     = "DELETE FROM platform_asset_db.endpoint_mapping WHERE partner_id = ? AND endpoint_id = ?"
	resourceTypeDesktop            = "Desktop"
	resourceTypeServer             = "Server"
	osTypeWindows                  = "Windows"
	//updatePEMQuery                 = `update platform_agent_db.partnerendpointmap set installed = %t where endpoint_id in (%s)`
)

var configObject Config

//Config is the struct rep. of config
type Config struct {
	ResourceType          string
	LegacyWebAPIURL       string
	IsDelete              bool
	NumOfWorkers          int
	CassandraHostForAsset string
	CSVPathForPartners    string
}

type endpointSummary struct {
	PartnerID        string
	SiteID           string
	ClientID         string
	EndpointID       string
	RegID            string
	FriendlyName     string
	CreatedTimeUTC   time.Time
	EndpointType     string
	ResourceType     string
	InsertedDTimeUTC time.Time
}

func create(line []string) endpointSummary {
	return endpointSummary{
		PartnerID:        line[0],
		SiteID:           line[1],
		ClientID:         line[2],
		EndpointID:       line[3],
		RegID:            line[4],
		FriendlyName:     line[5],
		CreatedTimeUTC:   time.Now().UTC(),
		EndpointType:     line[7],
		ResourceType:     line[8],
		InsertedDTimeUTC: time.Now().UTC(),
	}
}

func readAssetTables(filePath string) (map[string]map[string]endpointSummary, []string, error) {
	partnerAssetCollection := make(map[string]map[string]endpointSummary)
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

		partnerID := strings.TrimSpace(astCol.PartnerID)
		if partnerID == "" {
			continue
		}
		partnerIDs = append(partnerIDs, partnerID)
		regIDMap := partnerAssetCollection[partnerID]
		if regIDMap == nil {
			regIDMap = make(map[string]endpointSummary)
		}
		rID := strings.TrimSpace(astCol.RegID)
		asssetMigCollection := regIDMap[rID]
		if asssetMigCollection.EndpointID == "" {
			regIDMap[rID] = astCol
			partnerAssetCollection[partnerID] = regIDMap
		}
	}
	fmt.Println("------------------partnerAssetCollection----------------", partnerAssetCollection)
	return partnerAssetCollection, partnerIDs, err
}

//LoadConfiguration is a method to load configuration File
func LoadConfiguration(filePath string) (Config, error) {
	cfg := &Config{}
	file, err := os.Open(filePath)
	if err != nil {
		return *cfg, err
	}
	defer file.Close()

	deser := json.NewDecoder(file)
	err = deser.Decode(cfg)
	if err != nil {
		return *cfg, err
	}
	return *cfg, nil
}
func main() {
	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: syncSites config.json")
		return
	}
	configObject, err := LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}
	startTime := time.Now()
	fmt.Println("Migration Tool Started... Time started : ", startTime)
	excelPath := configObject.CSVPathForPartners
	legacyWebAPIURL = configObject.LegacyWebAPIURL
	assetClusterIP = configObject.CassandraHostForAsset
	noOfWorkers := configObject.NumOfWorkers
	resourceType := configObject.ResourceType
	IsDelete := configObject.IsDelete
	if legacyWebAPIURL != "" && resourceType != "" {
		legacyWebAPIURL = legacyWebAPIURL + "/v1/partner/%s/endpoints?resourceType=" + resourceType
	}
	cassError := getCassandraSession(assetClusterIP)
	if cassError != nil {
		fmt.Println("Error Occured while setting up the Tool with Cassandra, Error : ", cassError)
		return
	}
	partnerAssetDBData, partnerIDs, asstErr := readAssetTables(excelPath)
	if asstErr != nil || len(partnerIDs) < 1 {
		fmt.Println("Error Occured while readAssetTables from Excel, partnerIDs :  , Error : ", partnerIDs, asstErr)
		return
	}
	//Getting Unique PartnerIds
	uniquePartnerIDs, err := removeDuplicatePartnerIds(partnerIDs)
	if err != nil {
		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
		return
	}
	totalNoPartnersUnique := len(uniquePartnerIDs)
	fmt.Println("partnerIDs : ", totalNoPartnersUnique)
	jobs := make(chan string, totalNoPartnersUnique)
	results := make(chan PartnerMetrics, totalNoPartnersUnique)

	for w := 1; w <= noOfWorkers; w++ {
		go processParters(w, jobs, results, partnerAssetDBData, IsDelete)
	}
	fmt.Println("# of workers : ", noOfWorkers)

	for a := 0; a < totalNoPartnersUnique; a++ {
		jobs <- uniquePartnerIDs[a]
	}
	fmt.Println("# of jobs ", totalNoPartnersUnique)
	close(jobs)

	fmt.Println("Waiting for results ")
	ptrRes := make([]PartnerMetrics, len(partnerIDs))
	var a int
	for a = 0; a < totalNoPartnersUnique; a++ {
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

func getAllLegacyAssetsByPartnerID(partnerID string) (assetCollection, error) {
	assetColl := assetCollection{}
	partnerID = strings.TrimSpace(partnerID)
	requestURL := fmt.Sprintf(legacyWebAPIURL, partnerID)
	fmt.Println("requestURL : ", requestURL)
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
	EID    string
	FName  string
	SiteID string
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

func processParters(id int, jobs <-chan string, results chan<- PartnerMetrics, assetDataMap map[string]map[string]endpointSummary, IsDelete bool) {
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
			ls.RMM1Err = err.Error()
			results <- ls
			rmmoneerr = rmmoneerr + 1
			continue
		}
		rmm2ProcessingTime := time.Now()
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
			legacySiteID := strconv.Itoa(EndPointList.SiteID)
			regType := EndPointList.RegType
			resourceType := getResourceTypeByRegType(regType)
			if resourceType == "" {
				fmt.Printf("\n No resource type found for regID : %s, endpointID : %s, partnerId : %s, siteId : %s", regID, endpoint.EndpointID, endpoint.PartnerID, endpoint.SiteID)
			}
			if endpoint.EndpointID != "" {
				if endpoint.SiteID == legacySiteID {
					inSync = inSync + 1
					ls.FNameInSyncEndpoints = append(ls.FNameInSyncEndpoints, endpoint.EndpointID)
					continue
				}
				dbUpdateProcessingTime := time.Now()
				//To update the sites and delete the existing.
				if IsDelete {
					errs := UpdateAndDeleteSites(endpoint, legacySiteID, resourceType)
					fmt.Printf("\n error occurred for update \n: %+v", errs)
					//errs := UpdateFriendLyNameByEndpointID(endpoint.EID, partnerID, friendlyName)
					dbupdatecntr = dbupdatecntr + 1
					dbupdatett = dbupdatett + time.Since(dbUpdateProcessingTime).Seconds()
					if errs != nil {
						//fmt.Println("Error while getting the endpoints from Juno for the partner Id : ", partnerID)
						dbErr = dbErr + 1
						ls.UpdateDBErrorEndpoints = append(ls.UpdateDBErrorEndpoints, endpoint.EndpointID)

					} else {
						//fmt.Printf("Friendly Name updated. PartnerID : %s , Reg Id : %s , Endpoint ID : %s  and Friendly Name : %s\n", partnerID, regID, endpointID, friendlyName)
						counter = counter + 1
						ls.UpdatedEndpoints = append(ls.UpdatedEndpoints, endpoint.EndpointID)
					}
				} else {
					ls.UpdatedEndpoints = append(ls.UpdatedEndpoints, endpoint.EndpointID)
				}

			} else {
				//fmt.Printf("Endpoint id  not found for PartnerID : %s , reg Id : %s and friendly Name : %s\n", partnerID, regID, friendlyName)
				endpointNotFoundCounter = endpointNotFoundCounter + 1
				ls.RMM2NoEndpoints = append(ls.RMM2NoEndpoints, endpoint.EndpointID)

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

func getResourceTypeByRegType(regType string) string {
	switch regType {
	case "DPMA":
		return resourceTypeDesktop
	case "MSMA":
		return resourceTypeServer
	default:
		return ""
	}
}

//UpdateAndDeleteSites updates and delete the endpoints
func UpdateAndDeleteSites(endpoint endpointSummary, legacySiteID string, resourceType string) error {
	fmt.Println("endpoint : ", endpoint)
	asset, err := GetEndpointData(endpoint.PartnerID, endpoint.SiteID, endpoint.EndpointID)
	if asset == nil {
		return fmt.Errorf("No Asset Data found for Partner Id := %s, SiteId := %s and endpoint_id=%s", endpoint.PartnerID, endpoint.SiteID, endpoint.EndpointID)
	}
	if err != nil {
		return fmt.Errorf("Error Occurred while gettting the endpoint data for partner ID  :%s, endpointID : %s, SiteId : %s And err : %+v", endpoint.PartnerID, endpoint.EndpointID, endpoint.SiteID, err)
	}
	errForDel := Session.Query(deleteEndpointMappingQuery, endpoint.PartnerID, endpoint.EndpointID).Exec()
	if errForDel != nil {
		return fmt.Errorf("Error Occurred while Deleting the Mapping for PartnerID : %s, EndpointID : %s & Error : %+v", endpoint.PartnerID, endpoint.EndpointID, errForDel)
	}
	errForDel = Session.Query(deleteEndpointSummaryQuery, endpoint.PartnerID, endpoint.SiteID, endpoint.EndpointID).Exec()
	if errForDel != nil {
		return fmt.Errorf("Error Occurred while Deleting the Summary for PartnerID : %s, EndpointID : %s, SiteID : %s & Error : %+v", endpoint.PartnerID, endpoint.EndpointID, endpoint.SiteID, errForDel)
	}
	errForDel = Session.Query(deleteEndpointDetailsQuery, endpoint.PartnerID, endpoint.SiteID, endpoint.EndpointID).Exec()
	if errForDel != nil {
		return fmt.Errorf("Error Occurred while Deleting the Details for PartnerID : %s, EndpointID : %s, SiteID : %s & Error : %+v", endpoint.PartnerID, endpoint.EndpointID, endpoint.SiteID, errForDel)
	}
	asset.SiteID = legacySiteID
	asset.EndpointType = resourceType
	return SaveAssets(asset)

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

func getNilSafeValue(data map[string]interface{}, key string) string {
	val := data[key]
	if val == nil {
		return ""
	}
	return val.(string)
}

func Select(query string, value ...interface{}) ([]map[string]interface{}, error) {
	q := Session.Query(query, value...) //.Consistency(gocql.One)
	iter := Session.Query(query, value...).Iter()
	data, err := iter.SliceMap()
	defer q.Release()
	if err != nil {
		return nil, fmt.Errorf("ErrDbUnableToFetchRecord %+v", err)
	}
	err = iter.Close()
	if err != nil {
		return data, fmt.Errorf("ErrDbUnableToFetchRecord %+v", err)
	}
	return data, nil
}

//GetEndpointData will give entire data available for an endpoint attach to site of a partner in the system and as per filter level types provided
func GetEndpointData(partnerID, siteID, endpointID string) (*asset.AssetCollection, error) {
	var resAssetCollection *asset.AssetCollection
	var err error
	var dbRecords []map[string]interface{}
	dbRecords, err = Select(findEndpointSummaryQuery, partnerID, siteID, endpointID)
	if err != nil || len(dbRecords) <= 0 {
		return resAssetCollection, err
	}
	resAssetCollection = &asset.AssetCollection{}
	mapSummaryData(resAssetCollection, dbRecords[0])
	dbRecords, err = Select(findEndpointDetailsQuery, partnerID, siteID, endpointID)
	if err == nil && len(dbRecords) > 0 {
		mapDetailsData(resAssetCollection, dbRecords[0])
	}
	return resAssetCollection, err
}

func mapSummaryData(assetData *asset.AssetCollection, data map[string]interface{}) {
	assetData.CreateTimeUTC = data["agent_ts_utc"].(time.Time)
	assetData.CreatedBy = data["created_by"].(string)
	assetData.Name = utils.ToString(data["name"])
	assetData.Type = utils.ToString(data["type"])
	assetData.PartnerID = data["partner_id"].(string)
	assetData.ClientID = data["client_id"].(string)
	assetData.SiteID = data["site_id"].(string)
	assetData.EndpointID = data["endpoint_id"].(string)
	assetData.RegID = data["reg_id"].(string)
	assetData.FriendlyName = utils.ToString(data["friendly_name"])
	assetData.ResourceType = utils.ToString(data["resource_type"])
	assetData.EndpointType = utils.ToString(data["endpoint_type"])
	assetData.RemoteAddress = utils.ToString(data["remote_address"])
	assetData.Os = getOS(data)
	assetData.System = getSystem(data)
	//TODO This code will be removed once resourceType fix is in production for a month or so
	//This code changes allows us to categorize resourceType into Desktop or server even if the actual endpoint has been offline even before this fix went in
	if assetData.ResourceType == "" {
		if strings.Contains(strings.ToLower(assetData.Os.Product), strings.ToLower(resourceTypeServer)) &&
			strings.Contains(strings.ToLower(assetData.Os.Type), strings.ToLower(osTypeWindows)) {
			assetData.ResourceType = resourceTypeServer
		} else {
			assetData.ResourceType = resourceTypeDesktop
		}
	}
}

func mapDetailsData(assetData *asset.AssetCollection, data map[string]interface{}) {
	assetData.BaseBoard = getBaseBoard(data)
	assetData.Bios = getBios(data)
	assetData.RaidController = getRaidCtrlr(data)
	assetData.Drives = getAssetDrive(data)
	assetData.Memory = getMemory(data)
	assetData.Networks = getAssetNetwork(data)
	assetData.Processors = getAssetProcessor(data)
	assetData.InstalledSoftwares = getInstalledSoftwares(data)
	assetData.Keyboards = getKeyboards(data)
	assetData.Mouse = getMouse(data)
	assetData.Monitors = getMonitors(data)

	if data["physical_drives"] != nil {
		physDrives := data["physical_drives"].([]map[string]interface{})
		pDrives := make([]asset.PhysicalDrive, len(physDrives))
		i := 0
		for _, drive := range physDrives {
			pDrives[i] = asset.PhysicalDrive{
				Type:          utils.ToString(drive["type"]),
				PartitionData: partitionData(drive),
			}
			i++
		}
		assetData.PhysicalDrives = pDrives
	}
	assetData.Users = getUser(data)
	assetData.Shares = getShares(data)
	assetData.SoftwareLicenses = getSoftwareLicenses(data)

}

func getSoftwareLicenses(data map[string]interface{}) []asset.SoftwareLicenses {
	if data["software_licenses"] == nil {
		return nil
	}
	u := data["software_licenses"].([]map[string]interface{})
	l := len(u)

	if l <= 0 {
		return nil
	}
	softwarelicenses := make([]asset.SoftwareLicenses, l)

	for i := 0; i < l; i++ {
		softwarelicenses[i] = asset.SoftwareLicenses{
			ProductName: utils.ToString(u[i]["product_name"]),
			ProductID:   utils.ToString(u[i]["product_id"]),
			ProductKey:  utils.ToString(u[i]["product_key"]),
		}
	}
	return softwarelicenses
}

func getShares(data map[string]interface{}) []asset.Share {
	if data["shares"] == nil {
		return nil
	}
	u := data["shares"].([]map[string]interface{})
	l := len(u)

	if l <= 0 {
		return nil
	}
	shares := make([]asset.Share, l)

	for i := 0; i < l; i++ {
		shares[i] = asset.Share{
			Name:        utils.ToString(u[i]["name"]),
			Caption:     utils.ToString(u[i]["caption"]),
			Description: utils.ToString(u[i]["description"]),
			Path:        utils.ToString(u[i]["path"]),
			Access:      utils.ToString(u[i]["access"]),
			//UserAccess:  u[i]["user_access"].([]string),
			UserAccess: utils.ToStringArray(u[i]["user_access"]),
			//Type:       u[i]["type"].([]string),
			Type: utils.ToStringArray(u[i]["type"]),
		}
	}
	return shares
}

func getUser(data map[string]interface{}) []asset.User {
	if data["users"] == nil {
		return nil
	}

	u := data["users"].([]map[string]interface{})
	l := len(u)
	users := make([]asset.User, l)
	for i := 0; i < l; i++ {
		users[i] = asset.User{
			PasswordRequired:          utils.ToBool(u[i]["password_required"]),
			UserDisabled:              utils.ToBool(u[i]["user_disabled"]),
			UserLockout:               utils.ToBool(u[i]["user_lockout"]),
			Username:                  utils.ToString(u[i]["username"]),
			UserType:                  utils.ToString(u[i]["user_type"]),
			DomainName:                utils.ToString(u[i]["domain_name"]),
			PasswordChangeable:        utils.ToBool(u[i]["password_changeable"]),
			PasswordExpires:           utils.ToBool(u[i]["password_expires"]),
			UserID:                    utils.ToString(u[i]["userid"]),
			AccountType:               utils.ToString(u[i]["account_type"]),
			PasswordComplexityEnabled: utils.ToBool(u[i]["password_complexity_enabled"]),
			LockoutThreshold:          uint32(utils.ToInt(u[i]["lockout_threshold"])),
			LockoutDurationMins:       uint32(utils.ToInt(u[i]["lockout_duration_mins"])),
			SystemRole:                utils.ToString(u[i]["system_role"]),
			LastLogonTimestamp:        uint32(utils.ToInt(u[i]["lastlogon_timestamp"])),
			RemoteDesktopAllowed:      utils.ToBool(u[i]["remote_desktop_allowed"]),
			RemoteAccessAllowed:       utils.ToBool(u[i]["remote_access_allowed"]),
		}
	}
	return users
}

func getSystem(data map[string]interface{}) asset.AssetSystem {
	if data["system"] == nil {
		return asset.AssetSystem{}
	}
	systemMap := data["system"].(map[string]interface{})
	if len(systemMap) <= 0 {
		return asset.AssetSystem{}
	}
	var sys asset.AssetSystem
	sys.Category = utils.ToString(systemMap["category"])
	sys.Model = utils.ToString(systemMap["model"])
	sys.Product = utils.ToString(systemMap["product"])
	sys.SerialNumber = utils.ToString(systemMap["serial_number"])
	sys.SystemName = utils.ToString(systemMap["system_name"])
	sys.TimeZone = utils.ToString(systemMap["time_zone"])
	sys.TimeZoneDescription = utils.ToString(systemMap["time_zone_description"])
	sys.Domain = utils.ToString(systemMap["domain"])
	sys.DomainRole = utils.ToInt(systemMap["domain_role"])
	return sys
}

func getRaidCtrlr(data map[string]interface{}) asset.AssetRaidController {

	if data["raidcontroller"] == nil {
		return asset.AssetRaidController{}
	}
	raidCtrlMap := data["raidcontroller"].(map[string]interface{})
	if len(raidCtrlMap) <= 0 {
		return asset.AssetRaidController{}
	}
	var rc asset.AssetRaidController
	rc.HardwareRaid = utils.ToString(raidCtrlMap["hardware_raid"])
	rc.SoftwareRaid = utils.ToString(raidCtrlMap["software_raid"])
	rc.Vendor = utils.ToString(raidCtrlMap["vendor"])
	return rc
}

func getOS(data map[string]interface{}) asset.AssetOs {
	if data["os"] == nil {
		return asset.AssetOs{}
	}
	osMap := data["os"].(map[string]interface{})
	if len(osMap) <= 0 {
		return asset.AssetOs{}
	}
	var os asset.AssetOs
	os.Manufacturer = utils.ToString(osMap["manufacturer"])
	os.Product = utils.ToString(osMap["product"])
	os.OsLanguage = utils.ToString(osMap["os_language"])
	os.SerialNumber = utils.ToString(osMap["serial_number"])
	os.Version = utils.ToString(osMap["version"])
	os.InstallDate = utils.ToTime(osMap["installdate"])
	os.Type = utils.ToString(osMap["type"])
	os.Arch = utils.ToString(osMap["arch"])
	os.ServicePack = utils.ToString(osMap["service_pack"])
	os.OsInstalledDrive = utils.ToString(osMap["os_installed_drive"])
	os.BuildNumber = utils.ToString(osMap["build_number"])
	os.ProductID = utils.ToString(osMap["product_id"])
	os.ProductKey = utils.ToString(osMap["product_key"])
	return os

}

func getMemory(data map[string]interface{}) []asset.PhysicalMemory {
	if data["physical_memory"] == nil {
		return nil
	}

	mem := data["physical_memory"].([]map[string]interface{})
	l := len(mem)
	pMem := make([]asset.PhysicalMemory, l)
	for i := 0; i < l; i++ {
		pMem[i] = asset.PhysicalMemory{
			Manufacturer: utils.ToString(mem[i]["manufacturer"]),
			SerialNumber: utils.ToString(mem[i]["serial_number"]),
			SizeBytes:    uint64(utils.ToInt64(mem[i]["size_bytes"])),
		}
	}
	return pMem
}

func getBios(data map[string]interface{}) asset.AssetBios {
	if data["bios"] == nil {
		return asset.AssetBios{}
	}
	biosMap := data["bios"].(map[string]interface{})
	if len(biosMap) <= 0 {
		return asset.AssetBios{}
	}
	var bs asset.AssetBios
	bs.Manufacturer = utils.ToString(biosMap["manufacturer"])
	bs.Product = utils.ToString(biosMap["product"])
	bs.SerialNumber = utils.ToString(biosMap["serial_number"])
	bs.SmbiosVersion = utils.ToString(biosMap["smbios_version"])
	bs.Version = utils.ToString(biosMap["version"])
	return bs
}

func getAssetDrive(data map[string]interface{}) []asset.AssetDrive {
	if data["drives"] == nil {
		return nil
	}
	ad := data["drives"].([]map[string]interface{})
	l := len(ad)
	drives := make([]asset.AssetDrive, l)
	for i := 0; i < l; i++ {
		var noOfPartition int
		if ad[i]["number_of_partitions"] != nil {
			noOfPartition = utils.ToInt(ad[i]["number_of_partitions"])
		}
		drives[i] = asset.AssetDrive{
			Manufacturer:       utils.ToString(ad[i]["manufacturer"]),
			MediaType:          utils.ToString(ad[i]["media_type"]),
			InterfaceType:      utils.ToString(ad[i]["interface_type"]),
			LogicalName:        utils.ToString(ad[i]["logical_name"]),
			SerialNumber:       utils.ToString(ad[i]["serial_number"]),
			SizeBytes:          utils.ToInt64(ad[i]["size_bytes"]),
			Product:            utils.ToString(ad[i]["product"]),
			Partitions:         utils.ToStringArray(ad[i]["partitions"]),
			NumberOfPartitions: noOfPartition,
			PartitionData:      partitionData(ad[i]),
		}
	}
	return drives
}

func partitionData(drivesMap map[string]interface{}) []asset.DrivePartition {
	if drivesMap["partition_data"] == nil {
		return nil
	}
	partitionData := drivesMap["partition_data"].([]map[string]interface{})
	l := len(partitionData)
	drives := make([]asset.DrivePartition, l)
	for i := 0; i < l; i++ {
		drives[i].Description = utils.ToString(partitionData[i]["description"])
		drives[i].FileSystem = utils.ToString(partitionData[i]["file_system"])
		drives[i].Label = utils.ToString(partitionData[i]["label"])
		drives[i].Name = utils.ToString(partitionData[i]["name"])
		drives[i].SizeBytes = utils.ToInt64(partitionData[i]["size_bytes"])
	}
	return drives
}
func getAssetNetwork(data map[string]interface{}) []asset.AssetNetwork {
	if data["networks"] == nil {
		return nil
	}
	net := data["networks"].([]map[string]interface{})
	l := len(net)
	assetNet := make([]asset.AssetNetwork, l)
	for i := 0; i < l; i++ {
		assetNet[i] = asset.AssetNetwork{
			DefaultIPGateway:    utils.ToString(net[i]["default_ip_gateway"]),
			DefaultIPGateways:   utils.ToStringArray(net[i]["default_ip_gateways"]),
			DhcpEnabled:         utils.ToBool(net[i]["dhcp_enabled"]),
			DhcpServer:          utils.ToString(net[i]["dhcp_server"]),
			DhcpLeaseObtained:   utils.ToTime(net[i]["dhcp_lease_obtained"]),
			DhcpLeaseExpires:    utils.ToTime(net[i]["dhcp_lease_expires"]),
			MacAddress:          utils.ToString(net[i]["mac_address"]),
			IPEnabled:           utils.ToBool(net[i]["ip_enabled"]),
			IPv4:                utils.ToString(net[i]["ipv4"]),
			IPv4List:            utils.ToStringArray(net[i]["ipv4_list"]),
			IPv6:                utils.ToString(net[i]["ipv6"]),
			IPv6List:            utils.ToStringArray(net[i]["ipv6_list"]),
			Product:             utils.ToString(net[i]["product"]),
			SubnetMask:          utils.ToString(net[i]["subnet_mask"]),
			SubnetMasks:         utils.ToStringArray(net[i]["subnet_masks"]),
			Vendor:              utils.ToString(net[i]["vendor"]),
			DnsServers:          utils.ToStringArray(net[i]["dns_servers"]),
			LogicalName:         utils.ToString(net[i]["logical_name"]),
			WinsPrimaryServer:   utils.ToString(net[i]["wins_primary_server"]),
			WinsSecondaryServer: utils.ToString(net[i]["wins_secondary_server"]),
		}
	}
	return assetNet
}

func getAssetProcessor(data map[string]interface{}) []asset.AssetProcessor {
	if data["processors"] == nil {
		return nil
	}
	proc := data["processors"].([]map[string]interface{})
	l := len(proc)
	assetProc := make([]asset.AssetProcessor, l)
	for i := 0; i < l; i++ {
		assetProc[i] = asset.AssetProcessor{
			ClockSpeedMhz: utils.ToFloat64(proc[i]["clock_speed_mhz"]),
			Family:        utils.ToInt(proc[i]["family"]),
			Manufacturer:  utils.ToString(proc[i]["manufacturer"]),
			NumberOfCores: utils.ToInt(proc[i]["number_of_cores"]),
			ProcessorType: utils.ToString(proc[i]["processor_type"]),
			Product:       utils.ToString(proc[i]["product"]),
			SerialNumber:  utils.ToString(proc[i]["serial_number"]),
			Level:         utils.ToInt(proc[i]["level"]),
		}
	}
	return assetProc
}

func getInstalledSoftwares(data map[string]interface{}) []asset.InstalledSoftware {
	if data["installed_softwares"] == nil {
		return nil
	}
	sw := data["installed_softwares"].([]map[string]interface{})
	l := len(sw)
	insSw := make([]asset.InstalledSoftware, l)
	for i := 0; i < l; i++ {
		insSw[i] = asset.InstalledSoftware{
			Name:               utils.ToString(sw[i]["name"]),
			Version:            utils.ToString(sw[i]["version"]),
			InstallDate:        utils.ToTime(sw[i]["install_date"]),
			Publisher:          utils.ToString(sw[i]["publisher"]),
			UserName:           utils.ToString(sw[i]["user_name"]),
			LastAccessDateTime: utils.ToTime(sw[i]["last_access_datetime"]),
		}
	}
	return insSw
}

func getKeyboards(data map[string]interface{}) []asset.Keyboard {
	if data["keyboards"] == nil {
		return nil
	}
	kb := data["keyboards"].([]map[string]interface{})
	l := len(kb)
	keyboards := make([]asset.Keyboard, l)
	for i := 0; i < l; i++ {
		keyboards[i] = asset.Keyboard{
			DeviceID:    utils.ToString(kb[i]["device_id"]),
			Name:        utils.ToString(kb[i]["name"]),
			Description: utils.ToString(kb[i]["description"]),
		}
	}
	return keyboards
}

func getMouse(data map[string]interface{}) []asset.Mouse {

	if data["mouse"] == nil {
		return nil
	}
	mouse := data["mouse"].([]map[string]interface{})
	l := len(mouse)
	assetMouse := make([]asset.Mouse, l)
	for i := 0; i < l; i++ {
		assetMouse[i] = asset.Mouse{
			Name:            utils.ToString(mouse[i]["name"]),
			Buttons:         utils.ToInt(mouse[i]["buttons"]),
			DeviceID:        utils.ToString(mouse[i]["device_id"]),
			DeviceInterface: utils.ToInt(mouse[i]["device_interface"]),
			Manufacturer:    utils.ToString(mouse[i]["manufacturer"]),
			PointingType:    utils.ToInt(mouse[i]["pointing_type"]),
		}
	}
	return assetMouse
}

func getMonitors(data map[string]interface{}) []asset.Monitor {
	if data["monitors"] == nil {
		return nil
	}
	monitorArr := data["monitors"].([]map[string]interface{})
	monitors := make([]asset.Monitor, len(monitorArr))
	for i, value := range monitorArr {
		monitors[i].Name = utils.ToString(value["name"])
		monitors[i].DeviceID = utils.ToString(value["device_id"])
		monitors[i].Manufacturer = utils.ToString(value["manufacturer"])
		monitors[i].ScreenHeight = uint32(utils.ToInt(value["screen_height"]))
		monitors[i].ScreenWidth = uint32(utils.ToInt(value["screen_width"]))
	}
	return monitors
}

func getBaseBoard(data map[string]interface{}) asset.AssetBaseBoard {
	if data["baseboard"] == nil {
		return asset.AssetBaseBoard{}
	}
	baseBoardMap := data["baseboard"].(map[string]interface{})
	if len(baseBoardMap) <= 0 {
		return asset.AssetBaseBoard{}
	}

	var bb asset.AssetBaseBoard
	bb.Product = utils.ToString(baseBoardMap["product"])
	bb.Manufacturer = utils.ToString(baseBoardMap["manufacturer"])
	bb.Model = utils.ToString(baseBoardMap["model"])
	bb.SerialNumber = utils.ToString(baseBoardMap["serial_number"])
	bb.Version = utils.ToString(baseBoardMap["version"])
	bb.Name = utils.ToString(baseBoardMap["name"])
	bb.InstallDate = utils.ToTime(baseBoardMap["install_date"])
	return bb
}

//SaveAssets to parse Asset data to the DB data model and persist into DB
func SaveAssets(assetData *asset.AssetCollection) error {
	currUTCTime := time.Now().UTC()
	errForMap := Session.Query(insertEndpointMappingQuery, assetData.PartnerID, assetData.EndpointID, assetData.SiteID, assetData.ClientID, assetData.RegID, time.Now().UTC(), time.Now().UTC()).Exec()
	if errForMap != nil {
		return fmt.Errorf("Error Occurred while updating the mapping for endpointID : %s, partnerId : %s, err : %+v", assetData.EndpointID, assetData.PartnerID, errForMap)
	}
	err := addSummary(assetData, currUTCTime, currUTCTime)
	if err != nil {
		return fmt.Errorf("Error Occurred while saving the summary for endpointID : %s, partnerId : %s, err : %+v", assetData.EndpointID, assetData.PartnerID, err)
	}
	errForDetails := addDetails(assetData, currUTCTime, currUTCTime)
	if errForDetails != nil {
		return fmt.Errorf("Error Occurred while saving the Details for endpointID : %s, partnerId : %s, err : %+v", assetData.EndpointID, assetData.PartnerID, errForDetails)
	}
	return nil
}
func addSummary(asset *asset.AssetCollection, currUTCTime, agentUTCTime time.Time) error {
	return Session.Query(insertEndpointSummaryQuery,
		asset.PartnerID,
		asset.SiteID,
		asset.ClientID,
		asset.EndpointID,
		asset.RegID,
		asset.Name,
		asset.FriendlyName,
		asset.RemoteAddress,
		asset.Type,
		asset.ResourceType,
		asset.EndpointType,
		asset.System,
		asset.Os,
		asset.CreatedBy,
		agentUTCTime,
		currUTCTime,
	).Exec()
}

func addDetails(asset *asset.AssetCollection, currUTCTime, agentUTCTime time.Time) error {
	return Session.Query(insertEndpointDetailsQuery, asset.PartnerID, asset.SiteID, asset.EndpointID,
		asset.BaseBoard,
		asset.Bios,
		asset.Drives,
		asset.Memory,
		asset.Networks,
		asset.Processors,
		asset.RaidController,
		asset.InstalledSoftwares,
		asset.Keyboards,
		asset.Mouse,
		asset.Monitors,
		asset.PhysicalDrives,
		asset.Users,
		asset.Shares,
		asset.SoftwareLicenses,
		agentUTCTime,
		currUTCTime,
	).Exec()
}
