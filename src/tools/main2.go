// package main

// import (
// 	"bufio"
// 	"encoding/csv"
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"time"

// 	"fmt"
// 	"os"

// 	"github.com/gocql/gocql"
// )

// type assetCollection struct {
// 	EndPointList []EndPointList `json:"endPointList"`
// 	TotalPoint   int            `json:"totalCount"`
// }

// //EndPointList is the structuire
// type EndPointList struct {
// 	RegID        string `json:"regID,omitempty"`
// 	FriendlyName string `json:"friendlyName"`
// 	MachineName  string `json:"machineName"`
// 	SiteName     string `json:"siteName"`
// 	//OperatingSystem       string `json:"operatingSystem"`
// 	Availability int `json:"availability"`
// 	//IPAddress             string `json:"ipAddress"`
// 	RegType string `json:"regType"`
// 	//LmiStatus             int    `json:"lmiStatus"`
// 	//ResType               string `json:"resType"`
// 	SiteID int `json:"siteId"`
// 	//SmartDisk             int    `json:"smartDisk"`
// 	//Amt                   int    `json:"amt"`
// 	//MbSyncstatus          int    `json:"mbSyncstatus"`
// 	//EncryptedResourceName string `json:"encryptedResourceName"`
// 	//EncryptedSiteName     string `json:"encryptedSiteName"`
// 	// LmiHostID             string `json:"lmiHostId"`
// 	//RequestRegID string `json:"requestRegId"`
// 	EndpointID string `json:"endpointId"`
// }
// type Config struct {
// 	NewWebAPIURL            string
// 	OldWebAPIURL            string
// 	IsSite                  bool
// 	PartnerID               string
// 	SiteIDs                 string
// 	NumOfWorkers            int
// 	CassandraHostIPForAsset string
// 	CSVPathForPartners      string
// 	ResourceType            string
// }

// //Session session variable
// var Session *gocql.Session
// var clusterIP string
// var legacyWebAPIURL, newWebAPIURL string
// var totalNoPartnersFromCSV, totalNoPartners int
// var PartnerEndpoints map[string]map[string]fnEndpoint

// //var SiteIDs string

// const (
// 	queryToUpdateFriendlyName     = "UPDATE platform_asset_db.partner_asset_collection set friendly_name=? where partner_id=? AND endpoint_id=?"
// 	queryToGetEndpointBypartnerID = "SELECT endpoint_id, reg_id, friendly_name FROM  platform_asset_db.partner_asset_collection WHERE partner_id='%s'"
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

// 	startTime := time.Now()

// 	commandArgs := os.Args[1:]
// 	if len(commandArgs) != 1 {
// 		fmt.Println("Usage: YourScriptName <config.json>")
// 		fmt.Println("Example: runJob config.json")
// 		return
// 	}

// 	fmt.Println("Tool Started... Time started : ", startTime)

// 	configObject, err := LoadConfiguration(os.Args[1])
// 	if err != nil {
// 		fmt.Println("Error Occured while Loading the Configuration.", err)
// 	}

// 	fmt.Println("Migration Tool Started... Time started : ", startTime)
// 	excelPath := configObject.CSVPathForPartners
// 	noOfJobs := configObject.NumOfJobs

// 	endpoints, err := readAssetTables(excelPath)
// 	if err != nil {
// 		fmt.Println("Error Occured while readAssetTables, Error : ", err)
// 		return
// 	}
// 	jobs := make(chan string)
// 	results := make(chan PartnerMetrics, totalNoPartnersUnique)

// 	for w := 1; w <= noOfWorkers; w++ {
// 		go processParters(w, jobs, results)
// 	}
// 	fmt.Println("# of workers : ", noOfWorkers)

// 	for a := 0; a < len(uniquePartnerIDs); a++ {
// 		jobs <- uniquePartnerIDs[a]
// 	}
// 	fmt.Println("# of jobs ", totalNoPartnersUnique)
// 	close(jobs)

