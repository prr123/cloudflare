package main

import (
	"context"
	"fmt"
	"log"
	"os"
    yaml "github.com/goccy/go-yaml"
	"github.com/cloudflare/cloudflare-go"
)

type ApiObj struct {
    Api    string `yaml:"Api"`
    ApiKey string `yaml:"ApiKey"`
//    CADirUrl  string `yaml:"CA_DIR_URL"`
    Email     string `yaml:"Email"`
}


func main() {

    var apiObj ApiObj

    numArgs := len(os.Args)
    if numArgs > 2 {
        fmt.Printf("CLI args are not equal to 2: %d\n", numArgs)
        fmt.Printf("usage: readCFYaml yaml file\n")
        os.Exit(-1)
    }

    yamlFilNam := "cloudflareApi.yaml"

    if numArgs == 2 {yamlFilNam = os.Args[1]}

    log.Printf("Using yaml file: %s\n", yamlFilNam)

    // open file and decode
    buf, err := os.ReadFile(yamlFilNam)
    if err != nil {
        log.Fatalf("cannot open yaml File: os.ReadFile: %v\n", err)
    }

//    fmt.Printf("buf [%d]:\n%s\n", len(buf), string(buf))

    if err := yaml.Unmarshal(buf, &apiObj); err !=nil {
        log.Fatalf("error Unmarshalling Yaml File: %v\n", err)
    }

    // print results
    PrintApiObj (&apiObj)

	// Construct a new API object using a global API key
//	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	// alternatively, you can use a scoped API token
	cloudToken := "O5ART89fgxulItZ1l-o9PScX-uEGXN219dzo06Xi"
	api, err := cloudflare.NewWithAPIToken(cloudToken)
	if err != nil {
		log.Fatalf("api init: %v/n", err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatalf("api.UserDetails: %v\n", err)
	}

	// Print user details
	fmt.Println(u)

	fmt.Println("********************************************")

	PrintUserInfo(&u)
/*
	// try to list DNS Parameters

	var listDns cloudflare.ListDNSRecordsParams
	listDns.Name = "azulacademy.eu"

	var rc cloudflare.ResourceContainer
	rc.Level = cloudflare.ZoneRouteLevel
	rc.Identifier = "0e6e30d5edb4c1025817eb1678511cef"

	dnsRecs, resInfo, err := api.ListDNSRecords(ctx,&rc , listDns)
    if err != nil {
        log.Fatalf("api.ListDNSRecords: %v\n", err)
    }
	fmt.Printf("resInfo: %v\n", resInfo)
	fmt.Printf("Dns Records [%d]\n",len(dnsRecs))

	PrintResInfo(resInfo)
	PrintDnsRec(&dnsRecs)
*/
}

// https://github.com/cloudflare/cloudflare-go/blob/0d05fc09483641dde8abb4c64cf2f6016f590d79/user.go#L12
func PrintUserInfo (u *cloudflare.User) {

	var actTyp string

	fmt.Println("************** User Info **************")
	fmt.Printf("First Name:  %s\n", u.FirstName)
	fmt.Printf("Last Name:   %s\n", u.LastName)
	fmt.Printf("Email:       %s\n", u.Email)
	fmt.Printf("ID:          %s\n", u.ID)
	fmt.Printf("Country:     %s\n", u.Country)
	fmt.Printf("Zip Code:    %s\n", u.Zipcode)
	fmt.Printf("Phone:       %s\n", u.Telephone)
	fmt.Printf("2FA:         %t\n", u.TwoFA)
	timStr := (u.CreatedOn).Format("02 Jan 06 15:04 MST")
	fmt.Printf("Created:     %s\n", timStr)
	timStr = (u.ModifiedOn).Format("02 Jan 06 15:04 MST")
	fmt.Printf("Modified:    %s\n", timStr)
	fmt.Printf("ApiKey:      %s\n", u.APIKey)
	if len(u.Accounts) == 1 {
		act := u.Accounts[0]
		actTyp = act.Type
		if len(actTyp) == 0 {actTyp = "-"}
		fmt.Printf("account ID: %s Name: %s Type: %s\n", act.ID, act.Name, actTyp)
	} else {
		fmt.Printf("Accounts [%d]:\n", len(u.Accounts))
		fmt.Printf("Nu ID  Name  Type\n")
		for i:=0; i< len(u.Accounts); i++ {
			act := u.Accounts[i]
			actTyp = act.Type
			if len(actTyp) == 0 {actTyp = "-"}
			fmt.Printf("%d: %s %s %s\n", i+1, act.ID, act.Name, actTyp)
		}
	}
	fmt.Println("********** End User Info **************")

}

func PrintResInfo(res *cloudflare.ResultInfo) {

	fmt.Println("************** ResultInfo **************")
	fmt.Printf("Page:       %d\n", res.Page)
	fmt.Printf("PerPage:    %d\n", res.PerPage)
	fmt.Printf("TotalPages: %d\n", res.TotalPages)
	fmt.Printf("Count:      %d\n", res.Count)
	fmt.Printf("Total:      %d\n", res.Total)
	fmt.Println("********** End ResultInfo **************")
}

func PrintDnsRec(recs *[]cloudflare.DNSRecord) {
	fmt.Printf("************* DNS Recourds: %d ************\n", len(*recs))
	fmt.Println("number  type      name             value/ content            Id")
	for i:=0; i< len(*recs); i++ {
		fmt.Printf("Record[%d]: %-3s %s %s %s\n", i, (*recs)[i].Type, (*recs)[i].Name, (*recs)[i].Content, (*recs)[i].ID)
	}
}

func PrintApiObj (apiObj *ApiObj) {

    fmt.Println("********** Api Obj ************")
    fmt.Printf("API:     %s\n", apiObj.Api)
    fmt.Printf("APIKey:  %s\n", apiObj.ApiKey)
//    fmt.Printf("Ca Dir Url: %s\n", nchObj.CADirUrl)
    fmt.Printf("Email:   %s\n", apiObj.Email)
    fmt.Println("*******************************")
}
