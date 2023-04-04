// addDnsRec.go
// test program to delete a dns record from a domain/ zone
// Author: prr, azulsoftware
// Date: 31 March 2023
// copyright 2023 prr, azul software
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
    if numArgs < 1 {
        fmt.Printf("usage: addDnsRec [/yaml=file]\n")
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 2 {
        fmt.Printf("usage: addDnsRec [/yaml=file]\n")
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    yamlFilNam := "cloudflareApi.yaml"

    if numArgs == 2 {

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

	fmt.Println("************** before *********************")

//	cfLib.PrintDnsRec(&dnsRecs)

	// try to create DNS Record
	dnsPar := cloudflare.CreateDNSRecordParams{
		CreatedOn: time.Now(),
		Type: "TXT",
		Name: "azulacademy.eu",
		Content: "abacadabra",
		TTL: 30000,
		Comment: "test for acme",
	}

	var rc cloudflare.ResourceContainer
	//domains
	rc.Level = cloudflare.ZoneRouteLevel
	//domain id
	rc.Identifier = "0e6e30d5edb4c1025817eb1678511cef"


	dnsRec, err := api.CreateDNSRecord(ctx, &rc, dnsPar)
    if err != nil {
        log.Fatalf("api.CreateDNSRecord: %v\n", err)
    }

	fmt.Printf("Dns Record %v\n",dnsRec)

//	cfLib.PrintResInfo(resInfo)
//	fmt.Println("************** after *********************")
//	cfLib.PrintDnsRec(&dnsRecs)

}

