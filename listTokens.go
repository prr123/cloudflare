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
	"time"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

	dbg:= true
    numArgs := len(os.Args)

	cfDir := os.Getenv("Cloudflare")
    if len(cfDir) == 0 {log.Fatalf("could not get env: CloudFlare\n")}

    cfApiFilNam := cfDir + "/token/cfTok.yaml"

	useStr := "usage: listTokens [/yaml=apifile]\n"

	switch numArgs {
    case 1:

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

    if numArgs > 1 {

		flags := []string{"token", "dbg"}
		flagMap, err := util.ParseFlags(os.Args, flags)
		if err != nil {
			log.Fatalf("no flags found!: %v\n",err)
    	}

		val, ok := flagMap["token"]
		if !ok {
			log.Fatalf("no token value provided!")
		}
		tokFilnam, ok2 := val.(string)
		if !ok2 {
			log.Fatalf("token file value not a string!")
		}
		cfApiFilNam = cfDir +"/token/" + tokFilnam

		_, ok = flagMap["dbg"]
		if ok {
			dbg = true
		}


	}

    log.Printf("Using token apifile:    %s\n", cfApiFilNam)
	log.Printf("debug: %t\n", dbg)

    apiObj, err := cfLib.InitCfApi(cfApiFilNam)
    if err != nil {log.Fatalf("cfLib.InitCfApi: %v\n", err)}

    // print results
    if dbg {cfLib.PrintApiObj (apiObj.ApiObj)}

	// Most API calls require a Context
	ctx := context.Background()
	api := apiObj.API
//	apiobj := apiObj.ApiObj

	fmt.Println("********************************************")

	// first we need to retrieve account

	tokList, err := api.APITokens(ctx)
	if err != nil {log.Fatalf("APITokens: %v\n", err)}

	fmt.Printf("************ Token List [%d] **************\n", len(tokList))
	fmt.Printf("   seq     ID        Name        Value      Status  Exp \n")
	for i:=0; i<len(tokList); i++ {
		tok := tokList[i]
		fmt.Printf("  [%d]: %-20s| %-30s| %-5s| %-10s %-20s\n",i+1, tok.ID, tok.Name, tok.Value, tok.Status, tok.ExpiresOn.Format(time.RFC1123) )
	}

}