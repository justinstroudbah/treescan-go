package main

import (
	"github.com/alexflint/go-arg"   // CLI argument utility
	"github.com/antlr4-go/antlr/v4" // Go runtime for ANTLR
	"github.com/robertkrimen/otto"  // JavaScript engine
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"treescan-go/parser" // Generated ANTLR parser
	"treescan-go/treescanner"
	"treescan-go/util"
)

type Violation struct {
	NodeTYpeName   string
	SourceFileName string
	LineNumber     int
	Message        string
}

type apexListener struct {
	*parser.BaseApexParserListener
}

var SourceMap map[string]util.Rule

func (a apexListener) VisitTerminal(node antlr.TerminalNode) {

	//TODO implement me
}

func (a apexListener) VisitErrorNode(node antlr.ErrorNode) {
	//TODO implement me

}

func (a apexListener) EnterEveryRule(ctx antlr.ParserRuleContext) {

	// Set up convenience variables to be used by rule scripts
	var nodeType = reflect.TypeOf(ctx).String()
	var contextSource = ctx.GetText()
	var startLine = ctx.GetStart().GetLine()
	var stopLine = ctx.GetStop().GetLine()

	nodeTypeCompact := strings.Replace(nodeType, "*parser.", "", -1)

	vm := otto.New()
	vm.Set("START_LINE", startLine)
	vm.Set("STOP_LINE", stopLine)
	vm.Set("SOURCE", contextSource)
	vm.Set("CONTEXT", ctx)
	vm.Set("NODE_TYPE", nodeTypeCompact)

	if values, ok := util.ScriptSources[nodeTypeCompact]; ok {
		for _, value := range values {
			vm.Run(value)
		}
	}
}

func (a apexListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	//if Args.Debug {
	//	fmt.Println(reflect.TypeOf(ctx).String())
	//}
}

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

var Args struct {
	SourcePaths  string `arg:"-s" help:"Comma seperated list of file paths that will be scanned."`
	ReportPath   string `arg:"-o,env:TSGO_REPORT_PATH" help:"Where the reported scan results should be stored (aside from STDIO)"`
	Dump         bool   `arg:"-d" help:"Dump all results to stdout instead of scanning"`
	DumpFormat   string `arg:"-f" help:"Format of dump command"`
	ReportFormat string `arg:"-r" help:"Format of report command"`
	Debug        bool   `arg:"-x" help:"Enable debug mode"`
	Persist      bool   `arg:"-p" help:"Persist JavaScript VM session across a node type"`
	Languages    string `arg:"-l" help:"Comma separated list of languages"`
}

func main() {
	var startedAt = time.Now()
	arg.MustParse(&Args)

	// Pretty great documentation on how to integrate Otto:
	// https://github.com/robertkrimen/otto

	//"C:\\repos\\va-teams\\working\\va-teams\\force-app\\main\\default\\classes\\"
	var path = Args.SourcePaths
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	//SourceMap := make(map[string]util.Rule)

	SourceMap = util.ReadRuleConfiguration()

	if Args.Debug {
		//scanfile("C:\\repos\\va-teams\\working\\va-teams\\force-app\\main\\default\\classes\\test_ServiceResponse.cls")
	} else {
		for _, file := range files {
			var name = file.Name()
			if strings.HasSuffix(name, ".cls") {
				//scanfile(path + name)
			}
		}
	}
	var stoppedAt = time.Now()
	var runTime = stoppedAt.Sub(startedAt)
	var secondsToRun = runTime.Seconds()
	var stringRuntime = strconv.FormatFloat(secondsToRun, 'f', -1, 64)
	if Args.Debug {
		println("Execution time: ", stringRuntime, "sec.")
	}

	var s = treescanner.NewScanManager()
	s.Init()

	//mgr.Scan()
}

func scanfile(fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var sourceString = string(data)

	is := antlr.NewInputStream(sourceString)
	lexer := parser.NewApexLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewApexParser(stream)
	parser.ApexParserInit()
	p.BuildParseTrees = true

	antlr.ParseTreeWalkerDefault.Walk(&apexListener{}, p.CompilationUnit())

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
