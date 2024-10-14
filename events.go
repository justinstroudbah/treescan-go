package main

import (
	"reflect"
	"strings"
	"treescan-go/parser"
	"treescan-go/util"

	"github.com/antlr4-go/antlr/v4" // Go runtime for ANTLR
	"github.com/robertkrimen/otto"
)

var jsonDump []string

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
