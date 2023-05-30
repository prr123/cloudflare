package main

import (
//	"context"
	"fmt"
	"log"
	"os"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr :=  "usage: listDnsRec domain [/yaml=file]\n"

	zoneDir := os.Getenv("zoneDir")
	if len(zoneDir) == 0 {log.Fatalf("could not resolve env var zoneDir!")}

    DomainFilNam := zoneDir + "/cfDomainsShort.yaml"

	cfDir := os.Getenv("Cloudflare")
	if len(cfDir) == 0 {log.Fatalf("could not resolve env var Cloudflare!")}

    yamlApiFilNam := cfDir + "/token/cfDns.yaml"

	dbg:= true

    if numArgs < 2 {
        fmt.Printf(useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 3 {
        fmt.Printf(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//    yamlFilNam := "cloudflareApi.yaml"
	domain := os.Args[1]

	if len(domain) == 0 {
        fmt.Printf(useStr)
        log.Fatalf("no doman provided!\n")
	}
	if domain == "help" {
        fmt.Printf(useStr)
		os.Exit(1)
	}

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
		yamlApiFilNam, ok = val.(string)
		if !ok {
			log.Fatalf("no yaml file value not a string!")
		}
	}

	log.Printf("Lookup DNS Records for domain: %s\n", domain)
    log.Printf("Using token file: %s\n", yamlApiFilNam)
    log.Printf("Using domain file: %s\n", DomainFilNam)

    apiObj, err := cfLib.InitCfApi(yamlApiFilNam)
    if err != nil {log.Fatalf("cfLib.InitCfLib: %v\n", err)}

    // print results
    cfLib.PrintApiObj (apiObj.ApiObj)


	// Most API calls require a Context
//	ctx := context.Background()
//	api := apiObj.API

	fmt.Println("********************************************")

	//get Domain id file
    zoneList, err := cfLib.ReadZoneShortFile(DomainFilNam)
    if err != nil {log.Fatalf("ReadZoneFileShort: %v\n", err)}

    log.Printf("success reading all cf zones!\n")
    if dbg {cfLib.PrintZoneList(zoneList)}

    numZones := len(zoneList.Zones)
    if _, err := os.Stat(DomainFilNam); err != nil {
        log.Fatalf("no existing domain file: %v!\n", err)
    }
	log.Printf("Zones: %d\n",numZones)

	found := -1
    for i:=0; i<numZones; i++ {
        zone := zoneList.Zones[i]
		if zone.Name == domain {
			found = i
			break
		}
    }
	if found < 0 {
		log.Fatalf("domain: %s not found!\n", domain)
	}

	zone := zoneList.Zones[found]
	fmt.Printf("Zone[%d]: Name: %s Id: %s\n", found+1, zone.Name, zone.Id)


	dnsRecs,err := apiObj.ListDnsRecords(zone.Id)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }
	cfLib.PrintDnsRecs(dnsRecs)
}
