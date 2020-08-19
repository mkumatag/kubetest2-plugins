package deployer

import (
	"fmt"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers/common"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers/powervs"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/terraform"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"sigs.k8s.io/kubetest2/pkg/types"
	"sync"
)

const(
	Name = "tf"
)

type deployer struct {
	commonOptions      types.Options
	logsDir            string
	doInit             sync.Once
	tmpDir string
	provider providers.Provider
}

func (d *deployer) init() error {
	var err error
	d.doInit.Do(func() { err = d.initialize() })
	return err
}

func (d *deployer) initialize() error {
	d.provider = powervs.PowerVSProvider
	common.CommonProvider.Initialize()
	d.tmpDir = common.CommonProvider.ClusterName
	if _, err := os.Stat(d.tmpDir); os.IsNotExist(err) {
		err := os.Mkdir(d.tmpDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir: %s", d.tmpDir)
		}
	}else if !ignoreClusterDir{
		return fmt.Errorf("directory named %s already exist, please choose a different cluster-name", d.tmpDir)
	}
	return nil
}

var _ types.Deployer = &deployer{}

var (
	ignoreClusterDir bool
	autoApprove bool
)

func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	d := &deployer{
		commonOptions: opts,
		logsDir:       filepath.Join(opts.ArtifactsDir(), "logs"),
	}
	return d, bindFlags(d)
}

func bindFlags(d *deployer) *pflag.FlagSet {
	flags := pflag.NewFlagSet(Name, pflag.ContinueOnError)
	flags.BoolVar(
		&ignoreClusterDir, "ignore-cluster-dir", false, "Ignore the cluster folder if exists",
	)
	flags.BoolVar(
		&autoApprove, "auto-approve", true, "Terraform Auto Approve",
	)
	flags.MarkHidden("ignore-cluster-dir")
	flags.MarkHidden("auto-approve")
	common.CommonProvider.BindFlags(flags)
	powervs.PowerVSProvider.BindFlags(flags)

	return flags
}

func (d *deployer) Up() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("up failed to init: %s", err)
	}
	err := common.CommonProvider.DumpConfig(d.tmpDir)
	if err != nil {
		return fmt.Errorf("failed to dump common flags: %s", d.tmpDir)
	}
	err = d.provider.DumpConfig(d.tmpDir)
	if err != nil {
		return fmt.Errorf("failed to dumpconfig to: %s and err: %+v", d.tmpDir, err)
	}
	path, err := terraform.Apply(d.tmpDir, "powervs", autoApprove)
	if err != nil {
		return fmt.Errorf("terraform.Apply failed: %v", err)
	}
	fmt.Printf("path: %s", path)
	return nil
}

func (d *deployer) Down() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("down failed to init: %s", err)
	}
	err := terraform.Destroy(d.tmpDir, "powervs", autoApprove)
	if err != nil {
		return fmt.Errorf("terraform.Destroy failed: %v", err)
	}
	return nil
}

func (d *deployer) IsUp() (up bool, err error) {
	panic("implement me")
}

func (d *deployer) DumpClusterLogs() error {
	panic("implement me")
}

func (d *deployer) Build() error {
	panic("implement me")
}
