package main

import (
	"log"
	"os"
)

type Rules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	nodeType string   `json:"nodeType"`
	scripts  []string `json:"scripts"`
}

func gatherScripts() {

	var scriptPath = `.\scripts\`
	files, err := os.ReadDir(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	//structure and interperetation of computer programs
	for _, file := range files {
		var info = ScriptInfo{}
		info.FileName = scriptPath + file.Name()
		info.NodeType = file.Name()
		ScriptMap[scriptPath+file.Name()] = ScriptInfo{}
	}

}
