package main

import (
	"errors"
	"log"
	"os"

	"github.com/tflexsoom/duffle/internal/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "duffle",
		Usage: "duffle tool suite for duffle projects and modules.",
		Commands: []*cli.Command{
			{
				Name:   "compile",
				Usage:  "compile a local duffle project",
				Flags:  compileFlags,
				Action: multiProjectCmd("compile", compileSubCmd),
			},
			{
				Name:   "parse",
				Usage:  "parse a local duffle project",
				Flags:  parseFlags,
				Action: multiProjectCmd("parse", parseSubCmd),
			},
			// {
			// 	Name:   "ir",
			// 	Usage:  "parse a file into an intermediate representation",
			// 	Flags:  baseFlags,
			// 	Action: multiProjectCmd("ir", irSubCmd),
			// },
			{
				Name:   "typecheck",
				Usage:  "typecheck a duffle project",
				Flags:  baseFlags,
				Action: multiProjectCmd("typecheck", typecheckSubCmd),
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
)

func parseSubCmd(cCtx *cli.Context) error {
	return command.ParseOnly(command.ParserOptions{
		ProjectLocations: cCtx.Args().Slice(),
		OutputLocation:   cCtx.Path("output"),
		FunctionOnly:     cCtx.Bool("function"),
		DataOnly:         cCtx.Bool("data"),
		Verbose:          cCtx.Bool("verbose"),
	})
}

// func irSubCmd(cCtx *cli.Context) error {
// 	return command.IntermediateRepresentationOnly(
// 		command.IntermediateRepresentationOptions{
// 			ProjectLocations: cCtx.Args().Slice(),
// 			OutputLocation:   cCtx.Path("output"),
// 			Verbose:          cCtx.Bool("verbose"),
// 		},
// 	)
// }

func typecheckSubCmd(cCtx *cli.Context) error {
	return command.TypeCheckOnly(command.TypeCheckOptions{
		ProjectLocations: cCtx.Args().Slice(),
		OutputLocation:   cCtx.Path("output"),
		Verbose:          cCtx.Bool("verbose"),
	})
}

var compileFlags = append(parseFlags,
	&cli.StringFlag{
		Name:        "backend",
		Aliases:     []string{"B"},
		Usage:       "Backend tool for the output format",
		DefaultText: "binary_x86_64_exe",
	},
)

func compileSubCmd(cCtx *cli.Context) error {
	return command.Compile(command.CompilerOptions{
		ProjectLocations: cCtx.Args().Slice(),
		OutputLocation:   cCtx.Path("output"),
		FunctionOnly:     cCtx.Bool("function"),
		DataOnly:         cCtx.Bool("data"),
		Backend:          cCtx.String("backend"),
		Verbose:          cCtx.Bool("verbose"),
	})
}
