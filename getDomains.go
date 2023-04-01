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
        fmt.Printf("usage: getDomains [domainfile] [/yaml=apifile]\n")
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 3 {
        fmt.Printf("usage: getDomains [domainfile] [/yaml=apifile]\n")
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    yamlFilNam := "cloudflareApi.yaml"
	yamlDomainFilNam := "cfDomains.yaml"

    if numArgs > 1 {

		flags := []string{"yaml"}
		flagMap, err := util.ParseFlags(os.Args, flags)
		if err != nil {
			log.Fatalf("error parseFlags: %v\n",err)
    	}

		val, ok := flagMap["yaml"]
		if ok {
			yamlFilNamStr, ok2 := val.(string)
			if !ok2 {
				log.Fatalf("no yaml file value not a string!")
			} else {
				yamlFilNam = yamlFilNamStr
			}
			if numArgs == 3 {
				yamlDomainFilNam = os.Args[1]
			}
		} else {
			yamlDomainFilNam = os.Args[1]
		}
	}

    log.Printf("Using yaml apifile:    %s\n", yamlFilNam)
    log.Printf("Using yaml domainfile: %s\n", yamlDomainFilNam)

	// create yamlDomainFile
	if _, err := os.Stat(yamlDomainFilNam); err != nil {
		log.Printf("no existing domain file: %v!", err)
	} else {
		log.Printf("removing existing domain file!")
     	e := os.Remove(yamlDomainFilNam)
    	if e != nil {
        	log.Fatal("could not remove file %s: %v", yamlDomainFilNam, e)
    	}
	}

	yamlDomainFil, err := os.Create(yamlDomainFilNam)
	if err != nil {
        log.Fatal("could not create file %s: %v", yamlDomainFilNam, err)
	}
	defer yamlDomainFil.Close()

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

	cfLib.PrintZones(zones)

	err = cfLib.SaveZones(zones, yamlDomainFil)
    if err != nil {
        log.Fatalf("cfLib.SaveZones: %v\n", err)
    }

}

