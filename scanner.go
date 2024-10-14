package main

import "strings"

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
}
