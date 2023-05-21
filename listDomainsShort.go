// creDomainListShort
// program that creates a yaml file of all zones/ domains in an account
// Author: prr azul software
// Date: 16 April 2023
// copyright 2023 prr, azul software
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"path/filepath"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr := "usage: creDomainListShort [/save=domainfile] [/api=apifile]"


	if numArgs > 3 {
		fmt.Println(useStr)
        log.Fatalf("too many CLI args!\n")
    }

    cfDir := os.Getenv("CloudFlare")
    if len(cfDir) == 0 {
        log.Fatalf("could not get env: CloudFlare\n")
    }

    yamlApiFilNam := cfDir + "/token/cfZonesApi.yaml"

    zoneDir := os.Getenv("zoneDir")
    if len(zoneDir) == 0 {
        log.Fatalf("could not get env: zoneDir\n")
    }

	DomainFilNam := zoneDir + "/cfDomainsShort"

	flags := []string{"api", "save"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("error parseFlags: %v\n",err)
    }

	numFlags := len(flagMap)

	if numArgs > numFlags + 2 {
		fmt.Println(useStr)
		log.Fatalf("error unknown commands: %v\n", err)
	}

	if numArgs == numFlags + 2 {
		DomainFilNam = os.Args[1]
		if os.Args[1] == "help" {
			fmt.Println(useStr)
			os.Exit(-1)
		} else {
			fmt.Println(useStr)
			log.Fatalf("error unknown command: %s : %v\n", os.Args[1], err)
		}
	}

	domainExt := ".yaml"

	val, ok := flagMap["api"]
	if ok {
		yamlFilNamStr, ok2 := val.(string)
		if !ok2 {
			log.Fatalf("api flag value is not a string!")
		}
		yamlApiFilNam = yamlFilNamStr
	}

	saveVal, ok := flagMap["save"]
	if ok {
		saveStr, ok2 := saveVal.(string)
		if !ok2 {
			log.Fatalf("save flag value is not a string!")
		}
		DomainFilNam = saveStr
	}

	DomainFilNam = DomainFilNam + domainExt
    log.Printf("Using yaml apifile:    %s\n", yamlApiFilNam)
    log.Printf("Using yaml domainfile: %s\n", DomainFilNam)


	// create yamlDomainFile
	if _, err := os.Stat(DomainFilNam); err != nil {
		log.Printf("no existing domain file: %v!", err)
	} else {
		log.Printf("removing existing domain file!")
     	e := os.Remove(DomainFilNam)
    	if e != nil {
        	log.Fatalf("could not remove file %s: %v", DomainFilNam, e)
    	}
	}

	// make sure zones folder exists
	// todo
	//
	path, err := filepath.Abs(DomainFilNam)
	if err != nil {
		log.Fatalf("filepath error: %v", err)
	}

	_, err = os.Stat(filepath.Dir(path))
	if err != nil {
		log.Fatalf("folder %s err: %v: %v", filepath.Dir(path), err)
	}

	DomainFil, err := os.Create(DomainFilNam)
	if err != nil {
        log.Fatalf("could not create file %s: %v", DomainFilNam, err)
	}
	defer DomainFil.Close()

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

	cfLib.PrintZones(zones)


	zoneShortList := make([]cfLib.ZoneShort, len(zones))


	for i:=0; i< len(zones); i++ {
		zoneShortList[i].Name =zones[i].Name
		zoneShortList[i].Id =zones[i].ID
	}

	zoneList := &cfLib.ZoneList{
			AccountId: apiObj.AccountId,
			Email: apiObj.Email,
			ModTime: time.Now(),
			Zones: zoneShortList,
		}
	zoneList.Modified = zoneList.ModTime.Format(time.RFC1123)


	cfLib.PrintZoneList(zoneList)

//	os.Exit(1)

	err = cfLib.SaveZoneShortFile(zoneList, DomainFil)
    if err != nil {
        log.Fatalf("cfLib.SaveZoneShortFile: %v\n", err)
    }

}
