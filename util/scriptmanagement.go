package util

import (
	"encoding/json"
	"fmt"
	"os"
)

var Config Rules
var ScriptSources map[string][]string

type Rules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	NodeType string   `json:"nodeType"`
	Enabled  bool     `json:"enabled"`
	Priority int      `json:"priority"`
	Scripts  []string `json:"scripts"`
}

func ReadRuleConfiguration() map[string]Rule {

	var scriptMap map[string]Rule
	scriptMap = make(map[string]Rule)
	ScriptSources = make(map[string][]string)

	var configJSONPath = `C:\repos\go\treescan-go\scripts\rules.json`

	file, err := os.Open(configJSONPath)
	if err != nil {
		fmt.Println("Error opening configuration file:", configJSONPath)
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error decoding configuration file:", configJSONPath)
	}

	for _, rule := range Config.Rules {

		scriptMap[rule.NodeType] = rule
		for _, script := range rule.Scripts {
			source, err := os.ReadFile(script)
			if err != nil {
				fmt.Println("Error reading file:", script)
			}
			if values, ok := ScriptSources[rule.NodeType]; ok {
				ScriptSources[rule.NodeType] = append(values, string(source))
			} else {
				ScriptSources[rule.NodeType] = []string{string(source)}
			}
		}
	}

	return scriptMap
}
