package main

// This tool is created for Migration of data from old cassandra tables to new one.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/utils"
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
	CassandraHostForAsset string
	CSVPath               string
	LogFilePath           string
	NoOfWorkers           int
	FailureLogFilePath    string
}

const (
	resourceTypeDesktop               = "Desktop"
	resourceTypeServer                = "Server"
	osTypeWindows                     = "Windows"
	cEndpointIDSelectQuery            = "SELECT partner_id, client_id, site_id, endpoint_id,reg_id, friendly_name, agent_ts_utc, created_by,name,type,remote_address,resource_type, endpoint_type, baseboard, bios, drives, physical_memory, networks, os, processors, raidcontroller, system, installed_softwares, keyboards, mouse, monitors, physical_drives,users,services,shares,software_licenses FROM platform_asset_db.partner_asset_collection WHERE partner_id=? and endpoint_id=?"
	insertEndpointSummaryQuery        = "INSERT INTO platform_asset_db.endpoint_summary(partner_id, site_id, client_id,endpoint_id,reg_id, name,friendly_name,remote_address, type, resource_type, endpoint_type, system, os,created_by, agent_ts_utc,dc_ts_utc) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	insertEndpointMappingQuery        = "INSERT INTO platform_asset_db.endpoint_mapping((partner_id,endpoint_id,site_id, client_id,reg_id,dc_ts_utc,agent_ts_utc) VALUES(?,?,?,?,?,?,?)"
	insertEndpointServiceDetailsQuery = "INSERT INTO platform_asset_db.endpoint_service_details(partner_id,endpoint_id,service_name,details,dc_ts_utc,agent_ts_utc) VALUES(?,?,?,?,?,?)"
	insertEndpointDetailsQuery        = "INSERT INTO platform_asset_db.endpoint_details (partner_id,endpoint_id,baseboard,bios,drives,physical_memory,networks,processors,raidcontroller,installed_softwares, keyboards, mouse, monitors, physical_drives,users,shares,software_licenses) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

func main() {
	commandArgs := os.Args[1:]
	if len(commandArgs) != 1 {
		fmt.Println("Usage: YourScriptName <config.json>")
		fmt.Println("Example: migrate config.json")
		return
	}
	configObject, err := LoadConfiguration(os.Args[1])
	if err != nil {
		fmt.Println("Error Occured while Loading the Configuration.", err)
	}
	excelPath := configObject.CSVPath
	clusterIP := configObject.CassandraHostForAsset
	logFile := configObject.LogFilePath
	noOfWorkers := configObject.NoOfWorkers
	failureLogFilePath := configObject.FailureLogFilePath
	logfile, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
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
	jobs := make(chan endpointCollection, totalNoEndpoints)
	results := make(chan endpointCollection, totalNoEndpoints)
	for w := 1; w <= noOfWorkers; w++ {
		go migrateEndpoints(w, jobs, results)
	}
	var a int
	for a := 0; a < totalNoEndpoints; a++ {
		jobs <- endpoints[a]
	}

	close(jobs)
	var success, failure int
	failureEndpoints := make([]endpointCollection, 0)
	for a = 0; a < totalNoEndpoints; a++ {
		endpoint := <-results
		if endpoint.err != nil {
			log.Printf("Error Occurred : %v for partner Id : %s and endpoint : %s", endpoint.err, endpoint.partnerID, endpoint.endpointID)
			failure++
			failureEndpoints = append(failureEndpoints, endpoint)
		} else {
			success++
		}
	}
	close(results)
	createReport(failureEndpoints, failureLogFilePath)

}

func createReport(endpoints []endpointCollection, logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("\n Error occurred while creating the failure log file\n Err  : %v", err)
		return
	}
	logger = log.New(logfile, "", logFlags)
	logger.Printf("PARTNERID||ENDPOINTID||ERRORS \n")
	for _, endpoint := range endpoints {
		logger.Printf("\n%s||%s||%+v\n", endpoint.partnerID, endpoint.endpointID, endpoint.err)
	}
}
func migrateEndpoints(id int, jobs chan endpointCollection, results chan endpointCollection) {
	for endpoint := range jobs {
		asset, err := getEndpointDetails(endpoint)
		if err != nil {
			endpoint.err = err
			results <- endpoint
			continue
		}
		err = saveDetails(asset)
		if err != nil {
			endpoint.err = err
			results <- endpoint
			continue
		}
		results <- endpoint
	}
}

func getEndpointDetails(endpoint endpointCollection) (*asset.AssetCollection, error) {
	//asset := asset.AssetCollection{}
	if endpoint.endpointID != "" && endpoint.partnerID != "" {
		record, err := Select(cEndpointIDSelectQuery, endpoint.partnerID, endpoint.endpointID)
		if err != nil {
			return nil, err
		}
		if len(record) <= 0 {
			return nil, nil
		}
		//only one asset record per endpointId is maintained
		asset := dbMapToDBModel(record[0])
		return asset, nil
	}
	logger.Printf("\n Error occurred while getting details endpointId = %s and partnerID = %s \n", endpoint.endpointID, endpoint.partnerID)
	return nil, errors.New("Error_PARTNERID_ENDPOINTID_BLANK")
}

func saveDetails(asset *asset.AssetCollection) error {
	err := saveAssetSummary(asset)
	if err != nil {
		return err
	}
	var errString string
	err = saveEndpointMapping(asset)
	if err != nil {
		errString = errString + "\n" + err.Error()
	}
	err = saveEndpointDetails(asset)
	if err != nil {
		errString = errString + "\n" + err.Error()
	}
	saveEndpointServices(asset.PartnerID, asset.EndpointID, asset.Services)
	if errString != "" {
		return errors.New(errString)
	}
	return nil
}
func saveAssetSummary(asset *asset.AssetCollection) error {
	return Session.Query(insertEndpointSummaryQuery,
		asset.PartnerID,
		asset.EndpointID,
		asset.SiteID,
		asset.ClientID,
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
		time.Now().UTC(),
		time.Now().UTC(),
	).Exec()
}

func saveEndpointMapping(asset *asset.AssetCollection) error {
	return Session.Query(insertEndpointMappingQuery,
		asset.PartnerID,
		asset.EndpointID,
		asset.SiteID,
		asset.ClientID,
		asset.RegID,
		time.Now().UTC(),
		time.Now().UTC(),
	).Exec()
}

func saveEndpointDetails(asset *asset.AssetCollection) error {
	return Session.Query(insertEndpointDetailsQuery,
		asset.PartnerID,
		asset.EndpointID,
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
	).Exec()
}

func saveEndpointServices(partnerID string, endpointID string, services []asset.Service) {

	if len(services) < 1 {
		return
	}
	for i := 0; i <= len(services); i++ {
		service := services[i]
		svcName := service.Name
		svcDetails := asset.ServiceDetails{
			DisplayName:             service.Details.DisplayName,
			ExecutablePath:          service.Details.ExecutablePath,
			StartupType:             service.Details.StartupType,
			ServiceStatus:           service.Details.ServiceStatus,
			LogOnAs:                 service.Details.LogOnAs,
			StopEnableAction:        service.Details.StopEnableAction,
			DelayedAutoStart:        service.Details.DelayedAutoStart,
			Win32ExitCode:           service.Details.Win32ExitCode,
			ServiceSpecificExitCode: service.Details.ServiceSpecificExitCode,
		}
		Session.Query(insertEndpointServiceDetailsQuery, partnerID, endpointID, svcName, svcDetails, time.Now().UTC(), time.Now().UTC()).Exec()
	}

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
	partnerID  string
	endpointID string
	err        error
}

func create(line []string) endpointCollection {

	return endpointCollection{
		partnerID:  line[0],
		endpointID: line[1],
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
	cluster.Keyspace = "platform_asset_db"
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Printf("Error occured : %v", err)
		return err
	}
	return nil
}

func dbMapToDBModel(data map[string]interface{}) *asset.AssetCollection {
	assetData := &asset.AssetCollection{}
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

	assetData.BaseBoard = getBaseBoard(data)
	assetData.Bios = getBios(data)
	assetData.Os = getOS(data)
	assetData.RaidController = getRaidCtrlr(data)
	assetData.System = getSystem(data)
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
	assetData.Services = getServices(data)
	assetData.Shares = getShares(data)
	assetData.SoftwareLicenses = getSoftwareLicenses(data)

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
	return assetData
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
			UserAccess:  utils.ToStringArray(u[i]["user_access"]),
			Type:        utils.ToStringArray(u[i]["type"]),
		}
	}
	return shares
}

