// program that reads a yaml file that contains the zones and ids for a clouflare account
// Author: prr, azulsoftware
// Date: 10 April 2023
// copyright prr, azulsoftware
//

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
	useStr := "usage: readDomainsShort [domainfile] [/api=apifile]"

	if numArgs > 3 {
		fmt.Println(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//	cfDir := os.Getenv("Cloudflare")
//	if len(cfDir) == 0 {log.Fatalf("could not resolve env var Cloudflare!")}

	zoneDir := os.Getenv("zoneDir")
	if len(zoneDir) == 0 {log.Fatalf("could not resolve env var zoneDir!")}

//	domain := os.Args[1]
//    cfApiFilNam := cfDir + "/token/Api.yaml"
	DomainFilNam := zoneDir + "/cfDomainsShort.yaml"

	flags := []string{"api"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("error parseFlags: %v\n",err)
    }

	numFlags := len(flagMap)

	if numArgs > numFlags + 1 {
		fmt.Println(useStr)
		log.Fatalf("error more than one flag: %v\n",err)
	}

	if numArgs == numFlags +2 {
		DomainFilNam = os.Args[1]
		if os.Args[1] == "help" {
			fmt.Println(useStr)
			os.Exit(-1)
		}
	}

	domainExt := ".yaml"
//	jsonTyp := false
	if numFlags >0 {
/*
		val, ok := flagMap["api"]
		if ok {
			yamlFilNamStr, ok2 := val.(string)
			if !ok2 {
				log.Fatalf("api flag value is not a string!")
			}
//			yamlApiFilNam = yamlFilNamStr
		}
*/
		saveVal, ok := flagMap["filTyp"]
		if ok {
			saveStr, ok2 := saveVal.(string)
			if !ok2 {
				log.Fatalf("filTyp flag value is not a string!")
			}

			switch saveStr {
			case "yaml":
				domainExt = ".yaml"
			case "json":
				domainExt = ".json"
//				jsonTyp =true
			default:
				log.Fatalf("invalid save flag:!", saveStr)
			}
		}
	}

//	DomainFilNam = DomainFilNam + domainExt
//    log.Printf("Using yaml apifile:    %s\n", yamlApiFilNam)
    log.Printf("Using domainsfile: %s\n", DomainFilNam)
    log.Printf("domainsfile type: %s\n", domainExt)

	// create yamlDomainFile
/*
	if _, err := os.Stat(DomainFilNam); err != nil {
		log.Fatalf("no existing domain file: %v!\n", err)
	}
*/

/*
	infil, err := os.Open(DomainFilNam)
	if err != nil {
        log.Fatal("could not open file %s: %v", DomainFilNam, err)
	}

	if jsonTyp {
		log.Fatal("json read: still todo\n")

		err = cfLib.SaveZonesShortJson(zoneShortList, DomainFil)
    	if err != nil {
        	log.Fatalf("cfLib.SaveZonesShortJson: %v\n", err)
    	}
	}
*/

	zoneList, err := cfLib.ReadZoneShortFile(DomainFilNam)
	if err != nil {log.Fatalf("cfLib.ReadZonesShortFile: %v\n", err)}


	cfLib.PrintZoneList(zoneList)
}
