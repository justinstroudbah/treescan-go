package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"treescan-go/parser"

	"github.com/alexflint/go-arg"
	"github.com/antlr4-go/antlr/v4"
	"github.com/robertkrimen/otto"
)

type Violation struct {
	ScannerRule   Rule
	ParserContext antlr.ParserRuleContext
	StartLine     int
	StopLine      int
	SourceFile    string
	SourceCode    string
}

var Violations []Violation

type EventDump struct {
	NodeType   string
	StartLine  int
	StopLine   int
	SourceCode string
	FileName   string
}

// apexListener is where we inherit our events from, which in turn provide entry points for the scan
type apexListener struct {
	*parser.BaseApexParserListener
	FileName string
}

var Args struct {
	SourcePaths  string `arg:"-s" help:"Comma seperated list of file and directory paths that will be scanned."`
	Recursive    bool   `arg:"-v" help:"Recurse through directories supplied by SourcePaths"`
	ReportPath   string `arg:"-o,env:TSGO_REPORT_PATH" help:"Where the reported scan results should be stored (aside from STDIO)"`
	Dump         bool   `arg:"-d" help:"Dump all results to stdout instead of scanning"`
	DumpFormat   string `arg:"-f" help:"Format of dump command"`
	ReportFormat string `arg:"-r" help:"Format of report command"`
	Debug        bool   `arg:"-x" help:"Enable debug mode"`
	Persist      bool   `arg:"-p" help:"Persist JavaScript VM session across a node type"`
	Languages    string `arg:"-l" help:"Comma separated list of languages"`
	ConfigFile   string `arg:"-c" help:"Optional custom configuration file"`
}

var SourceMap map[string]string

type RuleConfiguration struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Name         string `json:"name"`
	NodeType     string `json:"nodeType"`
	Message      string `json:"message"`
	Description  string `json:"description"`
	Priority     int    `json:"priority"`
	Enabled      bool   `json:"enabled"`
	ScriptPath   string `json:"scriptPath"`
	ScriptSource string `json:"scriptSource"`
}

var RuleConfig RuleConfiguration

func main() {
	var startedAt = time.Now()
	arg.MustParse(&Args)
	SourceMap = make(map[string]string)

	// Get configuration, persist it somewhere
	loadConfiguration()

	// This is just to keep the compiler from barking at us for an unused variable

	sources := getSourceFiles(Args.SourcePaths, "cls")

	for sourceFileIndex := range len(sources) {
		sourceFile := sources[sourceFileIndex]
		ParseFile(sourceFile)
	}

	if len(Violations) > 0 {
		jsonData, err := json.MarshalIndent(Violations, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(jsonData))
	}
	// Metrics if we need them
	var stoppedAt = time.Now()
	if Args.Debug {
		var stringRuntime = strconv.FormatFloat(stoppedAt.Sub(startedAt).Seconds(), 'f', -1, 64)
		println("Execution time: ", stringRuntime, "sec.")
	}

}

func loadConfiguration() {

	var configJSONPath = `C:\repos\go\treescan-go\scripts\rules.json`
	if len(Args.ConfigFile) > 0 {
		configJSONPath = Args.ConfigFile
	}
	jsonConfigFile, err := os.Open(configJSONPath)
	if err != nil {
		fmt.Println("Error opening configuration file:", configJSONPath)
	}
	defer jsonConfigFile.Close()

	jsonDecoder := json.NewDecoder(jsonConfigFile)
	err = jsonDecoder.Decode(&RuleConfig)
	if err != nil {
		fmt.Println("Error decoding configuration file:", err)
	}

	for ruleIndex := range len(RuleConfig.Rules) {
		thisRule := RuleConfig.Rules[ruleIndex]
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
		RuleConfig.Rules[ruleIndex] = thisRule
	}
}

func getSourceFiles(rawPaths string, filterExtension string) []string {
	pathNames := strings.Split(rawPaths, ",")
	var sourceFilePaths []string

	for pathNameIndex := range len(pathNames) {
		pathName := pathNames[pathNameIndex]
		info, err := os.Stat(pathName)
		if err != nil {
			fmt.Println("File path to source code is invalid: ", pathName)
			log.Fatal(err)
		}
		if info.IsDir() {
			err := filepath.Walk(pathName,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if strings.HasSuffix(path, filterExtension) {
						sourceFilePaths = append(sourceFilePaths, path)
					}
					return nil
				})
			if err != nil {
				log.Println(err)
			}
		} else {
			sourceFilePaths = append(sourceFilePaths, pathName)
		}
	}
	return sourceFilePaths
}

