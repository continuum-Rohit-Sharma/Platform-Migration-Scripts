package main

import (
	"github.com/ContinuumLLC/platform-common-lib/src/downloader"
	"github.com/ContinuumLLC/platform-common-lib/src/downloader/http"
	"github.com/ContinuumLLC/platform-common-lib/src/webClient"
)

func main() {
	service := http.GetDownloader(webClient.TlsClient, webClient.ClientConfig{
		MaxIdleConns:                100,
		MaxIdleConnsPerHost:         10,
		IdleConnTimeoutMinute:       1,
		TimeoutMinute:               1,
		DialTimeoutSecond:           100,
		DialKeepAliveSecond:         100,
		TLSHandshakeTimeoutSecond:   100,
		ExpectContinueTimeoutSecond: 100,
	})

	service.Download(&downloader.Config{
		URL:              "https://integration.agent.exec.itsupport247.net/agent/v1/endpoint/0c520f40-4cc9-406f-ac6a-317aa6eb0da2/manifest",
		DownloadLocation: "/home/lokesh/Desktop",
		FileName:         "manifest.json",
		TransactionID:    "1",
	})
}
