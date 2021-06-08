package action

import (
	"fmt"
	"net"
	"strings"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ConfigMapNameCDFCluster = "cdf-cluster-host"
	ConfigMapNameCDF        = "cdf"
)

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Execute() error
}

type action struct {
	name string
	cfg  *Configuration
}

func (a *action) Name() string {
	return "[action]"
}

func (a *action) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	a.cfg.Debug()
	return nil
}

func (a *action) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}

func (a *action) iterate(address string, master bool, generator certs.Generator) error {
	dnsnames, ipaddrs, cn := a.combineSubject(address, master)
	for _, v := range CertificateSet {
		if !v.CanDep(master) {
			continue
		}
		dns := make([]string, 0)
		ips := make([]net.IP, 0)
		copy(dns, dnsnames)
		copy(ips, ipaddrs)

		v.CombineServerSan(dns, ips, cn, a.cfg.ServerCertSan, a.cfg.Cluster.KubeServiceIP)
		v.Validity = a.cfg.Validity
		v.UintTime = a.cfg.Unit

		log.Debugf("cert DNS: %s", v.DNSNames)
		log.Debugf("cert IPs: %s", v.IPs)
		log.Debugf("cert CN: %s", v.CN)

		err := generator.GenAndDump(v.Certificate, a.cfg.OutputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *action) combineSubject(address string, master bool) ([]string, []net.IP, string) {
	cn := address
	var dnsNames []string
	var ips []net.IP

	if !master {
		if utils.IsIPV4(address) {
			cn = fmt.Sprintf("host-%s", address)
		}

		return dnsNames, ips, cn
	}

	if a.cfg.Cluster.VirtualIP != "" {
		ips = append(ips, net.ParseIP(a.cfg.Cluster.VirtualIP))
	}
	if a.cfg.Cluster.LoadBalanceIP != "" {
		LBIP := net.ParseIP(a.cfg.Cluster.LoadBalanceIP)
		if LBIP != nil {
			ips = append(ips, LBIP)
		} else {
			dnsNames = append(dnsNames, a.cfg.Cluster.LoadBalanceIP)
		}
	}
	if utils.IsIPV4(address) {
		cn = fmt.Sprintf("host-%s", address)
		ips = append(ips, net.ParseIP(address))
	} else {
		dnsNames = append(dnsNames, address)
		shortNames := strings.Split(address, ".")
		// Add first short name
		if len(shortNames) > 0 {
			dnsNames = append(dnsNames, shortNames[0])
		}

		// Add master node ip
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP)
						break
					}
				}
			}
		}
	}

	return dnsNames, ips, cn
}

// ----------------------------------------------------------------------------
// Helpers...

func parseSan(san string) ([]string, []net.IP, net.IP) {
	if len(san) == 0 {
		return nil, nil, nil
	}

	var svcIp net.IP
	dns := make([]string, 0)
	ips := make([]net.IP, 0)

	subs := strings.Split(san, ",")
	for _, sub := range subs {
		if strings.HasPrefix(sub, "DNS:") {
			dns = append(dns, sub[4:])
		}
		if strings.HasPrefix(sub, "IP:") {
			ips = append(ips, net.ParseIP(sub[3:]))
		}
		if strings.HasPrefix(sub, "K8SSVCIP:") {
			svcIp = net.ParseIP(sub[9:])
		}
	}

	return dns, ips, svcIp
}
