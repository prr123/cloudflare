// listAccounts.go
// Author: prr, azul software
// Date 3 April 2023
// copyright prr, azul software
//
// usage listAccounts
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

	dbg:= true
    numArgs := len(os.Args)
	useStr := "usage: listAccounts [/yaml=apifile]\n"

	switch numArgs {
    case 1:
//		fmt.Printf(useStr)
//        log.Fatalf("insufficient CLI args!\n")
/*
    case 2:
		argByte := []byte(os.Args(1))
		if argByte[0] == '/' {
			log.Fatalf("no domain name provided!")
		}
		domainStr = os.Args[1]
*/
	case 2:
		cmdStr := os.Args[1]
		argByte := []byte(cmdStr)
		if argByte[0] != '/' {
			if cmdStr == "help" {
				fmt.Printf(useStr)
				os.Exit(-1)
			}
			fmt.Printf(useStr)
			log.Fatalf("invalid command!")
		}
	default:
        fmt.Printf(useStr)
        log.Fatalf("too many CLI args!\n")
    }

//	domain := os.Args[1]
    cfApiFilNam := "cloudflareApi.yaml"

    if numArgs ==2 {

		flags := []string{"yaml"}
		flagMap, err := util.ParseFlagsStart(os.Args, flags, 1)
		if err != nil {
			log.Fatalf("error parseFlags: %v\n",err)
    	}

		if dbg {
			fmt.Printf("** flagMap: %d **\n", len(flagMap))
			for k,v := range flagMap {
				fmt.Printf("k: %v val: %v\n", k, v)
			}
		}

		val, ok := flagMap["yaml"]
		if !ok {
			log.Fatalf("no yaml file provided as second cli argument!")
		}
		yamlFilNamStr, ok2 := val.(string)
		if !ok2 {
			log.Fatalf("yaml file value not a string!")
		}
		cfApiFilNam = yamlFilNamStr
	}

    log.Printf("Using yaml apifile:    %s\n", cfApiFilNam)

    apiObj, err := cfLib.InitCfLib(cfApiFilNam)
    if err != nil {
        log.Fatalf("cfLib.InitCfLib: %v\n", err)
    }
    // print results
    if dbg {cfLib.PrintApiObj (apiObj)}

	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	fmt.Println("********************************************")

	// first we need to retrieve account

//	par := cloudflare.AccountsListParams{
//		Name: apiObj.Email,
//	}


	acnt, _, err := api.Account(ctx, apiObj.AccountId)
	if err != nil {
		log.Fatalf("api.Account: %v\n", err)
	}

//	fmt.Printf("account: %v\n", acnt)

	cfLib.PrintAccount(&acnt)

}
