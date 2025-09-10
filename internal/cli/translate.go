package cli

import (
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var translateCommand = cli.Command{
	Name:      "translate-templates",
	ShortName: "tt",
	Usage:     "translate policy/procedure/narrative templates",
	ArgsUsage: "[path]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "provider, p",
			Usage: "LLM provider (openai, anthropic, ollama)",
		},
	},
	Action: translateTemplatesAction,
}

func translateTemplatesAction(c *cli.Context) error {
	path := c.Args().First()
	provider := c.String("provider")

	err := render.TranslateTemplates(path, provider)
	if err != nil {
		return errors.Wrap(err, "template translation failed")
	}
	return nil
}
