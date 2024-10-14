package main

import (
	"log"
	"os"
	"strings"
	"treescan-go/parser"

	"github.com/antlr4-go/antlr/v4"
)

func (config *RuleConfiguration) GetRulesForNode(nodeType string) []Rule {
	var result []Rule
	for ruleIndex := range len(CurrentConfiguration) {
		thisRule := CurrentConfiguration[ruleIndex]
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

	// Apex is NOT case sensitive, so that presents a problem...especially with SELECT and standard query conventions.
	// Maybe use a regex for this in the long term
	sourceString = strings.ToLower(sourceString)

	sourceInputStream := antlr.NewInputStream(sourceString)
	lexer := parser.NewApexLexer(sourceInputStream)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewApexParser(stream)
	parser.ApexParserInit()
	p.BuildParseTrees = true

	antlr.ParseTreeWalkerDefault.Walk(&apexListener{}, p.CompilationUnit())

}
