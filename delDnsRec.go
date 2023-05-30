// delDnsRec.go
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
	"bufio"
	"strconv"
//	"time"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

    zoneDir := os.Getenv("zoneDir")
    if len(zoneDir) == 0 {log.Fatalf("could not resolve env var zoneDir!")}

    DomainFilNam := zoneDir + "/cfDomainsShort.yaml"

    cfDir := os.Getenv("Cloudflare")
    if len(cfDir) == 0 {log.Fatalf("could not resolve env var Cloudflare!")}

    cfApiFilNam := cfDir + "/token/cfDns.yaml"

	useStr := "delDnsRec domain rec [/yaml=file]"
    if numArgs < 1 {
        fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 4 {
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

    log.Printf("Using api token: %s\n", cfApiFilNam)

    apiObj, err := cfLib.InitCfApi(cfApiFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfApi: %v\n", err)
    }
    // print results
    cfLib.PrintApiObj (apiObj.ApiObj)

    fmt.Println("********************************************")

    //get Domain id file
    zoneList, err := cfLib.ReadZoneShortFile(DomainFilNam)
    if err != nil {log.Fatalf("ReadZoneFileShort: %v\n", err)}

    log.Printf("success reading all cf zones!\n")
    cfLib.PrintZoneList(zoneList)

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


    dnsRecList,err := apiObj.ListDnsRecords(zone.Id)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }
    cfLib.PrintDnsRecs(dnsRecList)

    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("Enter Dns Rec number: ")
    str, err := reader.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    numStr := str[:len(str)-1]
    num, err := strconv.Atoi(numStr)
    if err!= nil {log.Fatalf("string conversion: %v\n", err)}
    fmt.Printf("DnsRec: %d\n", num)

	if num<1 || num > len(*dnsRecList) {log.Fatal("dns Rec number outside boundary!")}

	dnsRec := (*dnsRecList)[num-1]
	cfLib.PrintDnsRec(&dnsRec)

//	dnsRecId := dnsRec.ID
	fmt.Println("************** delDnsRec *********************")

	err = apiObj.DelDnsRec(zone.Id, dnsRec.ID)
    if err != nil {
        log.Fatalf("api.DelDNSRec: %v\n", err)
    }

	fmt.Println("Success Delete Dns Record")

//	fmt.Println("************** after deletion *********************")
//	cfLib.PrintDnsRec(&dnsRecs)

}

