// listUserInfo.go
// Author: prr, azul software
// Date: 29 March 2023
// copyright 2023 prr, azul software
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
    if numArgs > 2 {
        fmt.Printf("CLI args are not equal to 2: %d\n", numArgs)
        fmt.Printf("usage: readCFYaml yaml file\n")
        os.Exit(-1)
    }

    cfDir := os.Getenv("CloudFlare")
    if len(cfDir) == 0 {
        log.Fatalf("could not get env: CloudFlare\n")
    }

    yamlApiFilNam := cfDir + "/token/cfUser.yaml"

    if numArgs == 2 {yamlApiFilNam = os.Args[1]}

    log.Printf("Using yaml file: %s\n", yamlApiFilNam)

    apiObj, err := cfLib.InitCfApi(yamlApiFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfLib: %v\n", err)
    }
    // print results
    cfLib.PrintApiObj (apiObj.ApiObj)

/*
//	cloudToken := "O5ART89fgxulItZ1l-o9PScX-uEGXN219dzo06Xi"
	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}
*/
	// Most API calls require a Context
	ctx := context.Background()
	api := apiObj.API

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatalf("api.UserDetails: %v\n", err)
	}

	// Print user details
//	fmt.Println(u)
	fmt.Println("********************************************")

	cfLib.PrintUserInfo(&u)
}

