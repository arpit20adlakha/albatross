package helmcli

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"

	"github.com/gojekfarm/albatross/pkg/helmcli/config"
	"github.com/gojekfarm/albatross/pkg/helmcli/flags"
)

type Client interface {
	NewUpgrader(flags.UpgradeFlags) Upgrader
	NewInstaller(flags.InstallFlags) Installer
}

type Upgrader interface {
	Upgrade(ctx context.Context, relName, chartName string, values map[string]interface{}) (*release.Release, error)
}

type Installer interface {
	Install(ctx context.Context, relName, chartName string, values map[string]interface{}) (*release.Release, error)
}

func New() Client {
	return helmClient{}
}

type helmClient struct{}

func (c helmClient) NewUpgrader(flg flags.UpgradeFlags) Upgrader {
	//TODO: ifpossible envconfig could be moved to actionconfig new, remove pointer usage of globalflags
	envconfig := config.NewEnvConfig(&flg.GlobalFlags)
	actionconfig := config.NewActionConfig(envconfig, &flg.GlobalFlags)

	upgrade := action.NewUpgrade(actionconfig.Configuration)
	history := action.NewHistory(actionconfig.Configuration)
	installer := c.NewInstaller(flags.InstallFlags{
		DryRun:      flg.DryRun,
		Version:     flg.Version,
		GlobalFlags: flg.GlobalFlags,
	})

	upgrade.Namespace = flg.Namespace
	upgrade.Install = flg.Install
	upgrade.DryRun = flg.DryRun

	return &upgrader{
		action:      upgrade,
		envSettings: envconfig.EnvSettings,
		history:     history,
		installer:   installer,
	}
}

// NewInstaller returns a new instance of Installer struct
func (c helmClient) NewInstaller(flg flags.InstallFlags) Installer {
	envconfig := config.NewEnvConfig(&flg.GlobalFlags)
	actionconfig := config.NewActionConfig(envconfig, &flg.GlobalFlags)

	install := action.NewInstall(actionconfig.Configuration)
	install.Namespace = flg.Namespace
	install.DryRun = flg.DryRun

	return &installer{
		action:      install,
		envSettings: envconfig.EnvSettings,
	}
}
