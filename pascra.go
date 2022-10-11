package pascra

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/z89/pascra/cmd"
)

var VERSION = "1.0.0"

type Pascra struct {
	RootCmd *cobra.Command
}

func (pascra *Pascra) Start() {
	pascra.RootCmd.AddCommand(cmd.Fetch())

	if err := pascra.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func New() *Pascra {
	return &Pascra{
		RootCmd: &cobra.Command{
			Use:     "pascra",
			Short:   "pascra - a golang pastebin.com web scraper",
			Version: VERSION,

			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}
}
