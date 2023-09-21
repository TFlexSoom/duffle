package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/tflexsoom/deflemma/internal/compile"
)

func main() {
	app := &cli.App{
		Name:  "deflemma",
		Usage: "deflemma tool suite for deflemma projects and modules.",
		Commands: []*cli.Command{
			{
				Name:   "parse",
				Usage:  "parse a local deflemma project",
				Flags:  parseFlags,
				Action: multiProjectCmd("parse", parseSubCmd),
			},
			{
				Name:   "compile",
				Usage:  "compile a local deflemma project",
				Flags:  compileFlags,
				Action: multiProjectCmd("compile", compileSubCmd),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var baseFlags = []cli.Flag{
	&cli.PathFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Output file pathname for the backend output",
		Value:   "./a.out",
	},
	&cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "Print out debug information while performing work",
		Value:   false,
	},
}

func multiProjectCmd(
	cmdName string,
	cmdImpl func(cCtx *cli.Context) error,
) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		if cCtx.Args().Len() < 1 {
			return errors.New("Missing Project Destination! \"COMMAND [command options] [arguments...]\"")
		}

		return cmdImpl(cCtx)
	}
}

var parseFlags = append(baseFlags,
	&cli.BoolFlag{
		Name:  "function",
		Usage: "Only Parse Function Files",
		Value: false,
	},
	&cli.BoolFlag{
		Name:  "data",
		Usage: "Only Parse Data Files",
		Value: false,
	},
	&cli.BoolFlag{
		Name:  "memory",
		Usage: "Only Parse Memory Structure Files",
		Value: false,
	},
)

func parseSubCmd(cCtx *cli.Context) error {
	return compile.ParseOnly(compile.ParserOptions{
		ProjectLocations: cCtx.Args().Slice(),
		OutputLocation:   cCtx.Path("output"),
		FunctionOnly:     cCtx.Bool("function"),
		DataOnly:         cCtx.Bool("data"),
		MemoryOnly:       cCtx.Bool("memory"),
		Verbose:          cCtx.Bool("verbose"),
	})
}

var compileFlags = append(baseFlags,
	&cli.StringFlag{
		Name:        "backend",
		Aliases:     []string{"B"},
		Usage:       "Backend tool for the output format",
		DefaultText: "binary_x86_64_exe",
	},
)

func compileSubCmd(cCtx *cli.Context) error {
	return compile.Compile(compile.CompilerOptions{
		ProjectLocations: cCtx.Args().Slice(),
		OutputLocation:   cCtx.Path("output"),
		Backend:          cCtx.String("backend"),
		Verbose:          cCtx.Bool("verbose"),
	})
}
