// interim package to test Cloudflare Yaml read
//

package main

import (
    "fmt"
    "os"
    yaml "github.com/goccy/go-yaml"
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

	fmt.Printf("Using yaml file: %s\n", yamlFilNam)

    // open file and decode

    buf, err := os.ReadFile(yamlFilNam)
    if err != nil {
        fmt.Printf("cannot open yaml File: os.ReadFile: %v\n", err)
        os.Exit(-1)
    }

//    fmt.Printf("buf [%d]:\n%s\n", len(buf), string(buf))

    if err := yaml.Unmarshal(buf, &apiObj); err !=nil {
        fmt.Printf("error Unmarshalling Yaml File: %v\n", err)
        os.Exit(-1)
    }

    // print results
    PrintApiObj (&apiObj)
}

func PrintApiObj (apiObj *ApiObj) {

    fmt.Println("********** Api Obj ************")
    fmt.Printf("API:     %s\n", apiObj.Api)
    fmt.Printf("APIKey:  %s\n", apiObj.ApiKey)
//    fmt.Printf("Ca Dir Url: %s\n", nchObj.CADirUrl)
    fmt.Printf("Email:   %s\n", apiObj.Email)
    fmt.Println("*******************************")
}




