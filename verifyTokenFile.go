// verifyTokenFile.go
// Author: prr, azul software
// Date 13 July 2023
// copyright prr, azul software
//
// usage: verifyTokenFile token-file
//

package main

import (
//	"context"
	"fmt"
	"log"
	"os"

    util "github.com/prr123/utility/utilLib"
	"ns/cloudflare/cfLib"
)

func main() {

	dbg:= false
    numArgs := len(os.Args)

	cfDir := os.Getenv("Cloudflare")
    if len(cfDir) == 0 {log.Fatalf("could not get env: CloudFlare\n")}

//    cfApiFilnam := cfDir + "/token/"

	useStr := "verifyTokenFile /tokFil=tokenfile /dbg"
	helpStr :="program retrieves yaml file from token directory and verifies token"


	if numArgs==2 && os.Args[1] == "help" {
			fmt.Printf("help:\n%s\n", helpStr)
			fmt.Printf("usage: %s\n",useStr)
			os.Exit(1)
	}

   if numArgs == 1 {
        fmt.Printf("no flags provided!")
        fmt.Printf(useStr)
        os.Exit(1)
    }

    flags := []string{"tokFil", "dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("ParseFlags: %v\n",err)
    }

    val, ok := flagMap["tokFil"]
    if !ok {
        log.Fatalf("/out flag is missing!")
    }
    tokFilnam, _ := val.(string)
    if tokFilnam == "none" {
        log.Fatalf("toke file value not provided!")
    }

    cfApiFilnam := tokFilnam

    _, ok = flagMap["dbg"]
    if ok {
        dbg = true
		for k, v := range flagMap {
			fmt.Printf("flag: /%s value: %s\n",k, v)
		}
    }

	if dbg {
    	log.Printf("Using token apifile:    %s\n", cfApiFilnam)
		log.Printf("debug: %t\n", dbg)
	}

    apiObj, err := cfLib.InitCfApi(cfApiFilnam)
    if err != nil {log.Fatalf("cfLib.InitCfApi: %v\n", err)}

    // print results
    if dbg {cfLib.PrintApiObj (apiObj.ApiObj)}

	res, err := apiObj.VerifyToken()
	if err != nil {log.Fatalf("Verify Token: %v\n", err)}
	cfLib.PrintTokResp(res)
}
