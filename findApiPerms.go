// findApiPerms.go
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
	"strings"

//    util "github.com/prr123/utility/utilLib"
    yaml "github.com/goccy/go-yaml"
	"ns/cloudflare/cfLib"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)
	useStr := "findApiPerms [expression/help]"
	searchTerm :=""

	if numArgs <2 {
		fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("insufficient CLI args!\n")
    }
	if numArgs >2 {
		fmt.Printf("usage: %s\n", useStr)
        log.Fatalf("too many CLI args!\n")
    }

	if numArgs == 2 {
		if os.Args[1] == "help" {
			fmt.Printf("usage: %s\n", useStr)
			fmt.Printf("program looks through all permissions located in the ApiPerm file to find paermissions that include the search term\n")
			os.Exit(1)
		}
		searchTerm = os.Args[1]
	}

	log.Printf("search term: \"%s\"\n", searchTerm)

	zoneDir := os.Getenv("zoneDir")
	if len(zoneDir) == 0 {log.Fatalf("cannot resolve env var zoneDir!\n")}

	permFilnam := zoneDir + "/apiPerms.yaml"
    log.Printf("Using apiPerms file:    %s\n", permFilnam)

	data, err := os.ReadFile(permFilnam)
	if err != nil {log.Fatalf("ReadFile: %v\n", err)}

	permList := cfLib.ApiPermList{}

	err = yaml.Unmarshal(data, &permList)
	if err != nil {log.Fatalf("Unmarshal: %v\n", err)}

	log.Printf("success unmarshaling permList\n")

	for i:=0; i< len(permList.ApiPerms); i++ {
		nam := permList.ApiPerms[i].Name
		idx := strings.Index(nam, searchTerm)
		if idx > -1 {
			fmt.Printf("perm[%3d]: name: %-25s id: %s \n", i+1, nam, permList.ApiPerms[i].Id)
		}

	}
//	PrintApiPermsSummary(&permList)

	log.Println("*** success findApiPerms ***\n")
}

func PrintApiPermsSummary(permList *cfLib.ApiPermList) {

	fmt.Printf("****** ApiPermList ******\n")
	fmt.Printf("perms: %d\n", len(permList.ApiPerms))
	fmt.Printf("**** end ApiPermList ****\n")
}
