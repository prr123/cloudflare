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

	var domainStr string
	dbg:= true
    numArgs := len(os.Args)
	useStr := "usage: addDomain domain [/yaml=apifile]\n"

	switch numArgs {
    case 1:
		fmt.Printf(useStr)
        log.Fatalf("insufficient CLI args!\n")
    case 2:
		argByte := []byte(os.Args(1))
		if argByte[0] == '/' {
			log.Fata.f("no domain name provided!")
		}
		domainStr = os.Args[1]
	case 3:
		domainStr = os.Args[1]

	default: {
        fmt.Printf(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    cfApiFilNam := "cloudflareApi.yaml"

    if numArgs ==3 {

		flags := []string{"yaml"}
		flagMap, err := util.ParseFlags(os.Args, flags)
		if err != nil {
			log.Fatalf("error parseFlags: %v\n",err)
    	}

		val, ok := flagMap["yaml"]
		if !ok {
			log.Fatalf("no yaml file provided as second cli argument!")
		}
		yamlFilNamStr, ok2 := val.(string)
		if !ok2 {
			log.Fatalf("yaml file value not a string!")
		}

	}

    log.Printf("Using yaml apifile:    %s\n", cfApiFilNam)
    log.Printf("domain: %s\n", domainStr)

/*
	// create yamlDomainFile
	yaml.DmainFilNam := 
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
*/

    apiObj, err := cfLib.InitCfLib(cfApiFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfLib: %v\n", err)
    }
    // print results
    if dbg {cfLib.PrintApiObj (apiObj)}

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

	// first we need to retrieve account
	// todo store account id in api 

	// todo: support for full or partial
	zoneTyp := "partial"
	jump := true
	

	// todo check whether domain is registered with namecheap
	zoneNam := domainStr
	zone, err :=api.CreateZone(ctx, domain, jump, act, zoneTyp)
/*
	zones, err := api.ListZones(ctx)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }

	cfLib.PrintZones(zones)

	err = cfLib.SaveZones(zones, yamlDomainFil)
    if err != nil {
        log.Fatalf("cfLib.SaveZones: %v\n", err)
    }
*/
}

