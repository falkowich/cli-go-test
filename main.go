package main

import (
	"fmt"
	//	"os"

	"cli-go-test/cmd"
	"cli-go-test/internal"
	"github.com/alecthomas/kong"
)

var cli struct {
	Debug bool `help:"Enalbe debug mode."`

	Generate cmd.GenerateCmd `cmd:"" help:"Generate random mab's."`
	Import   cmd.ImportCmd   `cmd:"" help:"Import csv from filename."`
}

func main() {
	ctx := kong.Parse(&cli)
	if cli.Debug {
		fmt.Printf("Context struct: %+v\n", ctx)
		fmt.Printf("CLI struct: %+v\n", cli)
	}
	err := ctx.Run(&internal.Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)

}
