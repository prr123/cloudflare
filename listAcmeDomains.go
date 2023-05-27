// listAcmeDomains.go
// program that lists all domains in cloudflare and save the domain name in a yaml file
// Author: prr azul software
// Date 10 April 2023
// copyright prr azul software
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr := "usage: listAcmeDomains [domainfile] [/api=apifile] [/save]"

	if numArgs > 3 {
		fmt.Println(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    yamlApiFilNam := "/home/peter/yaml/cloudflareApi.yaml"
	AcmeDomainFilNam := "cfDomainsAcme"

	flags := []string{"api", "save"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("error parseFlags: %v\n",err)
    }

	numFlags := len(flagMap)
	log.Printf("flags: %d\n", numFlags)
	fmt.Printf("flagMap: %v\n", flagMap)

	if numArgs > numFlags + 2 {
		fmt.Println(useStr)
		log.Fatalf("error more than one cmd: %v\n",err)
	}

	if numArgs == numFlags +2 {
		if os.Args[1] == "help" {
			fmt.Println(useStr)
			os.Exit(-1)
		}
		AcmeDomainFilNam = os.Args[1]
	}

	domainExt := ".yaml"
//	jsonTyp := false
	saveFlag := false
	if numFlags >0 {
		val, ok := flagMap["api"]
		if ok {
			yamlFilNamStr, ok2 := val.(string)
			if !ok2 {
				log.Fatalf("api flag value is not a string!")
			}
			yamlApiFilNam = yamlFilNamStr
		}
		val, ok = flagMap["save"]
		if ok {
			saveFlag = true
		}
	}

	AcmeDomainFilNam = AcmeDomainFilNam + domainExt
    log.Printf("Using yaml apifile:    %s\n", yamlApiFilNam)
	if saveFlag {
	    log.Printf("Using yaml domainfile: %s\n", AcmeDomainFilNam)
	}
    log.Printf("Using saveFlag: %t\n", saveFlag)


	// create yamlDomainFile
	var DomainFil *os.File
	if saveFlag {
		if _, err := os.Stat(AcmeDomainFilNam); err != nil {
			log.Printf("no existing domain file: %v!", err)
		} else {
			log.Printf("removing existing domain file!")
     		e := os.Remove(AcmeDomainFilNam)
    		if e != nil {
        		log.Fatal("could not remove file %s: %v", AcmeDomainFilNam, e)
    		}
		}

		DomainFil, err = os.Create(AcmeDomainFilNam)
		if err != nil {
        	log.Fatal("could not create file %s: %v", AcmeDomainFilNam, err)
		}
		defer DomainFil.Close()
	}

    apiObj, err := cfLib.InitCfLib(yamlApiFilNam)
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

	// need to find Acme Records for each zone /domain
    var rc cloudflare.ResourceContainer
    rc.Level = cloudflare.ZoneRouteLevel

	acmeZones := make([]cfLib.ZoneAcme, len(zones))

    var listDns cloudflare.ListDNSRecordsParams
	count :=0
    for i:=0; i<len(zones); i++ {
//        fmt.Printf("Zone[%d]: Name: %s Id: %s\n", i+1, zone.Name, zone.Id)

    	rc.Identifier = zones[i].ID
    	dnsRecs, _, err := api.ListDNSRecords(ctx, &rc, listDns)
    	if err != nil {
        	log.Fatalf("domain[%d]: %s api.ListDNSRecords: %v\n", i+1, zones[i].Name, err)
    	}

		dnsId := ""
		for j:=0; j< len(dnsRecs); j++ {
			idx := strings.Index(dnsRecs[j].Name, "_acme-challenge.")
// dbg: fmt.Printf("rec [%d] name: %s idx %d\n", j, dnsRecs[j].Name, idx)
			if idx == 0 {
				dnsId = dnsRecs[j].ID
				break
			}
		}
		if len(dnsId) > 0 {
			acmeZones[count].Name = zones[i].Name
			acmeZones[count].Id = zones[i].ID
			acmeZones[count].AcmeId = dnsId
			count++
		}
	}

	log.Printf("found %d Domains containing an Acme challange record!", count)

	cfLib.PrintAcmeZones(acmeZones[:count])

	if saveFlag {
		err = cfLib.SaveAcmeDns(acmeZones[:count], DomainFil)
    	if err != nil {
       		log.Fatalf("cfLib.SaveZonesYaml: %v\n", err)
    	}
		log.Printf("success listAcmeDomains created Acme Domain File")
	} else {
		log.Printf("success listAcmeDomains no Acme Domain File creation")
	}
}
