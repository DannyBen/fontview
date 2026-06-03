package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"fontview/report"
)

type options struct {
	addr      string
	output    string
	writeHTML bool
	version   string
}

func Execute(args []string, version string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help":
			usage(os.Stdout)
			return nil
		case "--version":
			fmt.Println(version)
			return nil
		}
	}

	opts := options{version: version}
	flags := flag.NewFlagSet("fontview", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.addr, "addr", "0.0.0.0:3000", "server address")
	flags.BoolVar(&opts.writeHTML, "html", false, "write HTML instead of starting a server")
	flags.StringVar(&opts.output, "o", "fontview.html", "HTML output path")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w\n\n%s", err, usageText())
	}

	return report.Run(report.Options{
		Inputs:    flags.Args(),
		Addr:      opts.addr,
		Output:    opts.output,
		WriteHTML: opts.writeHTML,
		Version:   opts.version,
	})
}

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "fontview: %v\n", err)
}

func usage(w io.Writer) {
	fmt.Fprint(w, usageText())
}

func usageText() string {
	return `Usage:
  fontview [options] [FONT_OR_DIR...]

By default, fontview discovers font files in the current folder and starts a
local browser-friendly server. Pass one or more font files or directories to
limit the report.

Options:
  --addr ADDR   Server address (default 0.0.0.0:3000)
  --html        Write a standalone HTML report instead of serving
  -o PATH       HTML output path when using --html (default fontview.html)
  --version     Print version
  -h, --help    Show help
`
}