func getServices(data map[string]interface{}) []asset.Service {
	if data["services"] == nil {
		return nil
	}

	u := data["services"].([]map[string]interface{})
	l := len(u)

	if l <= 0 {
		return nil
	}
	services := make([]asset.Service, l)

	for i := 0; i < l; i++ {
		svcDetail := asset.ServiceDetails{
			DisplayName:             utils.ToString(u[i]["display_name"]),
			ExecutablePath:          utils.ToString(u[i]["executable_path"]),
			StartupType:             utils.ToString(u[i]["startup_type"]),
			ServiceStatus:           utils.ToString(u[i]["service_status"]),
			LogOnAs:                 utils.ToString(u[i]["logon_as"]),
			StopEnableAction:        utils.ToBool(u[i]["stop_enable_action"]),
			DelayedAutoStart:        utils.ToBool(u[i]["delayed_auto_start"]),
			Win32ExitCode:           uint32(utils.ToInt(u[i]["win32_exit_code"])),
			ServiceSpecificExitCode: uint32(utils.ToInt(u[i]["service_specific_exit_code"])),
		}
		services[i] = asset.Service{
			Name:    utils.ToString(u[i]["service_name"]),
			Details: svcDetail,
		}
	}
	return services
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

func Select(query string, value ...interface{}) ([]map[string]interface{}, error) {
	q := Session.Query(query, value...)
	iter := q.Iter()
	data, err := iter.SliceMap()
	defer q.Release()
	if err != nil {
		return nil, exception.New("ErrDbUnableToFetchRecord", err)
	}
	err = iter.Close()
	if err != nil {
		return data, exception.New("ErrDbUnableToFetchRecord", err)
	}
	return data, nil
}

const logFlags int = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile
