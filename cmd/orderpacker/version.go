package main

import (
	"fmt"
	log "log/slog"
	"os"
	"text/tabwriter"

	"github.com/obalunenko/version"
)

func printVersion() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	_, err := fmt.Fprintf(w, `
| app_name:	%s	|
| version:	%s	|
| short_commit:	%s	|
| commit:	%s	|
| build_date:	%s	|
| goversion:	%s	|
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
`,
		version.GetAppName(),
		version.GetVersion(),
		version.GetShortCommit(),
		version.GetCommit(),
		version.GetBuildDate(),
		version.GetGoVersion())
	if err != nil {
		log.Error("Print version error: %v", err)

		return
	}

	if err = w.Flush(); err != nil {
		log.Error("Flush error: %v", err)

		return
	}
}
