package treescanner

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ScanRule     Rule
	ScriptPath   string
	ScriptSource string
}

var config map[string][]Configuration

func readRuleConfiguration(configFilePath string) map[string][]Configuration {
	var fullConfigs map[string][]Configuration
	fullConfigs = make(map[string][]Configuration)
	var jsonRules Rules

	var configJSONPath = configFilePath

	file, err := os.Open(configJSONPath)
	if err != nil {
		fmt.Println("Error opening configuration file:", configJSONPath)
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&jsonRules)
	if err != nil {
		fmt.Println("Error decoding configuration file:", configJSONPath)
	}

	for _, rule := range jsonRules.Rules {

		for _, script := range rule.ScriptPaths {
			rawSource, err := os.ReadFile(script)
			if err != nil {
				fmt.Println("Error reading file:", script)
			}
			var config Configuration
			config.ScanRule = rule
			config.ScriptSource = string(rawSource)
			config.ScriptPath = script

			if values, ok := fullConfigs[rule.NodeType]; ok {
				fullConfigs[rule.NodeType] = append(values, config)
			} else {
				fullConfigs[rule.NodeType] = []Configuration{config}
			}
		}
	}
	return fullConfigs
}
