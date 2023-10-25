package static

import (
	"fmt"

	"github.com/jodydadescott/home-dns-server/types"
	"github.com/jodydadescott/home-dns-server/util"
	"go.uber.org/zap"
)

type Config = types.StaticConfig
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type InternalUseDomain = types.InternalUseDomain

const (
	source = "static"
)

type Client struct {
	domain *InternalUseDomain
}

func New(config *Config) []*Client {

	if config == nil {
		panic("config is required")
	}

	var clients []*Client

	for _, domain := range config.Domains {

		domainName := domain.Domain
		if domainName == "" {
			domainName = types.DefaultDomain
		}

		clients = append(clients, &Client{
			domain: &InternalUseDomain{
				Domain:       domainName,
				ARecords:     domain.ARecords,
				CnameRecords: domain.CnameRecords,
			},
		})
	}

	return clients
}

func (t *Client) GetDomain() (*InternalUseDomain, error) {

	domainName := t.domain.Domain
	if domainName == "" {
		domainName = types.DefaultDomain
	}

	for _, a := range t.domain.ARecords {
		if a.Hostname == "" {
			return nil, fmt.Errorf("A records must have a hostname")
		}

		if a.IP == "" {
			return nil, fmt.Errorf("A records must have a IP")
		}

		if a.Domain == "" {
			a.Domain = t.domain.Domain
		}

		a.SRC = source

		arpa, err := util.GetARPA(a.IP)
		if err != nil {
			return nil, err
		}

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: a.Hostname,
			Domain:   a.Domain,
			SRC:      "static",
		}

		p.SRC = source

		t.domain.AddPtrRecord(p)

		zap.L().Debug(fmt.Sprintf("Added %s %s A %s and %s PTR %s", a.SRC, a.GetKey(), a.GetValue(), p.GetKey(), p.GetValue()))

	}

	for _, r := range t.domain.CnameRecords {

		if r.AliasHostname == "" {
			return nil, fmt.Errorf("CNAME must have AliasHostname")
		}

		if r.TargetHostname == "" {
			return nil, fmt.Errorf("CNAME must have TargetHostname")
		}

		if r.AliasDomain == "" {
			r.AliasDomain = t.domain.Domain
		}

		if r.TargetDomain == "" {
			r.TargetDomain = t.domain.Domain
		}

		r.SRC = source

		zap.L().Debug(fmt.Sprintf("Added %s %s CNAME %s", r.SRC, r.GetKey(), r.GetValue()))

	}

	return t.domain, nil
}
