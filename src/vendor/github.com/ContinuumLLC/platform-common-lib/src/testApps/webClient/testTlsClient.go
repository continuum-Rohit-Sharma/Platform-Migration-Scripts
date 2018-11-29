package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ContinuumLLC/platform-common-lib/src/webClient"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://internal-continuum-agent-service-elb-int-1915575479.us-east-1.elb.amazonaws.com/agent/version", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err = webClient.ClientFactoryImpl{}.GetClientServiceByType(webClient.TlsClient).Do(req)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(res.Status)
		io.Copy(os.Stdout, res.Body)
	}
}