func (a apexListener) VisitTerminal(node antlr.TerminalNode) {

	// Not implemented
}

func (a apexListener) VisitErrorNode(node antlr.ErrorNode) {
	// Not implemented

}

// EnterEveryRule fires on all nodes in the tree. This is where we can determine whether or not the rules are interested in them by checking against the node type (the RuleContext)
func (a apexListener) EnterEveryRule(ctx antlr.ParserRuleContext) {

	// Set up convenience variables to be used by rule scripts
	var nodeType = reflect.TypeOf(ctx).String()
	var contextSource = ctx.GetText()
	var rawSource = SourceMap[a.FileName]
	var startLine = ctx.GetStart().GetLine()
	var stopLine = ctx.GetStop().GetLine()
	nodeTypeCompact := strings.Replace(nodeType, "*parser.", "", -1)

	vm := otto.New()
	vm.Set("START_LINE", startLine)
	vm.Set("STOP_LINE", stopLine)
	vm.Set("SOURCE", contextSource)
	vm.Set("CONTEXT", ctx)
	vm.Set("NODE_TYPE", nodeTypeCompact)
	vm.Set("FORMATTED_SOURCE", rawSource)

	var configCount = len(RuleConfig.Rules)

	for contextRuleIndex := range configCount {
		thisRule := RuleConfig.Rules[contextRuleIndex]
		var hasViolation = false
		if nodeTypeCompact == thisRule.NodeType {
			vm.Run(thisRule.ScriptSource)
			if value, err := vm.Get("HAS_VIOLATION"); err == nil {
				if value_bool, err := value.ToBoolean(); err == nil {
					hasViolation = value_bool
				}
			}
		}
		if hasViolation {
			var contextCodeLines = strings.Split(rawSource, "\r\n")
			//var sourceLineList = contextCodeLines[startLine:stopLine]
			var sourceLines = strings.Join(contextCodeLines[startLine-1:stopLine], "\n")
			var violationToAdd = new(Violation)
			violationToAdd.StartLine = startLine
			violationToAdd.StopLine = stopLine
			violationToAdd.SourceFile = a.FileName
			violationToAdd.SourceCode = sourceLines
			violationToAdd.ScannerRule = thisRule
			violationToAdd.ParserContext = ctx
			violationToAdd.ScannerRule.Message = thisRule.Message
			violationToAdd.ScannerRule.NodeType = nodeTypeCompact
			violationToAdd.ScannerRule.Name = thisRule.Name
			Violations = append(Violations, *violationToAdd)
		}
	}

}

func (a apexListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
}

func (config *RuleConfiguration) GetRulesForNode(nodeType string) []Rule {
	var result []Rule
	for ruleIndex := range len(RuleConfig.Rules) {
		thisRule := RuleConfig.Rules[ruleIndex]
		splitNodeType := strings.Split(nodeType, ",")
		nodeTypeCompact := splitNodeType[len(splitNodeType)-1]

		if thisRule.NodeType == nodeTypeCompact && thisRule.Enabled == true {
			result = append(result, thisRule)
		}
	}
	return result
}

func ParseFile(fileName string) {

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var sourceString = string(data)

	// We do want to save the source code (for now.) Ideally need to put off loading the full source until absolutely necessary
	SourceMap[fileName] = sourceString

	// Apex is NOT case sensitive, so that presents a problem...especially with SELECT and standard query conventions.
	// Maybe use a regex for this in the long term
	sourceStringForParsing := strings.ToLower(sourceString)

	sourceInputStream := antlr.NewInputStream(sourceStringForParsing)
	lexer := parser.NewApexLexer(sourceInputStream)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewApexParser(stream)
	parser.ApexParserInit()
	p.BuildParseTrees = true
	listener := apexListener{}
	listener.FileName = fileName

	antlr.ParseTreeWalkerDefault.Walk(&listener, p.CompilationUnit())

}
