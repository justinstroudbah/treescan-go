package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/antlr4-go/antlr/v4"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"treescan-go/parser"
)

type apexListener struct {
	*parser.BaseApexParserListener
}

var _sourceMap map[string]string

func (a apexListener) VisitTerminal(node antlr.TerminalNode) {

	//TODO implement me
}

func (a apexListener) VisitErrorNode(node antlr.ErrorNode) {
	//TODO implement me

}

func (a apexListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// Pass in reflect.TypeOf(x) and context
	vm := otto.New()
	vm.Set("START_LINE", ctx.GetStart().GetLine())
	vm.Set("STOP_LINE", ctx.GetStop().GetLine())
	vm.Set("CONTEXT", ctx.GetText())
	vm.Set("NODE_TYPE", reflect.TypeOf(ctx).String())
	vm.Run(`
	var c = CONTEXT.GetText();
	var x = 1;
	//console.log("Source: " + CONTEXT);
    //console.log("Type: " + NODE_TYPE); // 4
	`)
}

func (a apexListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	//TODO implement me
}

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

var args struct {
	SourcePaths  string `arg:"-s" help:"Comma seperated list of file paths that will be scanned."`
	ReportPath   string `arg:"-o,env:TSGO_REPORT_PATH" help:"Where the reported scan results should be stored (aside from STDIO)"`
	Dump         bool   `arg:"-d" help:"Dump all results to stdout instead of scanning"`
	DumpFormat   string `arg:"-f" help:"Format of dump command"`
	ReportFormat string `arg:"-r" help:"Format of report command"`
	Debug        bool   `arg:"-x" help:"Enable debug mode"`
	Languages    string `arg:"-l" help:"Comma separated list of languages"`
}

func main() {
	var startedAt = time.Now()
	arg.MustParse(&args)

	// Pretty great documentation on how to integrate Otto:
	// https://github.com/robertkrimen/otto

	_sourceMap = make(map[string]string)
	//"C:\\repos\\va-teams\\working\\va-teams\\force-app\\main\\default\\classes\\"
	var path = args.SourcePaths
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if args.Debug {
		scanfile("C:\\repos\\va-teams\\working\\va-teams\\force-app\\main\\default\\classes\\test_ServiceResponse.cls")
	}
	if !args.Debug {
		for _, file := range files {
			var name = file.Name()
			if strings.HasSuffix(name, ".cls") {
				scanfile(path + name)
				fmt.Println(file.Name())
			}
		}
	}
	var stoppedAt = time.Now()
	var runTime = stoppedAt.Sub(startedAt)
	var secondsToRun = runTime.Seconds()
	var stringRuntime = strconv.FormatFloat(secondsToRun, 'f', -1, 64)
	println("Execution time: %s\n", stringRuntime)
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
