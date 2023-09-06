package login

import (
	"github.com/taubyte/tau-cli/cli/common/options"
	"github.com/taubyte/tau-cli/flags"
	loginFlags "github.com/taubyte/tau-cli/flags/login"
	"github.com/taubyte/tau-cli/i18n"
	loginI18n "github.com/taubyte/tau-cli/i18n/login"
	loginLib "github.com/taubyte/tau-cli/lib/login"
	"github.com/taubyte/tau-cli/prompts"
	loginPrompts "github.com/taubyte/tau-cli/prompts/login"
	slices "github.com/taubyte/utils/slices/string"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name: "login",
	Flags: flags.Combine(
		flags.Name,
		loginFlags.Token,
		loginFlags.Provider,
		loginFlags.New,
		loginFlags.SetDefault,
	),
	ArgsUsage: i18n.ArgsUsageName,
	Action:    Run,
	Before:    options.SetNameAsArgs0,
}

func Run(ctx *cli.Context) error {
	// get available profiles after reading the config (will
	// autoload a new config if a config has not yet been set)
	// (options var should be renamed to something like nameOpts
	// to avoid variable/package name collisions)
	_default, options, err := loginLib.GetProfiles()
	if err != nil {
		return loginI18n.GetProfilesFailed(err)
	}

	// New: if --new or no selectable profiles
	if ctx.Bool(loginFlags.New.Name) || len(options) == 0 {
		return New(ctx, options)
	}

	// Selection
	var name string
	// if a name (alias -n) flag
	// was provided by the user, set
	// the name variable to the value
	// set by that flag if it exists in
	// the options returned on line 33,
	// or error if it doesn't exist
	if ctx.IsSet(flags.Name.Name) {
		name = ctx.String(flags.Name.Name)

		if !slices.Contains(options, name) {
			return loginI18n.DoesNotExistIn(name, options)
		}
		// if a name (alias -n) flag was
		// not provided by the user, prompt them
		// to select a profile from the options
		// returned on line 33 to use
	} else {
		name, err = prompts.SelectInterface(options, loginPrompts.SelectAProfile, _default)
		if err != nil {
			return err
		}
	}

	// select the chosen profile, overriding the default name
	// if the user provided the set-default (alias -d) flag
	return Select(ctx, name, ctx.Bool(loginFlags.SetDefault.Name))
}
