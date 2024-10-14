package main

import (
	"reflect"
	"strings"
	"treescan-go/parser"

	"github.com/antlr4-go/antlr/v4" // Go runtime for ANTLR
	"github.com/robertkrimen/otto"
)

var jsonDump []string

// apexListener is where we inherit our events from, which in turn provide entry points for the scan
type apexListener struct {
	*parser.BaseApexParserListener
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
	var startLine = ctx.GetStart().GetLine()
	var stopLine = ctx.GetStop().GetLine()

	nodeTypeCompact := strings.Replace(nodeType, "*parser.", "", -1)

	vm := otto.New()
	vm.Set("START_LINE", startLine)
	vm.Set("STOP_LINE", stopLine)
	vm.Set("SOURCE", contextSource)
	vm.Set("CONTEXT", ctx)
	vm.Set("NODE_TYPE", nodeTypeCompact)

}

func (a apexListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	//if Args.Debug {
	//	fmt.Println(reflect.TypeOf(ctx).String())
	//}
}
