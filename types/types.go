package types

import (
	"fmt"
	"regexp"

	"github.com/jinzhu/copier"
	"github.com/jodydadescott/unifi-go-sdk/unifi"
)

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

func (t ProtoType) String() string {

	switch t {

	case ProtoTypeUDP:
		return string(t)

	case ProtoTypeTCP:
		return string(t)

	}

	panic("Invalid proto type")
}

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

func (t *NetPort) SetProtoTypeTCP() {
	t.Proto = ProtoTypeTCP
}

func (t *NetPort) SetProtoTypeUDP() {
	t.Proto = ProtoTypeUDP
}

func (t *NetPort) GetIPColonPort() string {
	if t.ipColonPort == "" {
		t.ipColonPort = t.IP + ":" + fmt.Sprint(t.Port)
	}
	return t.ipColonPort
}

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

func (t *ARecord) GetKey() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

func (t *ARecord) GetValue() string {
	return t.IP
}

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

func (t *CNameRecord) GetKey() string {
	if t.fqdnAlias == "" {
		t.fqdnAlias = cleanHostname(t.AliasHostname) + "." + t.AliasDomain + "."
	}
	return t.fqdnAlias
}

func (t *CNameRecord) GetValue() string {
	if t.fqdnTarget == "" {
		t.fqdnTarget = cleanHostname(t.TargetHostname) + "." + t.TargetDomain + "."
	}
	return t.fqdnTarget
}

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

func (t *PTRrecord) GetKey() string {
	return t.ARPA
}

func (t *PTRrecord) GetValue() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

type Config struct {
	Notes       string        `json:"notes,omitempty" yaml:"notes,omitempty"`
	Unifi       *UnifiConfig  `json:"unifiConfig,omitempty" yaml:"unifiConfig,omitempty"`
	Listen      *NetPort      `json:"listen,omitempty" yaml:"listen,omitempty"`
	Static      *StaticConfig `json:"static,omitempty" yaml:"static,omitempty"`
	Nameservers []*NetPort    `json:"nameservers,omitempty" yaml:"nameservers,omitempty"`
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

func (t *Config) AddNameserver(r *NetPort) *Config {
	t.Nameservers = append(t.Nameservers, r)
	return t
}

type UnifiConfig struct {
	unifi.Config
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domain  string `json:"domain,omitempty" yaml:"domain,omitempty"`
}

// Clone return copy
func (t *UnifiConfig) Clone() *UnifiConfig {
	c := &UnifiConfig{}
	copier.Copy(&c, &t)
	return c
}

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

type Domain struct {
	Domain       string         `json:"domain,omitempty" yaml:"dnsDomain,omitempty"`
	ARecords     []*ARecord     `json:"aRecords,omitempty" yaml:"aRecords,omitempty"`
	CnameRecords []*CNameRecord `json:"cnameRecords,omitempty" yaml:"cnameRecords,omitempty"`
}

// Clone return copy
func (t *Domain) Clone() *Domain {
	c := &Domain{}
	copier.Copy(&c, &t)
	return c
}

func (t *StaticConfig) AddDomain(domain *Domain) *StaticConfig {
	t.Domains = append(t.Domains, domain)
	return t
}

func (t *Domain) AddARecord(r *ARecord) *Domain {
	t.ARecords = append(t.ARecords, r)
	return t
}

func (t *Domain) AddCNameRecord(r *CNameRecord) *Domain {
	t.CnameRecords = append(t.CnameRecords, r)
	return t
}

type InternalUseDomain struct {
	Domain       string         `json:"domain,omitempty" yaml:"dnsDomain,omitempty"`
	ARecords     []*ARecord     `json:"aRecords,omitempty" yaml:"aRecords,omitempty"`
	PtrRecords   []*PTRrecord   `json:"ptrRecords,omitempty" yaml:"ptrRecords,omitempty"`
	CnameRecords []*CNameRecord `json:"cnameRecords,omitempty" yaml:"cnameRecords,omitempty"`
}

// Clone return copy
func (t *InternalUseDomain) Clone() *InternalUseDomain {
	c := &InternalUseDomain{}
	copier.Copy(&c, &t)
	return c
}

func (t *InternalUseDomain) AddARecord(r *ARecord) *InternalUseDomain {
	t.ARecords = append(t.ARecords, r)
	return t
}

func (t *InternalUseDomain) AddCNameRecord(r *CNameRecord) *InternalUseDomain {
	t.CnameRecords = append(t.CnameRecords, r)
	return t
}

func (t *InternalUseDomain) AddPtrRecord(r *PTRrecord) *InternalUseDomain {
	t.PtrRecords = append(t.PtrRecords, r)
	return t
}