// 	fmt.Println("Waiting for results ")
// 	ptrRes := make([]PartnerMetrics, len(uniquePartnerIDs))
// 	var a int
// 	for a = 0; a < len(uniquePartnerIDs); a++ {
// 		ptrRes[a] = <-results
// 	}
// 	endTime := time.Since(startTime)
// 	fmt.Println("Done processing results. Total recieved : ", a)
// 	fmt.Println("Migration completed in : ", endTime)
// 	fmt.Printf("\n\n*** Partner Level Metrics ***\n\n")
// 	for a = 0; a < len(uniquePartnerIDs); a++ {
// 		fmt.Printf("Partner ID %s: RMM1Err: %s, RMM2Err: %s, NoEndpointsInRMM1: %v, NOE already in Sync: %d, NOE updated: %d, NOE with UpdateDBError: %d, NOE missing in RMM2: %d \n", ptrRes[a].ID, ptrRes[a].RMM1Err, ptrRes[a].RMM2Err, ptrRes[a].RMM1NoEndpointPartner, len(ptrRes[a].FNameInSyncEndpoints), len(ptrRes[a].UpdatedEndpoints), len(ptrRes[a].UpdateDBErrorEndpoints), len(ptrRes[a].RMM2NoEndpoints))
// 	}
// 	close(results)
// 	Session.Close()

// }

// func processParters(id int, jobs <-chan string, results chan<- PartnerMetrics) {

// 	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM1.0 errors: %d, RMM2.0 errors: %d, GoodCandidates: %d, Partners with no endpoints in RMM 1.0: %d, Common Partners with endpoints: %d\n", id, rmmoneerr, rmmtwoerr, goodones, unknwonptr, (goodones - unknwonptr))
// 	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM1.0 GetLegacyEndpointlist total time :  %f , instances: %d,  avg per partner: %f\n", id, rmm1tt, rmm1cntr, rmm1tt/float64(rmm1cntr))
// 	fmt.Printf("Goroutine id# %d, PerThreadPartnerMetrics : RMM2.0 GetJunoEndpointlist   total time :  %f , instances: %d,  avg per partner: %f\n", id, rmm2tt, rmm2cntr, rmm2tt/float64(rmm2cntr))

// 	fmt.Printf("Goroutine id# %d, PerThreadEndpointMetrics : Already in Sync: %d, DB Error : %d, EndpointNotfound : %d,  processed : %d \n", id, inSync, dbErr, endpointNotFoundCounter, counter)
// 	fmt.Printf("Goroutine id# %d, PerThreadEndpointMetrics : DBupdate total time :  %f , instances: %d,  avg per partner: %f\n", id, dbupdatett, dbupdatecntr, dbupdatett/float64(dbupdatecntr))
// }

// // func removeDuplicatePartnerIds(partnerIDs []string) ([]string, error) {
// // 	var err error
// // 	if len(partnerIDs) < 1 {
// // 		return partnerIDs, err
// // 	}

// // 	// Use map to record duplicates as we find them.
// // 	encountered := map[string]bool{}
// // 	result := []string{}

// // 	for v := range partnerIDs {
// // 		if encountered[partnerIDs[v]] == true {
// // 		} else {
// // 			encountered[partnerIDs[v]] = true
// // 			partnerID := strings.TrimSpace(partnerIDs[v])
// // 			result = append(result, partnerID)
// // 		}
// // 	}
// // 	return result, nil
// // }

// type asssetMigCollection struct {
// 	partnerID, endpointID, regID string
// }

// func create(line []string) asssetMigCollection {
// 	return asssetMigCollection{
// 		partnerID:  line[0],
// 		regID:      line[1],
// 		endpointID: line[2],
// 	}
// }

// func readAssetTables(filePath string) (asssetMigCollection, error) {
// 	astCol := asssetMigCollection{}
// 	csvFile, ferr := os.Open(filePath)
// 	if ferr != nil {
// 		return astCol, ferr
// 	}
// 	reader := csv.NewReader(bufio.NewReader(csvFile))
// 	var err error
// 	for {
// 		line, error := reader.Read()
// 		if error == io.EOF {
// 			break
// 		} else if error != nil {
// 			err = error
// 			log.Fatal(error)
// 		}
// 		astCol := create(line)
// 	}
// 	return astCol, err
// }
