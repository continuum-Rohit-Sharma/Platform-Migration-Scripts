package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ContinuumLLC/platform-common-lib/src/checksum"
	"github.com/ContinuumLLC/platform-common-lib/src/downloader"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/webClient"
)

const (
	fileName string = "checksumfile"
)

type serviceImpl struct {
	client webClient.ClientService
	log    logging.Logger
}

//GetDownloader is a definition returns a HTTP downloader instance
func GetDownloader(clientType webClient.ClientType, config webClient.ClientConfig) downloader.Service {
	return serviceImpl{
		client: webClient.ClientFactoryImpl{}.GetClientServiceByType(clientType, config),
		log:    logging.GetLoggerFactory().Get(),
	}
}

func (s serviceImpl) Download(conf *downloader.Config) error {
	err := s.downloadFile(conf)
	if err != nil {
		return fmt.Errorf("Failed to download file for url : %s with Error : %v", conf.URL, err)
	}
	if conf.CheckSumType == "" {
		return nil
	}
	mode := checksum.GetType(conf.CheckSumType)
	service, err := checksum.GetService(mode)
	if err != nil {
		return fmt.Errorf("Error Occurred because No Validator for CheckSum Type : %s is defined Err : %+v", conf.CheckSumType, err)
	}
	reader, verifyCS, err := s.getChecksum(conf)
	if err != nil {
		return fmt.Errorf("Error Occurred while Getting the Checksum, Err : %v", err)
	}

	defer reader.Close()

	_, err = service.Validate(reader, verifyCS)
	if err != nil {
		return fmt.Errorf("Checksum cannot be verified or Failed to download the checksum with Error : %v", err)
	}
	return nil
}

func (s serviceImpl) getChecksum(conf *downloader.Config) (io.ReadCloser, string, error) {
	dFLocation := filepath.Join(conf.DownloadLocation, conf.FileName)
	var verifyCS string
	var err error
	if conf.CheckSum == "" {
		verifyCS, err = s.getCheckSumValFromFile(conf)
		if err != nil {
			return nil, "", fmt.Errorf("Error Occurred while downloading the checksum, Err : %+v", err)
		}
	} else {
		verifyCS = conf.CheckSum
	}

	reader, err := os.Open(dFLocation)
	if err != nil {
		return nil, "", fmt.Errorf("Filed to open a file %s : %+v", dFLocation, err)
	}
	return reader, verifyCS, nil
}

func (s serviceImpl) getCheckSumValFromFile(conf *downloader.Config) (string, error) {
	newConf := downloader.Config{
		CheckSum:         conf.CheckSum,
		DownloadLocation: conf.DownloadLocation,
		FileName:         fileName,
		KeepOriginalName: conf.KeepOriginalName,
		TransactionID:    conf.TransactionID,
		URL:              fmt.Sprintf("%s.%s", conf.URL, conf.CheckSumType),
		CheckSumType:     conf.CheckSumType,
	}
	err := s.downloadFile(&newConf)
	if err != nil {
		return "", fmt.Errorf("Failed to download Checksum file for url : %s with Error : %v", conf.URL, err)
	}
	cSFileLoc := filepath.Join(newConf.DownloadLocation, newConf.FileName)
	return s.readCheckSumFromFile(cSFileLoc)
}

func (s serviceImpl) readCheckSumFromFile(path string) (string, error) {
	md5Data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(md5Data[:]), nil
}

func (s serviceImpl) downloadFile(conf *downloader.Config) error {
	req, err := http.NewRequest(http.MethodGet, conf.URL, nil)
	if err != nil {
		s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "Failed to create request for url : %s with Error : %v", conf.URL, err)
		return err
	}

	for key, value := range conf.Header {
		req.Header.Set(key, value)
	}

	res, err := s.client.Do(req)
	if err != nil {
		s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "Failed to execute request for url : %s with Error : %v", conf.URL, err)
		return err
	}
	err = os.MkdirAll(conf.DownloadLocation, os.ModePerm)
	if err != nil {
		return err
	}
	generateFileName(conf, res)
	s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "Recieved Response with status %s for url : %s", res.Status, conf.URL)
	body := res.Body
	defer body.Close()
	return s.createFile(conf, res.Body)
}

//### Should be moved to common file at the time of adding new Downloader like FTP ####
func (s serviceImpl) createFile(conf *downloader.Config, data io.ReadCloser) error {
	dst := conf.DownloadLocation + string(os.PathSeparator) + conf.FileName
	out, err := os.Create(dst)
	if err != nil {
		s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "Failed to create File : %s with Error : %v", dst, err)
		return err
	}
	defer out.Close()

	w, err := io.Copy(out, data)
	if err != nil {
		s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "Failed to copy File : %s with Error %v", dst, err)
		return err
	}
	s.log.LogWithTransactionf(logging.TRACE, conf.TransactionID, "%d of bytes copied for File : %s", w, dst)
	return nil
}

func generateFileName(config *downloader.Config, resp *http.Response) {
	if config.FileName != "" {
		return
	}

	header := resp.Header.Get("content-disposition")
	filename := resp.Request.URL.String()
	if header != "" {
		_, params, err := mime.ParseMediaType(header)
		if err == nil && params["filename"] != "" {
			filename = params["filename"]
		}
	}
	if !config.KeepOriginalName {
		h := md5.New()
		h.Write([]byte(resp.Request.URL.String()))
		config.FileName = hex.EncodeToString(h.Sum(nil)) + filepath.Ext(filename)
		return
	}
	config.FileName = filepath.Base(filename)
}
