package main

import (
	"github.com/alexflint/go-arg"
	"github.com/antlr4-go/antlr/v4"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"treescan-go/parser"
)

//type SourceNode struct {
//	Name     string
//	StartRow int
//	StartCol int
//	EndRow   int
//	EndCol   int
//	Text     string
//	Parent   SourceNode
//	Children []SourceNode
//}

type apexListener struct {
	*parser.BaseApexParserListener
}

func (a apexListener) VisitTerminal(node antlr.TerminalNode) {

	//TODO implement me
}

func (a apexListener) VisitErrorNode(node antlr.ErrorNode) {
	//TODO implement me

}

func (a apexListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// Pass in reflect.TypeOf(x) and context
	vm := otto.New()
	vm.Set("context", ctx)
	vm.Run(`
    abc = 2 + 2;
	console.log(context.GetParent().GetText());
    console.log("The value of abc is " + abc); // 4
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
	arg.MustParse(&args)

	// Pretty great documentation on how to integrate Otto:
	// https://github.com/robertkrimen/otto

	data, err := ioutil.ReadFile("C:\\repos\\va-teams\\working\\va-teams\\force-app\\main\\default\\classes\\VCR\\Tests\\VCR_VisitRepoTest.cls")
	if err != nil {
		log.Fatal(err)
	}
	is := antlr.NewInputStream(string(data))
	lexer := parser.NewApexLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewApexParser(stream)
	parser.ApexParserInit()
	p.BuildParseTrees = true

	antlr.ParseTreeWalkerDefault.Walk(&apexListener{}, p.CompilationUnit())
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
