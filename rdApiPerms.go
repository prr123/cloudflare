// rdApiPerms.go
// program that reads the ApiPermisionGroupsList
// Author: prr azul software
// Date: 31 May 2023
// copyright prr, azul software
//

package main

import (
//	"context"
	"fmt"
	"log"
	"os"

//    util "github.com/prr123/utility/utilLib"
    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr := "rdApiPerms [help]"

	if numArgs <1 {
		fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("no CLI args!\n")
    }
	if numArgs >2 {
		fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("too many CLI args!\n")
    }

	if numArgs == 2 {
		if os.Args[1] != "help" {
			fmt.Printf("usage: %s\n", useStr)
    	    log.Fatalf("invalid cli command %s!\n", os.Args[1])
		}
	}

	zoneDir := os.Getenv("zoneDir")
	if len(zoneDir) == 0 {log.Fatalf("cannot resolve env var zoneDir!\n")}

	permFilnam := zoneDir + "/apiPerms.yaml"
    log.Printf("Using yaml file:    %s\n", permFilnam)

	data, err := os.ReadFile(permFilnam)
	if err != nil {log.Fatalf("ReadFile: %v\n", err)}

	permList := cfLib.ApiPermList{}

	err = yaml.Unmarshal(data, &permList)
	if err != nil {log.Fatalf("Unmarshal: %v\n", err)}

	PrintApiPerms(&permList)

	log.Println("*** success rdApiPerms ***\n")
}

func PrintApiPerms(permList *cfLib.ApiPermList) {

	fmt.Printf("****** ApiPermList ******\n")
	fmt.Printf("perms: %d\n", len(permList.ApiPerms))
	fmt.Printf("**** end ApiPermList ****\n")
}
