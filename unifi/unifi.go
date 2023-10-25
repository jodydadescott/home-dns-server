package unifi

import (
	"fmt"

	"github.com/jodydadescott/home-dns-server/types"
	"github.com/jodydadescott/home-dns-server/util"
	"github.com/jodydadescott/unifi-go-sdk/unifi"
	"go.uber.org/zap"
)

type Config = types.UnifiConfig
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type InternalUseDomain = types.InternalUseDomain

type Client struct {
	unifiClient *unifi.Client
	domain      string
}

const (
	source = "unifi"
)

func New(config *Config) *Client {

	if config == nil {
		panic("config is required")
	}

	if config.Hostname == "" {
		panic("Hostname is required")
	}

	if config.Username == "" {
		panic("Username is required")
	}

	if config.Password == "" {
		panic("Password is required")
	}

	domain := types.DefaultDomain
	if config.Domain != "" {
		domain = config.Domain
	}

	return &Client{
		domain:      domain,
		unifiClient: unifi.New(&config.Config),
	}
}

func (t *Client) GetDomain() (*InternalUseDomain, error) {

	clients, err := t.unifiClient.GetClients()
	if err != nil {
		return nil, err
	}

	domain := &InternalUseDomain{
		Domain: t.domain,
	}

	for _, client := range clients {

		name := client.Name

		if name == "" {
			name = client.Hostname
		}

		if name == "" {
			zap.L().Debug(fmt.Sprintf("Client with MAC=%s and IP=%s does not have a name", client.Mac, client.IP))
			continue
		}

		if client.IP == "" {
			zap.L().Debug(fmt.Sprintf("Client %s does not have an IP", name))
			continue
		}

		arpa, err := util.GetARPA(client.IP)
		if err != nil {
			zap.L().Debug(fmt.Sprintf("Client %s has an invalid IP; error %s", name, err.Error()))
			continue
		}

		a := &ARecord{
			Hostname: name,
			Domain:   t.domain,
			IP:       client.IP,
			SRC:      source,
		}

		domain.AddARecord(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: name,
			Domain:   t.domain,
			SRC:      source,
		}

		domain.AddPtrRecord(p)

		zap.L().Debug(fmt.Sprintf("Added %s %s A %s and %s PTR %s", a.SRC, a.GetKey(), a.GetValue(), p.GetKey(), p.GetValue()))
	}

	return domain, nil

}
