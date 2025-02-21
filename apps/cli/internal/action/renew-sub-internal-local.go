package action

import (
	"github.com/shipengqi/example.v1/apps/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/apps/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/apps/cli/internal/sysc"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"strings"
)

type renewSubInternalLocal struct {
	*action

	generator certs.Generator
}

func NewRenewSubInternalLocal(cfg *Configuration) Interface {
	c := &renewSubInternalLocal{
		action: newAction("renew-sub-internal-local", cfg),
	}

	key, err := c.parseCAKey()
	if err != nil {
		panic(err)
	}

	g, err := infra.New(cfg.CACert, key)
	if err != nil {
		panic(err)
	}
	c.generator = g

	return c
}

func (a *renewSubInternalLocal) Name() string {
	return a.name
}

func (a *renewSubInternalLocal) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	if a.cfg.Env.RunOnMaster {
		log.Debug("renew certificate secrets on master node")
		err := a.iterateSecrets(a.generator)
		if err != nil {
			return err
		}
	}

	return a.iterate(a.cfg.Host, true, true, a.generator)
}

func (a *renewSubInternalLocal) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	hostname, err := sysc.Hostname()
	if err != nil {
		log.Warnf("sysc.Hostname(): %v", err)
	} else {
		a.cfg.Host = hostname
	}
	log.Debugf("get local hostname: %s", hostname)

	a.cfg.Debug()

	return nil
}
