// creTokFile.go
// program that create a token file 
// Author: prr, azulsoftware
// Date: 12 July 2023
// copyright 2023 prr, azul software
//

package main

import (
//	"context"
	"fmt"
	"log"
	"os"
//	"time"
	"strings"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "creTokFile /token=token /file=yamlfile [/dbg]"
	helpStr := "program create a yaml token file from a token and verify the token."
	dbg := false

    cfDir := os.Getenv("Cloudflare")
    if len(cfDir) == 0 {log.Fatalf("could not resolve Cloudflare\n")}

//    cfTokenFilnam := cfDir + "/token/"

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf("usage: %s\n", useStr)
		os.Exit(1)
	}

    if numArgs < 2 {
        fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs > 4 {
        fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("too many CLI args!\n")
    }

	flags := []string{"token", "file", "dbg"}

	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("error parseFlags: %v\n",err)
    }

	val, ok := flagMap["file"]
	if !ok {
        fmt.Printf("usage: %s\n", useStr)
		log.Fatalf("no file flag!")
	}

	outFilnam, ok := val.(string)
	if len(outFilnam) == 0 {
        fmt.Printf("usage: %s\n", useStr)
		log.Fatalf("no filename provided with /file flag!")
	}
	if idx:=strings.Index(outFilnam, ".yaml"); idx == -1 {outFilnam += ".yaml"}

	tokval, ok := flagMap["token"]
	if !ok {
        fmt.Printf("usage: %s\n", useStr)
		log.Fatalf("no token flag!")
	}
	token := tokval.(string)

	if len(token) == 0 {
        fmt.Printf("usage: %s\n", useStr)
		log.Fatalf("no token provided with /token flag!")
	}

	_, ok = flagMap["dbg"]
	if ok { dbg = true }

	cfTokenFilnam := outFilnam

	if dbg {
		for k,v := range flagMap {
			fmt.Printf("flag: /%s value: %s\n", k, v)
		}
		log.Printf("Using token: %s\n", token)
		log.Printf("Using yaml file: %s\n", cfTokenFilnam)
	}

//	fmt.Println("************** token *********************")


//	cfLib.PrintToken(token)

    fmt.Println("********************************************")

	err = cfLib.CreateTokFile(cfTokenFilnam, token, dbg)
	if err != nil {log.Fatalf("CreateTokenFilnam: %v", err) }

	log.Printf("success creating token file!\n")
}
