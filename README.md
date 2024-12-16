# treescan-go

## Update

Been experimenting with using tree sitter and TypeScript as a solution but it's difficult to see past the performance considerations. While Rust is nice and will arguably find itself dominating the Linux kernel at some point, I don't see that happening any time in the near-near future. Go, on the other hand, has a proven record of satisfying results within the devops world (beyond Salesforce, which is important to think about if there is a need for quality Salesforce implementations. Developers should stop thinking as admins.)

I have an immediate ask for something that can provide metrics (you have `[count of all variables]` variables in your codebase, `[count per analysis rule]` violate the `[the rule we're interested in]` rule.) This fits the bill.

I've already tested this against a large codebase, so now the trick is that there needs to be a way to define quick-hit rules and render useful results.

This is a rough skeleton of how things should work. Very much a work in progress, assume there be dragons everywhere.

### Can I run this?

Go for it. [Don't expect to solve the issues that faced the Brazilian economy in 1994](https://en.mercopress.com/2014/07/02/two-decades-of-the-real-the-currency-that-helped-brazil-trust-financial-stability) or anything. ~~Right now only one of the arguments is actually supported, `-s` or `--scan`. Right now it only accepts one path and doesn't recurse.~~ `-d` or `--debug` also works.

### How Will it Look Eventually?

Here's the help text:
```text
Usage: treescan-go.exe [--sourcepaths SOURCEPATHS] [--reportpath REPORTPATH] [--dump] [--dumpformat DUMPFORMAT] [--reportformat REPORTFORMAT] [--debug] [--languages LANGUAGES]

Options:
  --sourcepaths SOURCEPATHS, -s SOURCEPATHS
                         Comma seperated list of file paths that will be scanned.
  --reportpath REPORTPATH, -o REPORTPATH
                         Where the reported scan results should be stored (aside from STDIO) [env: TSGO_REPORT_PATH]
  --reportformat REPORTFORMAT, -r REPORTFORMAT
                         Format of report command
  --debug, -x            Enable debug mode
  --languages LANGUAGES, -l LANGUAGES
                         Comma separated list of languages
  --help, -h             display this help and exit

Commands:
  scan [options]        Initiate a static analysis scan with the supplied options
  measure [options]     Get metrics for the supplied codebase that is friendly for deeper analysis
```
