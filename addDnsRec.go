// addDnsRec.go
// test program to delete a dns record from a domain/ zone
// Author: prr, azulsoftware
// Date: 31 March 2023
// copyright 2023 prr, azul software
//

package main

import (
//	"context"
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

	useStr := "addDnsRec domain [/yaml=file]"

    cfDir := os.Getenv("Cloudflare")
    if len(cfDir) == 0 {log.Fatalf("could not resolve Cloudflare\n")}

    cfApiFilNam := cfDir + "/token/cfDns.yaml"

    zoneDir := os.Getenv("zoneDir")
    if len(zoneDir) == 0 {log.Fatalf("could not resolve env var zoneDir!")}

    DomainFilNam := zoneDir + "/cfDomainsShort.yaml"

    if numArgs < 2 {
        fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 3 {
        fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("too many CLI args!\n")
    }

	domain := os.Args[1]

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
		cfApiFilNam, ok = val.(string)
		if !ok {
			log.Fatalf("no yaml file value not a string!")
		}
	}

    log.Printf("Using yaml file: %s\n", cfApiFilNam)

    fmt.Println("********************************************")

    //get Domain id file
    zoneList, err := cfLib.ReadZoneShortFile(DomainFilNam)
    if err != nil {log.Fatalf("ReadZoneFileShort: %v\n", err)}

    log.Printf("success reading all cf zones!\n")
	cfLib.PrintZoneList(zoneList)

    numZones := len(zoneList.Zones)

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

    apiObj, err := cfLib.InitCfApi(cfApiFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfLib: %v\n", err)
    }
    // print results
    cfLib.PrintApiObj (apiObj.ApiObj)


	fmt.Println("************** before *********************")

	// try to create DNS Record
	dnsPar := cloudflare.CreateDNSRecordParams{
		CreatedOn: time.Now(),
		Type: "TXT",
		Name: "azultest",
		Content: "abacadabra",
		TTL: 30000,
		Comment: "test for acme",
	}

	dnsRec, err := apiObj.AddDnsRecord(zone.Id, &dnsPar)
    if err != nil {
        log.Fatalf("api.CreateDNSRecord: %v\n", err)
    }

	cfLib.PrintDnsRec(dnsRec)

}

