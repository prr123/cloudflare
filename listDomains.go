package main

import (
	"context"
	"fmt"
	"log"
	"os"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
    if numArgs < 1 {
        fmt.Printf("usage: listDomains [/yaml=file]\n")
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 2 {
        fmt.Printf("usage: listDomains [/yaml=file]\n")
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    yamlFilNam := "cloudflareApi.yaml"

    if numArgs == 3 {

		flags := []string{"yaml"}
		flagMap, err := util.ParseFlags(os.Args, flags)
		if err != nil {
			log.Fatalf("error parseFlags: %v\n",err)
    	}

		val, ok := flagMap["yaml"]
		if !ok {
			log.Fatalf("no yaml file specified!")
		}
		yamlFilNam, ok = val.(string)
		if !ok {
			log.Fatalf("no yaml file value not a string!")
		}
	}

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

	zones, err := api.ListZones(ctx)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }

	PrintZones(zones)
}

func PrintZones(zones []cloudflare.Zone) {

	fmt.Printf("************** Zones/Domains [%d] *************\n", len(zones))

	for i:=0; i< len(zones); i++ {
		zone := zones[i]
		fmt.Printf("%d %-20s %s\n",i+1, zone.Name, zone.ID)
	}
}
