// cloudflare support library
// Author: prr, azulsoftware
// Date: 20 Mar 2023
// copyright 2023 prr, azul software
//
package main

import (
    "os"
    "log"
	"fmt"

	"ns/cloudflare/cfLib"

)

func main() {

	numArgs := len(os.Args)
	if numArgs > 2 {
		fmt.Printf("CLI args are not equal to 2: %d\n", numArgs)
		fmt.Printf("usage: readCFYaml yaml file\n")
		os.Exit(-1)
	}

	yamlFilNam := "cloudflareApi.yaml"

	if numArgs == 2 {yamlFilNam = os.Args[1]}

	log.Printf("Using yaml file: %s\n", yamlFilNam)

	apiObj, err := cfLib.InitCfLib(yamlFilNam)
	if err != nil {
		log.Fatalf("cfLib.InitCfLib: %v\n", err)
	}
    // print results
    cfLib.PrintApiObj (apiObj)
}

