# treescan-go

This is the brother-from-another-mother of treescan-rs...which I may abandon in favor of this.

It is a static analysis tool for Salesforce source code. Rules are written in JavaScript, the engine itself is this thing.

This is a rough skeleton of how things should work. Very much a work in progress, assume there be dragons everywhere.

### Progress/Notes/Cautions

* One thing I'm considering...when a given rule is entered (a node, if you will) I can allow the JavaScript VM's state to persist across that particular enter event for all the scripts associated with that context. 
* Got some more stuff working, some stuff I need to figure out.
  * Rules are now being loaded via the config json. 
  * It's probably my Go ignorance: in the EnterEveryRUle event, I can't seem to get the values for the JavaScript source unless I loop through the map. ** That can't be right. Feedback appreciated **

* The build is broken currently, but not by much. Tinkering with different ways to handle rule files. My idea is as follows:
  * One JSON file (see scripts directory) that specifies the rule filenames and the node context they scan (method invocations, class declarations, etc.)
  * Each rule is a separate .js script file in the same directory as the above config file.
  * Why? To keep the process lean. Don't do any processing if you don't have to.

### Current Issues

This is just stuff I know about and will be fixing soon:
* The ANTLR4 grammar doesn't like SOQL and SOSL. I'll get to the bottom of it. It doesn't break things, but it causes the parser to get confused.
* Everything is very much in an embryonic, ad-hoc state. I have proven (to myself, at least) that:
    * We can scan apex. Lots of it.
    * We can execute JavaScript against contextual information about that apex whenever an ANTLR 'rule' is entered.
    * This JavaScript has one instance per rule entrance. There is no multithreading yet.
    * Go is a perfectly adequate language for this.
    * **For a certain repository, running this CLI tool against the entire classes directory takes 21 seconds. Wow. This is with zero optimization.**

### Can I run this?

Go for it. ~~Right now only one of the arguments is actually supported, `-s` or `--scan`. Right now it only accepts one path and doesn't recurse.~~ `-d` or `--debug` also works.

### How Will it Look Eventually?

Here's the help text:
```text
Usage: treescan-go.exe [--sourcepaths SOURCEPATHS] [--reportpath REPORTPATH] [--dump] [--dumpformat DUMPFORMAT] [--reportformat REPORTFORMAT] [--debug] [--languages LANGUAGES]

Options:
  --sourcepaths SOURCEPATHS, -s SOURCEPATHS
                         Comma seperated list of file paths that will be scanned.
  --reportpath REPORTPATH, -o REPORTPATH
                         Where the reported scan results should be stored (aside from STDIO) [env: TSGO_REPORT_PATH]
  --dump, -d             Dump all results to stdout instead of scanning
  --dumpformat DUMPFORMAT, -f DUMPFORMAT
                         Format of dump command
  --reportformat REPORTFORMAT, -r REPORTFORMAT
                         Format of report command
  --debug, -x            Enable debug mode
  --languages LANGUAGES, -l LANGUAGES
                         Comma separated list of languages
  --help, -h             display this help and exit
```
