package unifi

import (
	"fmt"
	"strings"

	"github.com/jodydadescott/home-dns-server/types"
	"github.com/jodydadescott/home-dns-server/util"
	"github.com/jodydadescott/unifi-go-sdk"
	"go.uber.org/zap"
)

type Config = types.UnifiConfig
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type Domain = types.Domain

type Client struct {
	unifiClient *unifi.Client
	config      *Config
	ignoreMacs  []string
	domain      string
}

const (
	source = "unifi"
)

func New(config *Config) *Client {

	if config == nil {
		panic("config is required")
	}

	config = config.Clone()

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
		config:      config,
		domain:      domain,
		unifiClient: unifi.New(&config.Config),
	}
}

func (t *Client) GetDomain() (*Domain, error) {

	clients, err := t.unifiClient.GetClients()
	if err != nil {
		return nil, err
	}

	enrichedConfigs, err := t.unifiClient.GetEnrichedConfiguration()
	if err != nil {
		return nil, err
	}

	domain := &Domain{
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
			SRC:      source + ":unifi-client",
		}

		domain.AddARecords(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: name,
			Domain:   t.domain,
			SRC:      source + ":unifi-client",
		}

		domain.AddPtrRecords(p)

		zap.L().Debug(fmt.Sprintf("Added %s %s A %s and %s PTR %s", a.SRC, a.GetKey(), a.GetValue(), p.GetKey(), p.GetValue()))
	}

	for _, enrichedConfig := range enrichedConfigs {

		name := strings.ToLower(enrichedConfig.Configuration.Name)
		ip := strings.Split(enrichedConfig.Configuration.IPSubnet, "/")[0]

		if name == "" {
			zap.L().Debug("Interface is missing its name")
			continue
		}

		if name == "default" {
			zap.L().Debug("Skipping default interface")
		}

		if ip == "" {
			zap.L().Debug(fmt.Sprintf("Interface %s is missing its ip", enrichedConfig.Configuration.Name))
			continue
		}

		interfaceName := "inf-" + name + "-"
		interfaceName += strings.Replace(ip, ".", "-", -1)

		arpa, err := util.GetARPA(ip)
		if err != nil {
			zap.L().Debug(fmt.Sprintf("Interface %s has an invalid IP; error %s", name, err.Error()))
			continue
		}

		a := &ARecord{
			Hostname: interfaceName,
			Domain:   t.domain,
			IP:       ip,
			SRC:      source + ":unifi-interface",
		}

		domain.AddARecords(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: interfaceName,
			Domain:   t.domain,
			SRC:      source + ":unifi-interface",
		}

		domain.AddPtrRecords(p)

		zap.L().Debug(fmt.Sprintf("Added %s %s A %s and %s PTR %s", a.SRC, a.GetKey(), a.GetValue(), p.GetKey(), p.GetValue()))
	}

	devices, err := t.unifiClient.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices.NetworkDevices {

		if device.Name == "" {
			zap.L().Debug("Device is missing its name")
			continue
		}

		if device.IP == "" {
			zap.L().Debug(fmt.Sprintf("Interface %s is missing its ip", device.Name))
			continue
		}

		if t.config.IgnoreMac(device.Mac) {
			zap.L().Debug(fmt.Sprintf("Ignoring mac %s with name %s", device.Mac, device.Name))
			continue
		}

		arpa, err := util.GetARPA(device.IP)
		if err != nil {
			zap.L().Debug(fmt.Sprintf("Device %s has an invalid IP; error %s", device.Name, err.Error()))
			continue
		}

		a := &ARecord{
			Hostname: device.Name,
			Domain:   t.domain,
			IP:       device.IP,
			SRC:      source + ":unifi-device",
		}

		domain.AddARecords(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: device.Name,
			Domain:   t.domain,
			SRC:      source + ":unifi-device",
		}

		domain.AddPtrRecords(p)

		zap.L().Debug(fmt.Sprintf("Added %s %s A %s and %s PTR %s", a.SRC, a.GetKey(), a.GetValue(), p.GetKey(), p.GetValue()))

	}

	return domain, nil

}
