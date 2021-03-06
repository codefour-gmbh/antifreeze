package application

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/api/applications"
	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/configuration/coreconfig"
	"github.com/cloudfoundry/cli/cf/errors"
	. "github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/flags"
)

type DeleteApp struct {
	ui        terminal.UI
	config    coreconfig.Reader
	appRepo   applications.ApplicationRepository
	routeRepo api.RouteRepository
	appReq    requirements.ApplicationRequirement
}

func init() {
	commandregistry.Register(&DeleteApp{})
}

func (cmd *DeleteApp) MetaData() commandregistry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}
	fs["r"] = &flags.BoolFlag{ShortName: "r", Usage: T("Also delete any mapped routes")}

	return commandregistry.CommandMetadata{
		Name:        "delete",
		ShortName:   "d",
		Description: T("Delete an app"),
		Usage: []string{
			T("CF_NAME delete APP_NAME [-f -r]"),
		},
		Flags: fs,
	}
}

func (cmd *DeleteApp) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires app name as argument\n\n") + commandregistry.Commands.CommandUsage("delete"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
	}

	return reqs
}

func (cmd *DeleteApp) SetDependency(deps commandregistry.Dependency, pluginCall bool) commandregistry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.appRepo = deps.RepoLocator.GetApplicationRepository()
	cmd.routeRepo = deps.RepoLocator.GetRouteRepository()
	return cmd
}

func (cmd *DeleteApp) Execute(c flags.FlagContext) {
	appName := c.Args()[0]

	if !c.Bool("f") {
		response := cmd.ui.ConfirmDelete(T("app"), appName)
		if !response {
			return
		}
	}

	cmd.ui.Say(T("Deleting app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{
			"AppName":   terminal.EntityNameColor(appName),
			"OrgName":   terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"SpaceName": terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	app, apiErr := cmd.appRepo.Read(appName)

	switch apiErr.(type) {
	case nil: // no error
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(T("App {{.AppName}} does not exist.", map[string]interface{}{"AppName": appName}))
		return
	default:
		cmd.ui.Failed(apiErr.Error())
	}

	if c.Bool("r") {
		for _, route := range app.Routes {
			apiErr = cmd.routeRepo.Delete(route.Guid)
			if apiErr != nil {
				cmd.ui.Failed(apiErr.Error())
			}
		}
	}

	apiErr = cmd.appRepo.Delete(app.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
	}

	cmd.ui.Ok()
}
