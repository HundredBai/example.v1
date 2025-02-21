package certmng

import (
	action2 "github.com/shipengqi/example.v1/apps/cli/internal/action"
	"github.com/shipengqi/example.v1/apps/cli/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createOptions struct {
	caCert        string
	caKey         string
	nodeType      string
	host          string
	serverCertSan string
	outputDir     string
	validity      int
}

func (o *createOptions) combine(f *pflag.FlagSet, cfg *action2.Configuration) {
	if f.Changed(caCertFlagName) {
		cfg.CACert = o.caCert
	}
	if f.Changed(caKeyFlagName) {
		cfg.CAKey = o.caKey
	}
	if f.Changed(hostFlagName) {
		cfg.Host = o.host
	}
	if f.Changed(serverCertSanFlagName) {
		cfg.ServerCertSan = o.serverCertSan
	}
	if f.Changed(outputFlagName) {
		cfg.OutputDir = o.outputDir
	}

	// default value
	cfg.Validity = o.validity
	cfg.NodeType = o.nodeType

}

func newCreateCmd(cfg *action2.Configuration) *cobra.Command {
	o := &createOptions{}
	c := &cobra.Command{
		Use:   "create [options]",
		Short: "Create the internal certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action2.NewCreate(cfg)
			return action2.Execute(c)
		},
	}

	c.Flags().SortFlags = false
	c.DisableFlagsInUseLine = true
	f := c.Flags()
	f.StringVar(&o.caCert, caCertFlagName, "", "CA certificate file path.")
	f.StringVar(&o.caKey, caKeyFlagName, "", "CA key file path.")
	f.StringVar(&o.nodeType, nodeTypeFlagName, types.NodeTypeControlPlane, nodeTypeFlagDesc)
	f.StringVar(&o.host, hostFlagName, "", "The host FQDN or IP address.")
	f.StringVar(&o.serverCertSan, serverCertSanFlagName, "", "server-cert-san for node.")
	f.IntVarP(&o.validity, validityFlagName, "V", 365, validityFlagDesc)

	return c
}
