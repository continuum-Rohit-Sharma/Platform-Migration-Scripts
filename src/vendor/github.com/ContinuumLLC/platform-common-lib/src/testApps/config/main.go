package main

import (
	"fmt"
	"os"

	"github.com/ContinuumLLC/platform-common-lib/src/config"
)

func main() {
	cfg := config.Configuration{
		FilePath:      os.Args[1],
		Content:       "{\"tree1\\\\b\" : {\"tree1.prop2\" : \"prop2\"}}",
		TransationID:  "111",
		PartialUpdate: true,
	}

	srv := config.GetConfigurationService()
	u, err := srv.Update(cfg)
	fmt.Printf("Updated Config   : %+v\nProcessing Error : %+v\n", u, err)
}
