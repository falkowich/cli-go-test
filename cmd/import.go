package cmd

import (
	"cli-go-test/internal"
	"fmt"
)

type ImportCmd struct {
	Filename string `help:"Filename of cs tom import from"`
}

func (i *ImportCmd) Run(ctx *internal.Context) error {
	fmt.Println("import", i.Filename)
	return nil
}
