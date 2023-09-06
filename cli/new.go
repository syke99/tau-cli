package cli

import (
	"github.com/pterm/pterm"
	"github.com/taubyte/tau-cli/cli/commands/autocomplete"
	"github.com/taubyte/tau-cli/cli/commands/current"
	"github.com/taubyte/tau-cli/cli/commands/dream"
	"github.com/taubyte/tau-cli/cli/commands/exit"
	"github.com/taubyte/tau-cli/cli/commands/login"
	"github.com/taubyte/tau-cli/cli/commands/resources/application"
	"github.com/taubyte/tau-cli/cli/commands/resources/builds"
	"github.com/taubyte/tau-cli/cli/commands/resources/builds/build"
	"github.com/taubyte/tau-cli/cli/commands/resources/database"
	"github.com/taubyte/tau-cli/cli/commands/resources/domain"
	"github.com/taubyte/tau-cli/cli/commands/resources/function"
	"github.com/taubyte/tau-cli/cli/commands/resources/library"
	"github.com/taubyte/tau-cli/cli/commands/resources/logs"
	"github.com/taubyte/tau-cli/cli/commands/resources/messaging"
	"github.com/taubyte/tau-cli/cli/commands/resources/network"
	"github.com/taubyte/tau-cli/cli/commands/resources/project"
	"github.com/taubyte/tau-cli/cli/commands/resources/service"
	"github.com/taubyte/tau-cli/cli/commands/resources/smartops"
	"github.com/taubyte/tau-cli/cli/commands/resources/storage"
	"github.com/taubyte/tau-cli/cli/commands/resources/website"
	"github.com/taubyte/tau-cli/cli/commands/version"
	"github.com/taubyte/tau-cli/cli/common"
	"github.com/taubyte/tau-cli/flags"
	"github.com/taubyte/tau-cli/states"
	"github.com/urfave/cli/v2"
)

func New() (*cli.App, error) {
	// configure the Env and Color flags to
	// read from the TAU_ENV and TAU_COLOR
	// environment variables, respectively
	globalFlags := []cli.Flag{
		flags.Env,
		flags.Color,
	}

	app := &cli.App{
		UseShortOptionHandling: true,
		Flags:                  globalFlags,
		EnableBashCompletion:   true,
		Before: func(ctx *cli.Context) error {
			// pass ctx to states.New() to
			// create a new internal context
			// with a cancel func (set to Context
			// and ContextC, respectively, inside
			// ../states/context.go) to guarantee
			// there isn't a nil context internally
			// and allow for internal context cancellation
			states.New(ctx.Context)

			// check to see if a (terminal) color
			// has been set in the context (via reading
			// from the loaded TAU_COLOR env var above),
			// and if so, guarantee that the value is one
			// of the valid options
			color, err := flags.GetColor(ctx)
			if err != nil {
				// return err if not nil to prevent
				// unnecessary completion of App setup
				return err
			}

			// disable terminal color(s) if Color flag
			// was set to "never" (why a switch and not
			// an if statement? possibly future options?)
			switch color {
			case flags.ColorNever:
				pterm.DisableColor()
			}

			// return nil if App setup successful
			return nil
		},
		// add login, current, exit, and dream
		// Commands to App
		Commands: []*cli.Command{
			login.Command,
			current.Command,
			exit.Command,
			dream.Command,
		},
	}

	// attach the below Commands as SubCommands
	// to all the base commands found in
	// tau-cli/cli/common/base_commands.go
	common.Attach(app,
		project.New,
		application.New,
		network.New,
		database.New,
		domain.New,
		function.New,
		library.New,
		messaging.New,
		service.New,
		smartops.New,
		storage.New,
		website.New,
		builds.New,
		build.New,
		logs.New,
	)

	// add the autocomplete and version commands to the app
	//
	// ( curious as to why these aren't added with the rest of
	// the commands on line 78? adding these two commands here creates
	// unnecessary memory (the slice of *cli.Command is an extra
	// slice, not needed if added above, and if the array backing
	// app.Commands doesn't have the capacity, it will have to be
	// re-allocated and re-sliced, creating more memory for the GC
	// to have to clean up) unless command ordering is of importance)
	app.Commands = append(app.Commands, []*cli.Command{
		autocomplete.Command,
		version.Command,
	}...)

	return app, nil
}
