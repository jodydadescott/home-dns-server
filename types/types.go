package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/jodydadescott/unifi-go-sdk"
)

// ProtoType is the protocol type. Currently only UDP and TCP are supported.
type ProtoType string

const (
	ProtoTypeInvalid ProtoType = "INVALID"
	ProtoTypeUDP               = "udp"
	ProtoTypeTCP               = "tcp"
)

const (
	DefaultDomain = "home"
)

var space = regexp.MustCompile(`\s+`)

// String returns string value of the protocol type
func (t ProtoType) String() string {

	switch t {

	case ProtoTypeUDP:
		return string(t)

	case ProtoTypeTCP:
		return string(t)

	}

	panic("Invalid proto type")
}

// Netport is the IP, Port and Protocol type
type NetPort struct {
	IP          string    `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port        int       `json:"port,omitempty" yaml:"port,omitempty"`
	Proto       ProtoType `json:"proto,omitempty" yaml:"proto,omitempty"`
	ipColonPort string    `json:"-"`
}

// Clone return copy
func (t *NetPort) Clone() *NetPort {
	c := &NetPort{}
	copier.Copy(&c, &t)
	return c
}

// SetProtoTypeTCP sets proto type to TCP
func (t *NetPort) SetProtoTypeTCP() {
	t.Proto = ProtoTypeTCP
}

// SetProtoTypeTCP sets proto type to UDP
func (t *NetPort) SetProtoTypeUDP() {
	t.Proto = ProtoTypeUDP
}

// GetIPColonPort returns the IP + colong + port as a string
func (t *NetPort) GetIPColonPort() string {
	if t.ipColonPort == "" {
		t.ipColonPort = t.IP + ":" + fmt.Sprint(t.Port)
	}
	return t.ipColonPort
}

// ARecord is a DNS A Record
type ARecord struct {
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty"`
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IP       string `json:"ip,omitempty" yaml:"ip,omitempty"`
	SRC      string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdn     string `json:"-"`
}

// Clone return copy
func (t *ARecord) Clone() *ARecord {
	c := &ARecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *ARecord) GetKey() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

// GetValue returns the value for the record type
func (t *ARecord) GetValue() string {
	return t.IP
}

// CNameRecord is a DNS CNAME Record
type CNameRecord struct {
	AliasHostname  string `json:"aliasHostname,omitempty" yaml:"aliasHostname,omitempty"`
	AliasDomain    string `json:"aliasDomain,omitempty" yaml:"aliasDomain,omitempty"`
	TargetHostname string `json:"targetHostname,omitempty" yaml:"targetHostname,omitempty"`
	TargetDomain   string `json:"targetDomain,omitempty" yaml:"targetDomain,omitempty"`
	SRC            string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdnAlias      string `json:"-"`
	fqdnTarget     string `json:"-"`
}

// Clone return copy
func (t *CNameRecord) Clone() *CNameRecord {
	c := &CNameRecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *CNameRecord) GetKey() string {
	if t.fqdnAlias == "" {
		t.fqdnAlias = cleanHostname(t.AliasHostname) + "." + t.AliasDomain + "."
	}
	return t.fqdnAlias
}

// GetValue returns the value for the record type
func (t *CNameRecord) GetValue() string {
	if t.fqdnTarget == "" {
		t.fqdnTarget = cleanHostname(t.TargetHostname) + "." + t.TargetDomain + "."
	}
	return t.fqdnTarget
}

// PTRrecord is a DNS PTR Record
type PTRrecord struct {
	ARPA     string `json:"arpa,omitempty" yaml:"arpa,omitempty"`
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty"`
	SRC      string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdn     string `json:"-"`
}

// Clone return copy
func (t *PTRrecord) Clone() *PTRrecord {
	c := &PTRrecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *PTRrecord) GetKey() string {
	return t.ARPA
}

// GetValue returns the value for the record type
func (t *PTRrecord) GetValue() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

