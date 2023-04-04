// cloudflare support library
// Author: prr, azul software
// Date: 29 March 2023
// Copyright 2023 prr, azul software

package cfLib

import (
	"fmt"
	"os"

    yaml "github.com/goccy/go-yaml"
    "github.com/cloudflare/cloudflare-go"
	json "github.com/goccy/go-json"
)


type ApiObj struct {
    Api    string `yaml:"Api"`
    ApiKey string `yaml:"ApiKey"`
    ApiToken string `yaml:"ApiToken"`
	AccountId string `yaml:"AccountId"`
//    CADirUrl  string `yaml:"CA_DIR_URL"`
    Email     string `yaml:"Email"`
	YamlFile	string
}

func InitCfLib(yamlFilNam string) (apiObjRef *ApiObj, err error) {
	var apiObj ApiObj

    // open file and decode
    buf, err := os.ReadFile(yamlFilNam)
    if err != nil {
        return nil, fmt.Errorf("cannot open yaml File: os.ReadFile: %v\n", err)
    }

//    fmt.Printf("buf [%d]:\n%s\n", len(buf), string(buf))

    if err := yaml.Unmarshal(buf, &apiObj); err !=nil {
		return nil, fmt.Errorf("error Unmarshalling Yaml File: %v\n", err)
    }

	if apiObj.Api != "cloudflare" {
		return nil, fmt.Errorf("Api is not cloudflare!")
	}

	apiObj.YamlFile = yamlFilNam

	return &apiObj, nil
}

func SaveZones(zones []cloudflare.Zone, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}
	jsdata, err := json.Marshal(zones)
	if err != nil {return fmt.Errorf("could not convertZone: %v", err)}

	_, err = outfil.Write(jsdata)
	if err != nil {return fmt.Errorf("jsdata os.Write: %v", err)}
	return nil
}

func PrintZones(zones []cloudflare.Zone) {

    fmt.Printf("************** Zones/Domains [%d] *************\n", len(zones))

    for i:=0; i< len(zones); i++ {
        zone := zones[i]
        fmt.Printf("%d %-20s %s\n",i+1, zone.Name, zone.ID)
    }
}



func PrintApiObj (apiObj *ApiObj) {

    fmt.Println("********** Api Obj ************")
    fmt.Printf("API:       %s\n", apiObj.Api)
    fmt.Printf("APIKey:    %s\n", apiObj.ApiKey)
    fmt.Printf("APIToken:  %s\n", apiObj.ApiToken)
    fmt.Printf("AccountId: %s\n", apiObj.AccountId)
    fmt.Printf("Email:     %s\n", apiObj.Email)
    fmt.Println("*******************************")
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
