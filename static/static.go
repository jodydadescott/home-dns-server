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
type Domain = types.Domain

const (
	source = "config"
)

type Client struct {
	domain *Domain
}

func New(config *Config) []*Client {

	if config == nil {
		panic("config is required")
	}

	config = config.Clone()

	var clients []*Client

	for _, domain := range config.Domains {

		if domain.Domain == "" {
			domain.Domain = types.DefaultDomain

		}

		clients = append(clients, &Client{domain: domain})
	}

	return clients
}

func (t *Client) GetDomain() (*Domain, error) {

	domainName := t.domain.Domain
	if domainName == "" {
		domainName = types.DefaultDomain
	}

	ptrRecordsMap := make(map[string]*PTRrecord)

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

		arpa, err := util.GetARPA(a.IP)
		if err != nil {
			return nil, err
		}

		a.SRC = source + ":static"

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: a.Hostname,
			Domain:   a.Domain,
			SRC:      source + ":dynamic",
		}

		ptrRecordsMap[p.GetKey()] = p

		// t.domain.AddPtrRecord(p)

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

		r.SRC = source + ":static"

		zap.L().Debug(fmt.Sprintf("Added %s %s CNAME %s", r.SRC, r.GetKey(), r.GetValue()))

	}

	for _, p := range t.domain.PtrRecords {
		existing := ptrRecordsMap[p.GetKey()]
		if existing == nil {

			arpa, err := util.GetARPA(p.ARPA)
			if err != nil {
				return nil, err
			}

			p.ARPA = arpa
			p.SRC = source + ":static"
			ptrRecordsMap[p.GetKey()] = p

			zap.L().Debug(fmt.Sprintf("Added %s %s PTR %s", p.SRC, p.GetKey(), p.GetValue()))
		} else {
			p.SRC = source + ":static-and-dynamic"
			zap.L().Debug(fmt.Sprintf("PTR %s already existed with %s; source upted to %s", p.GetKey(), p.GetValue(), p.SRC))
		}
	}

	var ptrRecords []*PTRrecord

	for _, v := range ptrRecordsMap {
		ptrRecords = append(ptrRecords, v)
	}

	t.domain.PtrRecords = ptrRecords
	return t.domain, nil
}
