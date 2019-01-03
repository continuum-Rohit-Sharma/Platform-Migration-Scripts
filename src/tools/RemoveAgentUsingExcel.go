package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gocql/gocql"
// 	"github.com/xuri/excelize"
// )

// type assetCollection struct {
// 	EndPointList []EndPointList `json:"endPointList"`
// 	TotalPoint   int            `json:"totalCount"`
// }

// //EndPointList is the structuire
// type EndPointList struct {
// 	RegID                 string `json:"regID,omitempty"`
// 	FriendlyName          string `json:"friendlyName"`
// 	MachineName           string `json:"machineName"`
// 	SiteName              string `json:"siteName"`
// 	OperatingSystem       string `json:"operatingSystem"`
// 	Availability          int    `json:"availability"`
// 	IPAddress             string `json:"ipAddress"`
// 	RegType               string `json:"regType"`
// 	LmiStatus             int    `json:"lmiStatus"`
// 	ResType               string `json:"resType"`
// 	SiteID                int    `json:"siteId"`
// 	SmartDisk             int    `json:"smartDisk"`
// 	Amt                   int    `json:"amt"`
// 	MbSyncstatus          int    `json:"mbSyncstatus"`
// 	EncryptedResourceName string `json:"encryptedResourceName"`
// 	EncryptedSiteName     string `json:"encryptedSiteName"`
// 	LmiHostID             string `json:"lmiHostId"`
// 	RequestRegID          string `json:"requestRegId"`
// 	EndpointID            string `json:"endpointId"`
// }
// type fnEndpoint struct {
// 	PartnerID     string
// 	EndpointID    string `json:"endpointId"`
// 	SystemName    string `json:"systemName"`
// 	BiosSerialNum string `json:"serialNumber"`
// 	RegID         string `json:"regID"`
// 	ToDelete      bool
// }

// type resultValue struct {
// 	partnerID string
// 	endpoints string
// }

// //Session session variable
// var Session *gocql.Session
// var clusterIP string
// var newWebAPIURL, getAssetURL string
// var totalNoPartnersFromCSV, totalNoPartners int
// var mapSystemNameFnEndpoint map[string][]fnEndpoint
// var duplicateEndpoints [][]string
// var mapDuplicateEndpoints map[string][][]string
// var doDelete bool
// var totalCount, totalPartners, noOfWorkers int
// var configObject Config

// type Config struct {
// 	PrintResult         bool
// 	RequestWaitInSecond int64
// 	ServerAddress       string
// 	NewWebAPIURL        string
// 	IsDelete            bool
// 	NumOfWorkers        int
// 	CassandraHost       string
// 	CSVPathForPartners  string
// 	//Messages            []MailboxMessage
// }

// const (
// 	queryToDeleteEndpoint                   = `DELETE FROM platform_asset_db.partner_asset_collection where partner_id= '%s' AND endpoint_id = '%s'`
// 	queryToDeleteEndpointsFromAgentVersions = `DELETE FROM platform_asset_db.agent_version_details where partner_id= '%s' AND endpoint_id = '%s'`
// 	cSelectMultipleEndpointQuery            = `select partner_id, endpoint_id,  DcDateTimeUTC from platform_agent_db.agent_heartbeat where partner_id ='%s' and endpoint_id in (%s)`
// )

// //LoadConfiguration is a method to load configuration File
// func LoadConfiguration(filePath string) (Config, error) {
// 	mbMessage := &Config{}
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return *mbMessage, err
// 	}
// 	defer file.Close()

// 	deser := json.NewDecoder(file)
// 	err = deser.Decode(mbMessage)
// 	if err != nil {
// 		return *mbMessage, err
// 	}
// 	return *mbMessage, nil
// }
// func main() {
// 	commandArgs := os.Args[1:]
// 	if len(commandArgs) != 1 {
// 		fmt.Println("Usage: removeDuplicates configuration.json")
// 		fmt.Println("Example: removeDuplicates config.json")
// 		return
// 	}

// 	startTime := time.Now()

// 	fmt.Println("Migration Tool Started... Time started : ", startTime)
// 	configObject, err := LoadConfiguration(os.Args[1])
// 	if err != nil {
// 		fmt.Println("Error Occured while Loading the Configuration.", err)
// 	}
// 	excelPath := configObject.CSVPathForPartners
// 	newWebAPIURL = configObject.NewWebAPIURL
// 	clusterIP = configObject.CassandraHost
// 	noOfWorkers = configObject.NumOfWorkers
// 	doDelete = configObject.IsDelete

// 	// if newWebAPIURL != "" {
// 	// 	newWebAPIURL = newWebAPIURL + "/asset/v1/partner/%s/endpoints/%s?field=bios&field=system"
// 	// 	//fmt.Printf("New System Service URL %s\n", newWebAPIURL)
// 	// }
// 	cassError := getCassandraSession(clusterIP)
// 	if cassError != nil {
// 		fmt.Println("Error Occured while setting up the Tool with Cassandra, Error : ", cassError)
// 		return
// 	}
// 	fnEndpoints, err := readEndpoints(excelPath)
// 	if err != nil {
// 		fmt.Println("Error Occured while getting the partners from Excel, Error : ", err)
// 		return
// 	}
// 	// totalCountOfEndpoints := len(fnEndpoints)
// 	// jobs := make(chan fnEndpoint, totalCountOfEndpoints)
// 	// results := make(chan fnEndpoint, totalCountOfEndpoints)
// 	for _, endpoint := range fnEndpoints {
// 		if endpoint.EndpointID != "" && endpoint.PartnerID != "" {
// 			deleteDupEndpoints(endpoint.EndpointID, endpoint.PartnerID)
// 			deleteDupEndpointsFromAgent(endpoint.EndpointID, endpoint.PartnerID)
// 		}

