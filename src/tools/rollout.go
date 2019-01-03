package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/xuri/excelize"
)

var (
	configObject Config
	logger       *log.Logger
	noOfWorkers  int
	//Session is gocql session used across the queries
	Session *gocql.Session
)

//Config Struct to hold Configuration for Load Generation
type Config struct {
	PrintResult         bool
	RequestWaitInSecond int64
	ServerAddress       string
	NumOfWorkers        int
	CSVPathForPartners  string
	Offset              int
	Messages            []MailboxMessage
	ManifestVersion     string
	AgentCassandraDB    string
	LogFile             string
}

// Mailbox ...
type Mailbox struct {
	Endpoints []string
	SourceURL string
	Message   MailboxMessage
}

//MailboxMessage ...
type MailboxMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Path         string    `json:"path"`
	Message      string    `json:"message"`
}

const (
	insertPartnerManifest = "INSERT INTO platform_agent_db.partner_manifest (partner_id,manifest_version, dc_created_ts_utc) VALUES(?,?,?);"
)

func main() {
	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: RollOut config_rollout.json")
		return
	}
	configObject, err := LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}
	excelPath := configObject.CSVPathForPartners
	noOfWorkers = configObject.NumOfWorkers
	offset := configObject.Offset
	mailboxMsgs := configObject.Messages
	clusterIP := configObject.AgentCassandraDB
	manifestVersion := configObject.ManifestVersion
	logFile := configObject.LogFile
	var logNo string
	if offset == -1 {
		logNo = "0"
	} else {
		logNo = strconv.Itoa(offset)

	}

	logFileName := fmt.Sprintf(logFile, logNo)
	logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("\n Error occurred while creating the log file\n Err  : %v", err)
		return
	}
	logger = log.New(logfile, "", logFlags)
	startTime := time.Now()
	logger.Println("Tool Started... Time started : ", startTime)
	///////////////////////////////////////////////////////////////////////
	err = getCassandraSession(clusterIP)
	if err != nil {
		logger.Printf("\n Error occurred while creating the Cassandra session for IP : %s, Error is : %v\n", clusterIP, err)
		return
	}
	endpoints, err := readEndpoints(excelPath)
	if err != nil {
		logger.Printf("\n Error occurred while reading the endpoints from excel \n")
	}
	if len(endpoints) == 0 {
		logger.Printf("\n No endpoints got read from Excel \n")
	}
	totalNoEndpoints := len(endpoints)

	if offset+1 >= totalNoEndpoints {
		logger.Println("No Endpoints to process. So exiting.")
		return
	}

	msg := mailboxMsgs[0]
	msg.TimestampUTC = time.Now()
	endpointIds := make([]string, 0)
	partnerIDs := make(map[string]string)
	newOffset := 0
	for cnt, endpoint := range endpoints {
		if (cnt > offset) && (cnt <= offset+noOfWorkers) {
			endpointIds = append(endpointIds, endpoint.endpointID)
			partnerIDs[endpoint.partnerID] = endpoint.partnerID
			newOffset = cnt
		}
	}
	//fmt.Println("map : ", partnerIDs)
	for partnerID, _ := range partnerIDs {
		err = saveManifestVersionForPartner(manifestVersion, partnerID)
		if err != nil {
			logger.Printf("\n Error occurred while saving the mainfest in DB for partner Id : %s and error is : %v\n", partnerID, err)
			return
		}

	}
	err = sendMessages(endpointIds, msg, configObject)
	if err != nil {
		logger.Printf("\n Error Occurred while sending the message and error is : %v\n", err)
		return

	}
	logger.Println("Endpoints Taken for Autoupdate : ", endpointIds)
	configObject.Offset = newOffset
	logger.Println("Updating the log file with new offset with value : ", newOffset)
	errWritingFile := WriteFile("config_rollout.json", &configObject)
	if errWritingFile != nil {
		logger.Printf("\n Error occurred while writing the file, error is : %v \n", errWritingFile)
	}
	///////////////////////////////////////
	endTime := time.Since(startTime)
	logger.Printf("\nRollout completed in for %d to %d in time : %v \n", offset, newOffset-1, endTime)

}

func saveManifestVersionForPartner(manifestVersion string, partnerID string) error {
	if manifestVersion == "" || partnerID == "" {
		return errors.New("Blank manifest or partner ID")
	}
	partnerID = strings.TrimSpace(partnerID)
	//query := fmt.Sprintf(insertPartnerManifest, partnerID)
	err := Session.Query(insertPartnerManifest, partnerID, manifestVersion, time.Now().UTC()).Exec()
	return err
}

// //AutoUpdateMailboxMsg is mailbox message for auto-update
// var AutoUpdateMailboxMsg = apiModel.MailboxMessage{
// 	Path:    "",
// 	Message: `{"action": "update","appname": "JunoAgent"}`,
// 	Name:    "JunoAutoupdate",
// 	Type:    "APP",
// 	Version: "1.0",
// }

func sendMessages(endpoints []string, mbMsg MailboxMessage, configObject Config) error {
	msg := Mailbox{
		Endpoints: endpoints,
		Message:   mbMsg,
	}
	resp := make([]string, 0)
	encoded, err := json.Marshal(msg)
	if err != nil {
		logger.Println("Error Occurred while marshalling the mailbox message")
		return err
	}
	payload := strings.NewReader(string(encoded))
	req, err := http.NewRequest("POST", configObject.ServerAddress, payload)
	if err != nil {
		logger.Println(time.Now(), " Error while creating request ", err, " for endpoints ", endpoints)
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	logger.Println("Message : ", msg)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(time.Now(), " Error while geting response ", err, " for endpoints ", endpoints)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println(time.Now(), " Error while Reading message ", err, " for endpoints ", endpoints)
		return err
	}
	resp = append(resp, string(body))
	logger.Println("Response for the endpoints is : ", resp)
	return nil
}

//WriteFile is to write the file at the desired location
func WriteFile(filePath string, tObject interface{}) error {
	if filePath == "" {
		logger.Println("Error Occurred while opening the file ; ", filePath)
		return errors.New("errorWhileOpening")
	}

	fp, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error Occurred while creating the file ; ", filePath)
		return err
	}
	return Serialize(fp, tObject)
}

//Serialize serializes the given interface
func Serialize(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	err := enc.Encode(v)
	if err != nil {
		return err
	}
	return nil
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

type endpointCollection struct {
	partnerID, endpointID, regID string
}

func create(line []string) endpointCollection {

	return endpointCollection{
		partnerID:  line[0],
		regID:      line[1],
		endpointID: line[2],
	}
}

func readEndpoints(filePath string) ([]endpointCollection, error) {

	endpointCollections := make([]endpointCollection, 0)
	xlFile, err := excelize.OpenFile(filePath)
	if err != nil {
		logger.Println("Error opening excel sheet, Err : ", err)
		return endpointCollections, err
	}
	rows := xlFile.GetRows("Sheet1")
	for _, row := range rows {
		astCol := create(row)
		if astCol.endpointID != "" {
			endpointCollections = append(endpointCollections, astCol)
		}
	}
	return endpointCollections, err

}

func getCassandraSession(clusterIP string) error {
	var err error
	cluster := gocql.NewCluster(clusterIP)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "platform_agent_db"
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Printf("Error occured : %v", err)
		return err
	}
	return nil
}

const logFlags int = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile
