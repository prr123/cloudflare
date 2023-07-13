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

	dbg:= false
    numArgs := len(os.Args)

//	cfDir := os.Getenv("Cloudflare")
//   if len(cfDir) == 0 {log.Fatalf("could not get env: CloudFlare\n")}

//    cfApiFilnam := cfDir + "/token/cfTok.yaml"
    cfApiFilnam := "cfTok.yaml"

	useStr := "usage: creDnsToken /out=file [/dbg]\n"
	helpStr := "this program creates a new  token with the ability to change Dns records."

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf(useStr)
		os.Exit(1)
	}

    if numArgs == 1 {
		fmt.Printf("no flags provided!")
		fmt.Printf(useStr)
		os.Exit(1)
	}

	flags := []string{"out", "dbg"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("ParseFlags: %v\n",err)
    }

	val, ok := flagMap["out"]
	if !ok {
		log.Fatalf("/out flag is missing!")
	}
	tokFilnam, _ := val.(string)
	if tokFilnam == "none" {
		log.Fatalf("toke file value not provided!")
	}
	newTokenFilnam := tokFilnam

	_, ok = flagMap["dbg"]
	if ok {
		dbg = true
	}

	if dbg {
		log.Printf("Using token apifile: %s\n", cfApiFilnam)
		log.Printf("New token file:      %s\n", newTokenFilnam)
		log.Printf("debug: %t\n", dbg)
	}

    apiObj, err := cfLib.InitCfApi(cfApiFilnam)
    if err != nil {log.Fatalf("cfLib.InitCfApi: %v\n", err)}

    // print results
    if dbg {cfLib.PrintApiObj(apiObj.ApiObj)}

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


	actStr := "com.cloudflare.api.account." + apiObj.ApiObj.AccountId
	if dbg {fmt.Printf("Account: %s\n", actStr)}

	res := make(map[string]interface{})
//	res["com.cloudflare.api.account.d0e0781201c0536742831e308ce406fb"] = "*"
	res[actStr] = "*"

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
		Name: "TestDnsChange",
		NotBefore: &startTime,
		ExpiresOn: &endTime,
		Policies: policies,
	}

	newTok, err := api.CreateAPIToken(ctx, tok)
	if err != nil {log.Fatalf("CreateApiToken: %v\n", err)}

	if dbg {cfLib.PrintToken(newTok)}

    err = cfLib.CreateTokFile(newTokenFilnam, newTok.Value, dbg)
    if err != nil {log.Fatalf("CreateTokFile: %v", err) }

	log.Printf("success creating Dns Token!")
}
