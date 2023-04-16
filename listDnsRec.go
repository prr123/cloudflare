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
	useStr :=  "usage: listDnsRec domain [/yaml=file]\n"
    DomainFilNam := "cfDomainsShort.yaml"

    if numArgs < 2 {
        fmt.Printf(useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 3 {
        fmt.Printf(useStr)
        log.Fatalf("too many CLI args!\n")
    }

    yamlFilNam := "cloudflareApi.yaml"
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
		yamlFilNam, ok = val.(string)
		if !ok {
			log.Fatalf("no yaml file value not a string!")
		}
	}

	log.Printf("Lookup DNS Records for domain: %s\n", domain)
    log.Printf("Using yaml file: %s\n", yamlFilNam)
    log.Printf("Using yaml domainsfile: %s\n", DomainFilNam)

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

	//get Domain id file
    if _, err := os.Stat(DomainFilNam); err != nil {
        log.Fatalf("no existing domain file: %v!\n", err)
    }

    infil, err := os.Open(DomainFilNam)
    if err != nil {
        log.Fatal("could not open file %s: %v", DomainFilNam, err)
    }

    var zoneShortList *[]cfLib.ZoneShort

	zoneShortList, err = cfLib.ReadZonesShortYaml(infil)
	if err != nil {
		log.Fatalf("cfLib.ReadZonesShortYaml: %v\n", err)
	}

	found := -1
    for i:=0; i<len((*zoneShortList)); i++ {
        zone := (*zoneShortList)[i]
		if zone.Name == domain {
			found = i
			break
		}
    }
	if found < 0 {
		log.Fatalf("domain: %s not found!\n", domain)
	}

	zone := (*zoneShortList)[found]
	fmt.Printf("Zone[%d]: Name: %s Id: %s\n", found+1, zone.Name, zone.Id)

	// try to list DNS Parameters
	var listDns cloudflare.ListDNSRecordsParams

	var rc cloudflare.ResourceContainer
	rc.Level = cloudflare.ZoneRouteLevel
//	rc.Identifier = "0e6e30d5edb4c1025817eb1678511cef"
	rc.Identifier = zone.Id

//	dnsRecs, resInfo, err := api.ListDNSRecords(ctx, &rc, listDns)
	dnsRecs, _, err := api.ListDNSRecords(ctx, &rc, listDns)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }
//	fmt.Printf("resInfo: %v\n", resInfo)
//	fmt.Printf("Dns Records [%d]\n",len(dnsRecs))

//	cfLib.PrintResInfo(resInfo)
	cfLib.PrintDnsRecs(&dnsRecs)
}

