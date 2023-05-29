// creDnsToken.go
// Author: prr, azul software
// Date 28 May 2023
// copyright prr, azul software
//
// usage creDnsToken
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
	"github.com/cloudflare/cloudflare-go"
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

	permGroup :=cloudflare.APITokenPermissionGroups {
		ID: "4755a26eedb94da69e1066d98aa820be",
		Name: "DNS Write",
		Scopes: nil,
	}

	permGroups := make([]cloudflare.APITokenPermissionGroups, 1)
	permGroups[0] = permGroup

	res := make(map[string]interface{})
	res["com.cloudflare.api.account.d0e0781201c0536742831e308ce406fb"] = "*"

	policy := cloudflare.APITokenPolicies{
			Effect: "allow",
			Resources: res,
			PermissionGroups: permGroups,
		}

	policies := make([]cloudflare.APITokenPolicies, 1)
	policies[0] = policy



	startTime := time.Now().UTC().Round(time.Second)
//.Format("2005-12-30T01:02:03Z")
	endTime := time.Now().UTC().AddDate(0,2,0).Round(time.Second)
//.Format("2005-12-30T01:02:03Z")


	// first we need to retrieve account
	tok:= cloudflare.APIToken{
		Name: "testToken",
		NotBefore: &startTime,
		ExpiresOn: &endTime,
		Policies: policies,
	}

	newTok, err := api.CreateAPIToken(ctx, tok)
	if err != nil {log.Fatalf("CreateApiToken: %v\n", err)}
//	tokList, err := api.APITokens(ctx)
//	if err != nil {log.Fatalf("APITokens: %v\n", err)}

	cfLib.PrintToken(newTok)
}
