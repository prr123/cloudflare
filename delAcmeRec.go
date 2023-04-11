// delAcmeRec.go
// test program that deletes an acme challenge record from a domain/ zone
// Author: prr, azulsoftware
// Date: 11 April 2023
// copyright 2023 prr, azul software
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
//	"time"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr:= "delAcmeRec domain [/acme=domain][/api=file]\n"
    if numArgs < 2 {
        fmt.Println(useStr)
        log.Fatalf("insufficient CLI args! Need to specify a domain!\n")
    }
	if numArgs > 4 {
        fmt.Println(useStr)
        log.Fatalf("too many CLI args!\n")
    }

	domain := os.Args[1]

    yamlFilNam := "cloudflareApi.yaml"
	acmeFilNam := "cfDomainsAcme.yaml"

    if numArgs > 2 {

		flags := []string{"api, acme"}
		flagMap, err := util.ParseFlags(os.Args, flags)
		if err != nil {
			log.Fatalf("error parseFlags: %v\n",err)
    	}

		val, ok := flagMap["api"]
		if ok {
			yamlFilNam, ok = val.(string)
			if !ok {
				log.Fatalf("api file value not a string!")
			}
		}
		val, ok = flagMap["acme"]
		if ok {
			acmeFilNam, ok = val.(string)
			if !ok {
				log.Fatalf("acme file value not a string!")
			}
		}
	}

	log.Printf("domain: %s\n", domain)
    log.Printf("Using acme file: %s\n", acmeFilNam)
    log.Printf("Using api file: %s\n", yamlFilNam)

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

	// reading acme file
	zoneList, err := cfLib.ReadAcmeZones(acmeFilNam)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}

	fmt.Printf("zones: %d\n", len(*zoneList))

	found := -1
	for i:=0; i< len(*zoneList); i++ {
		if (*zoneList)[i].Name == domain {
			found = i
			break
		}
	}
	if found == -1 {
		log.Fatalf("domain % not found in zoneList!/n", domain)
	}

	log.Printf("domain: %s\n", (*zoneList)[found].Name)
	log.Printf("Zone Id: %s\n", (*zoneList)[found].Id)
	log.Printf("Acme Rec Id: %s\n", (*zoneList)[found].AcmeId)
	os.Exit(0)

	fmt.Println("************** delDnsRec *********************")

	dnsRecId := "b12cd8fae338120e4aced6378fa8d5e5"

	var rc cloudflare.ResourceContainer
	//domains
	rc.Level = cloudflare.ZoneRouteLevel
	//domain id
	rc.Identifier = "0e6e30d5edb4c1025817eb1678511cef"

	err = api.DeleteDNSRecord(ctx, &rc, dnsRecId)
    if err != nil {
        log.Fatalf("api.DeleteDNSRecord: %v\n", err)
    }

	fmt.Println("Success Delete Dns Record")

//	cfLib.PrintResInfo(resInfo)
//	fmt.Println("************** after *********************")
//	cfLib.PrintDnsRec(&dnsRecs)

}

