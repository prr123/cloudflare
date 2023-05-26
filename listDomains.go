// listDomains.go
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

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr := "usage: listDomains [domainfile] [/save=json/yaml] [/api=apifile]"

	if numArgs > 4 {
		fmt.Println(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]

	cfDir := os.Getenv("CloudFlare")
	if len(cfDir) == 0 {
		log.Fatalf("could not get env: CloudFlare\n")
	}

    yamlApiFilNam := cfDir + "/token/cfZonesApi.yaml"
	DomainFilNam := "cfDomainsLong"

	flags := []string{"api","save"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("error parseFlags: %v\n",err)
    }

	numFlags := len(flagMap)

	if numArgs > numFlags + 2 {
		fmt.Println(useStr)
		log.Fatalf("error more than one cmd: %v\n",err)
	}

	if numArgs == numFlags +2 {
		DomainFilNam = os.Args[1]
		if os.Args[1] == "help" {
			fmt.Println(useStr)
			os.Exit(-1)
		}
	}

	domainExt := ".yaml"
	jsonTyp := false
	if numFlags >0 {
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

			switch saveStr {
			case "yaml":
				domainExt = ".yaml"
			case "json":
				domainExt = ".json"
				jsonTyp =true
			default:
				log.Fatalf("invalid save flag:!", saveStr)
			}
		}
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
        	log.Fatal("could not remove file %s: %v", DomainFilNam, e)
    	}
	}

	DomainFil, err := os.Create(DomainFilNam)
	if err != nil {
        log.Fatal("could not create file %s: %v", DomainFilNam, err)
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

/*
	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}
*/
	// Most API calls require a Context
	ctx := context.Background()

	fmt.Println("********************************************")

	zones, err := api.ListZones(ctx)
    if err != nil {
        log.Fatalf("api.ListZones: %v\n", err)
    }

	cfLib.PrintZones(zones)

	if jsonTyp {
		err = cfLib.SaveZonesJson(zones, DomainFil)
    	if err != nil {
        	log.Fatalf("cfLib.SaveZonesJson: %v\n", err)
    	}
	} else {
		err = cfLib.SaveZonesYaml(zones, DomainFil)
    	if err != nil {
        	log.Fatalf("cfLib.SaveZonesYaml: %v\n", err)
    	}
	}
}