// Config is the main user level config
type Config struct {
	Notes       string        `json:"notes,omitempty" yaml:"notes,omitempty"`
	Unifi       *UnifiConfig  `json:"unifiConfig,omitempty" yaml:"unifiConfig,omitempty"`
	Listeners   []*NetPort    `json:"listeners,omitempty" yaml:"listeners,omitempty"`
	Static      *StaticConfig `json:"static,omitempty" yaml:"static,omitempty"`
	Nameservers []*NetPort    `json:"nameservers,omitempty" yaml:"nameservers,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

// AddNameserver adds the specified nameserver to the config
func (t *Config) AddNameservers(nameservers ...*NetPort) *Config {
	for _, v := range nameservers {
		t.Nameservers = append(t.Nameservers, v)
	}
	return t
}

// AddNameserver adds the specified nameserver to the config
func (t *Config) AddListeners(listeners ...*NetPort) *Config {
	for _, v := range listeners {
		t.Listeners = append(t.Listeners, v)
	}
	return t
}

// UnifiConfig is the config for Unifi servers
type UnifiConfig struct {
	unifi.Config
	Enabled    bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domain     string   `json:"domain,omitempty" yaml:"domain,omitempty"`
	IgnoreMacs []string `json:"ignoreMacs,omitempty" yaml:"ignoreMacs,omitempty"`
}

// Clone return copy
func (t *UnifiConfig) Clone() *UnifiConfig {
	c := &UnifiConfig{}
	copier.Copy(&c, &t)
	return c
}

// AddIgnoreMac add a MAC that will be ignored when process the Unifi config
func (t *UnifiConfig) AddIgnoreMacs(macs ...string) *UnifiConfig {
	for _, v := range macs {
		t.IgnoreMacs = append(t.IgnoreMacs, v)
	}
	return t
}

// IgnoreMac is a convenience function that returns true if the MAC exist in the ignore slice
func (t *UnifiConfig) IgnoreMac(mac string) bool {
	mac = strings.ToLower(mac)
	for _, m := range t.IgnoreMacs {
		if mac == strings.ToLower(m) {
			return true
		}
	}
	return false
}

// StaticConfig are records from config that are statically defined
type StaticConfig struct {
	Enabled bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domains []*Domain `json:"domains,omitempty" yaml:"domains,omitempty"`
}

// Clone return copy
func (t *StaticConfig) Clone() *StaticConfig {
	c := &StaticConfig{}
	copier.Copy(&c, &t)
	return c
}

// Domain is a collection of A & CNAME records with a common domain. If the domain is
// not set then a default domain will be used. The same default domain will be used
// of CNAME target domains if not configured. It is not normally required to add PTR
// records as they will be automatically generated when the A record is created.
type Domain struct {
	Domain       string         `json:"domain,omitempty" yaml:"dnsDomain,omitempty"`
	ARecords     []*ARecord     `json:"aRecords,omitempty" yaml:"aRecords,omitempty"`
	CnameRecords []*CNameRecord `json:"cnameRecords,omitempty" yaml:"cnameRecords,omitempty"`
	PtrRecords   []*PTRrecord   `json:"ptrRecords,omitempty" yaml:"ptrRecords,omitempty"`
}

// Clone return copy
func (t *Domain) Clone() *Domain {
	c := &Domain{}
	copier.Copy(&c, &t)
	return c
}

// AddDomain is a convenience function that adds the specified Domaain to the StaticConfig
func (t *StaticConfig) AddDomains(domains ...*Domain) *StaticConfig {
	for _, v := range domains {
		t.Domains = append(t.Domains, v)
	}
	return t
}

// AddARecord is a convenience that adds the specified ARecord to the Domain
func (t *Domain) AddARecords(records ...*ARecord) *Domain {
	for _, v := range records {
		t.ARecords = append(t.ARecords, v)
	}
	return t
}

func (t *Domain) AddPtrRecords(records ...*PTRrecord) *Domain {
	for _, v := range records {
		t.PtrRecords = append(t.PtrRecords, v)
	}
	return t
}

// CNameRecord is a convenience that adds the specified CNameRecord to the Domain
func (t *Domain) AddCNameRecords(records ...*CNameRecord) *Domain {
	for _, v := range records {
		t.CnameRecords = append(t.CnameRecords, v)
	}
	return t
}
