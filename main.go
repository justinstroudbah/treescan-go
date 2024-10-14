package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alexflint/go-arg"
)

var Args struct {
	SourcePaths  string `arg:"-s" help:"Comma seperated list of file paths that will be scanned."`
	ReportPath   string `arg:"-o,env:TSGO_REPORT_PATH" help:"Where the reported scan results should be stored (aside from STDIO)"`
	Dump         bool   `arg:"-d" help:"Dump all results to stdout instead of scanning"`
	DumpFormat   string `arg:"-f" help:"Format of dump command"`
	ReportFormat string `arg:"-r" help:"Format of report command"`
	Debug        bool   `arg:"-x" help:"Enable debug mode"`
	Persist      bool   `arg:"-p" help:"Persist JavaScript VM session across a node type"`
	Languages    string `arg:"-l" help:"Comma separated list of languages"`
	ConfigFile   string `arg:"-c" help:"Optional custom configuration file"`
}

type FlatNode struct {
	NodeType   string
	LineNumber int
	RawSource  string
}

type RuleConfiguration struct {
	Rules []Rule `json:"rule"`
}

type Rule struct {
	Name         string `json:"name"`
	NodeType     string `json:"nodeType"`
	Priority     int    `json:"priority"`
	Enabled      bool   `json:"enabled"`
	ScriptPath   string `json:"scriptPath"`
	ScriptSource string `json:"scriptSource"`
}

var CurrentConfiguration []Rule

func main() {
	var startedAt = time.Now()
	arg.MustParse(&Args)

	config := loadConfiguration()

	CurrentConfiguration := config.Rules

	var _ = CurrentConfiguration

	var stoppedAt = time.Now()
	if Args.Debug {
		var stringRuntime = strconv.FormatFloat(stoppedAt.Sub(startedAt).Seconds(), 'f', -1, 64)
		println("Execution time: ", stringRuntime, "sec.")
	}

}

func loadConfiguration() RuleConfiguration {
	var ruleConfig RuleConfiguration

	var configJSONPath = `C:\repos\go\treescan-go\scripts\rules.json` // Change this to be configurable at run time
	jsonConfigFile, err := os.Open(configJSONPath)
	if err != nil {
		fmt.Println("Error opening configuration file:", configJSONPath)
	}
	defer jsonConfigFile.Close()

	jsonDecoder := json.NewDecoder(jsonConfigFile)
	err = jsonDecoder.Decode(&ruleConfig)
	if err != nil {
		fmt.Println("Error decoding configuration file:", err)
	}

	for ruleIndex := range len(ruleConfig.Rules) {
		thisRule := ruleConfig.Rules[ruleIndex]
		filePath := thisRule.ScriptPath
		fileSource := thisRule.ScriptSource

		if len(fileSource) > 0 {
			continue
		}

		sourceFromFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		thisRule.ScriptSource = string(sourceFromFile)
	}
	return ruleConfig
}
