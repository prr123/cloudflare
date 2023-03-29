package main

import (
	"context"
	"fmt"
	"log"
	"os"

//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
    if numArgs > 2 {
        fmt.Printf("CLI args are not equal to 2: %d\n", numArgs)
        fmt.Printf("usage: readCFYaml yaml file\n")
        os.Exit(-1)
    }

    yamlFilNam := "cloudflareApi.yaml"

    if numArgs == 2 {yamlFilNam = os.Args[1]}

    log.Printf("Using yaml file: %s\n", yamlFilNam)

    apiObj, err := cfLib.InitCfLib(yamlFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfLib: %v\n", err)
    }
    // print results
    cfLib.PrintApiObj (apiObj)

	// Construct a new API object using a global API key
//	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	// alternatively, you can use a scoped API token

	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	fmt.Println("********************************************")

	// try to list DNS Parameters

	var listDns cloudflare.ListDNSRecordsParams
	listDns.Name = "azulacademy.eu"

	var rc cloudflare.ResourceContainer
	rc.Level = cloudflare.ZoneRouteLevel
	rc.Identifier = "0e6e30d5edb4c1025817eb1678511cef"

	dnsRecs, resInfo, err := api.ListDNSRecords(ctx,&rc , listDns)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }
//	fmt.Printf("resInfo: %v\n", resInfo)
//	fmt.Printf("Dns Records [%d]\n",len(dnsRecs))

	cfLib.PrintResInfo(resInfo)
	cfLib.PrintDnsRec(&dnsRecs)

}