// 		// getEndpointByEndpointID(endpoint)

// 	}
// 	// for w := 0; w <= totalCountOfEndpoints; w++ {
// 	// 	 getEndpointByEndpointID()
// 	// }
// 	// fmt.Println("# of workers : ", noOfWorkers)
// 	// for a := 0; a < totalCountOfEndpoints; a++ {
// 	// 	jobs <- fnEndpoints[a]
// 	// }
// 	// fmt.Println("# of jobs ", totalCountOfEndpoints)
// 	// close(jobs)
// 	// fmt.Println("Waiting for results ")
// 	// //ptrRes := make([]PartnerMetrics, len(uniquePartnerIDs))
// 	// var a int
// 	// endpoints := make([]fnEndpoint, 0)
// 	// for a = 0; a < totalCountOfEndpoints; a++ {
// 	// 	endpoint := <-results
// 	// 	if endpoint.ToDelete {
// 	// 		endpoints = append(endpoints, endpoint)
// 	// 	}
// 	// }
// 	// endTime := time.Since(startTime)
// 	// fmt.Println("Done processing results. Total recieved : ", a)
// 	// fmt.Println("Identification completed in : ", endTime)
// 	// fmt.Printf("\n\n*** Duplicate endpoints Metrics ***\n")
// 	// fmt.Printf("\n | Partner ID | Endpoint ID \n")
// 	// for _, endpoint := range endpoints {
// 	// 	fmt.Printf("\n | %s | %s ", endpoint.PartnerID, endpoint.EndpointID)
// 	// }
// 	// if doDelete && len(endpoints) > 0 {
// 	// 	fmt.Println("Came for Deletion of the endpoints")
// 	// 	for _, endpoint := range endpoints {
// 	// 		deleteDupEndpoints(endpoint.EndpointID, endpoint.PartnerID)
// 	// 	}
// 	// } else {
// 	// 	fmt.Println("Nothing to delete")
// 	// }
// 	// close(results)
// 	Session.Close()
// }

// func getCassandraSession(clusterIP string) error {
// 	var err error
// 	cluster := gocql.NewCluster(clusterIP)
// 	cluster.Consistency = gocql.Quorum
// 	cluster.Keyspace = "platform_asset_db"
// 	Session, err = cluster.CreateSession()
// 	if err != nil {
// 		fmt.Printf("Error occured : %v", err)
// 		return err
// 	}
// 	return nil
// }

// type hbEndpoint struct {
// 	EndpointID    string
// 	PartnerID     string
// 	DcDateTimeUTC time.Time
// }

// //Heartbeat struct for capturing heartbeat
// type Heartbeat struct {
// 	RegID            string
// 	DcDateTimeUTC    time.Time
// 	AgentDateTimeUTC time.Time
// 	HeartbeatCounter int64
// 	EndpointID       string
// 	installed        bool
// }

// func getEndpointByEndpointID(endpoint fnEndpoint) {
// 	statusCode := 0
// 	partnerID := strings.TrimSpace(endpoint.PartnerID)
// 	endpointID := strings.TrimSpace(endpoint.EndpointID)
// 	requestURL := fmt.Sprintf(newWebAPIURL, partnerID, endpointID)
// 	fmt.Println("Request URL : ", requestURL)
// 	res, err := http.Get(requestURL)
// 	if err != nil {
// 		fmt.Printf("\n Error occurred for Partner Id : %s & endpoint Id : %s \n", partnerID, endpointID)
// 		return
// 	} else {
// 		if res != nil {
// 			statusCode = res.StatusCode
// 		}
// 	}
// 	defer res.Body.Close()
// 	fmt.Println("statusCode : ", statusCode)

// 	if statusCode == 404 {
// 		fmt.Printf(" \n Deleting for endpoint id : %s and partner ID : %s", endpoint.EndpointID, endpoint.PartnerID)
// 		deleteDupEndpoints(endpoint.EndpointID, endpoint.PartnerID)
// 	}
// 	return
// }

// func deleteDupEndpoints(endpointID string, partnerID string) error {
// 	fmt.Printf(" \n Deleting the Assets for Endpoint id : %s and partner Id : %s ", endpointID, partnerID)
// 	query := fmt.Sprintf(queryToDeleteEndpoint, partnerID, endpointID)
// 	err := Session.Query(query).Exec()
// 	return err
// }

// func deleteDupEndpointsFromAgent(endpointID string, partnerID string) error {
// 	fmt.Printf("\n Deleting the Agent version for Endpoint id from agent versions: %s", endpointID)
// 	//sEndpointIDs := strings.Join(endpointIDs, "','")
// 	//sEndpointIDs = fmt.Sprintf("'%s'", sEndpointIDs)
// 	query := fmt.Sprintf(queryToDeleteEndpointsFromAgentVersions, partnerID, endpointID)
// 	err := Session.Query(query).Exec()
// 	return err
// }
// func readEndpoints(sheetPath string) ([]fnEndpoint, error) {
// 	endpoints := make([]fnEndpoint, 0)
// 	xlFile, err := excelize.OpenFile(sheetPath)
// 	if err != nil {
// 		fmt.Println("Error opening excel sheet")
// 		return endpoints, err
// 	}
// 	rows := xlFile.GetRows("Sheet1")
// 	for cnt, row := range rows {
// 		if cnt == 0 {
// 			continue
// 		}
// 		endpoint := fnEndpoint{
// 			EndpointID: row[2],
// 			PartnerID:  row[3],
// 			RegID:      row[1],
// 		}
// 		endpoints = append(endpoints, endpoint)
// 	}
// 	fmt.Println(endpoints)
// 	return endpoints, err
// }
