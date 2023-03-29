// cloudflare support library
// Author: prr, azul software
// Date: 29 March 2023
// Copyright 2023 prr, azul software

package cfLib

import (
	"fmt"
	"os"
    yaml "github.com/goccy/go-yaml"
)


type ApiObj struct {
    Api    string `yaml:"Api"`
    ApiKey string `yaml:"ApiKey"`
    ApiToken string `yaml:"ApiToken"`
//    CADirUrl  string `yaml:"CA_DIR_URL"`
    Email     string `yaml:"Email"`
	YamlFile	string
}

func InitClLib(yamlFilNam string) (apiObjRef *ApiObj, err error) {
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


func PrintApiObj (apiObj *ApiObj) {

    fmt.Println("********** Api Obj ************")
    fmt.Printf("API:      %s\n", apiObj.Api)
    fmt.Printf("APIKey:   %s\n", apiObj.ApiKey)
    fmt.Printf("APIToken: %s\n", apiObj.ApiToken)
//    fmt.Printf("Ca Dir Url: %s\n", nchObj.CADirUrl)
    fmt.Printf("Email:    %s\n", apiObj.Email)
    fmt.Println("*******************************")
}
