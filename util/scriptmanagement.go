package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Rules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	NodeType string   `json:"nodeType"`
	Scripts  []string `json:"scripts"`
}

func ReadRuleConfiguration() map[string][]string {

	var scriptMap map[string][]string
	scriptMap = make(map[string][]string)

	var configJSONPath = `C:\repos\go\treescan-go\scripts\rules.json`

	file, err := os.Open(configJSONPath)
	if err != nil {
		fmt.Println("Error opening configuration file:", configJSONPath)
	}
	defer file.Close()

	var rules Rules
	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&rules)
	if err != nil {
		fmt.Println("Error decoding configuration file:", configJSONPath)
	}

	for _, rule := range rules.Rules {
		for _, script := range rule.Scripts {
			sourceData, err := ioutil.ReadFile(script)
			if err != nil {
				fmt.Println("Error opening script file:", script)
			}
			scriptMap[rule.NodeType] = append(scriptMap[rule.NodeType], string(sourceData))
		}
	}

	return scriptMap
}
