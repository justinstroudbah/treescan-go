package treescanner

import "fmt"

type Violation struct {
	SourceFilePath string
	LineNumber     int
	Message        string
	Priority       int
}

type Rules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	RuleName    string   `json:"name"`
	NodeType    string   `json:"node_type"`
	Priority    int      `json:"priority"`
	Enabled     bool     `json:"enabled"`
	ScriptPaths []string `json:"scripts"`
}

var violations map[string][]Violation

type ScanManager struct {
}

func NewScanManager() *ScanManager {
	config = make(map[string][]Configuration)
	violations = make(map[string][]Violation)
	manager := new(ScanManager)
	return manager

}

func (scanner *ScanManager) Init() bool {
	var configPath = "C:\\repos\\go\\treescan-go\\scripts\\rules.json"
	var foo map[string][]Configuration
	foo = make(map[string][]Configuration)
	foo = readRuleConfiguration(configPath)

	for key, value := range foo {
		println(key)
		fmt.Println(value, "\n")
	}
	return true
}

func (scanner *ScanManager) Scan(pathList []string) map[string][]Violation {

	return violations
}

func (scanner *ScanManager) GetConfiguration(pathList []string) map[string][]Configuration {
	return config
}
