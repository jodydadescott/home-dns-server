package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/miekg/dns"
	"go.uber.org/zap"

	"github.com/jodydadescott/home-dns-server/types"
)

const (
	defaultDnsProto  = types.ProtoTypeUDP
	defaultDnsPort   = 53
	defaultDnsDomain = "home"
)

type NetPort = types.NetPort
type ProtoType = types.ProtoType
type InternalUseDomain = types.InternalUseDomain
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type CNameRecord = types.CNameRecord

type Config struct {
	Providers   []Provider
	Debug       bool
	Listen      *NetPort
	Nameservers []*NetPort
}

func (t *Config) AddProvider(provider Provider) {
	t.Providers = append(t.Providers, provider)
}

type Client struct {
	debug        bool
	listen       *NetPort
	domains      []string
	udpDnsClient *dns.Client
	tcpDnsClient *dns.Client
	providers    []Provider
	aRecords     map[string]*ARecord
	ptrRecords   map[string]*PTRrecord
	cnameRecords map[string]*CNameRecord
	nameservers  []*NetPort
}

type Provider interface {
	GetDomain() (*InternalUseDomain, error)
}

func New(config *Config) *Client {

	if config == nil {
		panic("config is required")
	}

	listen := config.Listen

	if listen == nil {
		listen = &NetPort{}
	}

	switch listen.Proto {

	case types.ProtoTypeUDP, types.ProtoTypeTCP:

	default:
		listen.Proto = types.ProtoTypeUDP
	}

	if listen.Port <= 0 {
		listen.Port = defaultDnsPort
	}

	var nameservers []*NetPort
	for _, nameserver := range config.Nameservers {

		switch nameserver.Proto {

		case types.ProtoTypeUDP, types.ProtoTypeTCP:

		default:
			nameserver.Proto = types.ProtoTypeUDP
		}

		if nameserver.Port <= 0 {
			nameserver.Port = defaultDnsPort
		}

		nameservers = append(nameservers, nameserver)
	}

	return &Client{
		debug:        config.Debug,
		listen:       listen,
		udpDnsClient: &dns.Client{Net: "udp", SingleInflight: true},
		tcpDnsClient: &dns.Client{Net: "tcp", SingleInflight: true},
		aRecords:     make(map[string]*ARecord),
		ptrRecords:   make(map[string]*PTRrecord),
		cnameRecords: make(map[string]*CNameRecord),
		nameservers:  nameservers,
		providers:    config.Providers,
	}
}

func (t *Client) Run(ctx context.Context) error {

	u := &uniq{}

	for _, provider := range t.providers {

		domain, err := provider.GetDomain()

		u.add(domain.Domain)

		if err != nil {
			return err
		}

		for _, r := range domain.ARecords {
			t.aRecords[r.GetKey()] = r
		}

		for _, r := range domain.PtrRecords {
			t.ptrRecords[r.GetKey()] = r
		}

		for _, r := range domain.CnameRecords {
			t.cnameRecords[r.GetKey()] = r
		}

	}

	for _, v := range u.names {
		zap.L().Debug(fmt.Sprintf("Adding domain %s to be handled locally", v))
		dns.HandleFunc(v+".", t.handleLocal)
	}

	dns.HandleFunc("10.in-addr.arpa.", t.handleLocal)
	dns.HandleFunc("168.192.in-addr.arpa.", t.handleLocal)
	dns.HandleFunc("0.0.16.127.in-addr.arpa.", t.handleLocal)
	dns.HandleFunc("0.0.168.192.in-addr.arpa.", t.handleLocal)

	if len(t.nameservers) > 0 {
		for _, v := range t.nameservers {
			zap.L().Debug(fmt.Sprintf("Forwarding to nameserver %s : %s", v.IP, v.Proto.String()))
		}

		dns.HandleFunc(".", t.handleRemote)

	} else {
		zap.L().Debug("Forwarding to nameservers is not enabled")
	}

	server := &dns.Server{Addr: t.listen.IP + ":" + strconv.Itoa(t.listen.Port), Net: t.listen.Proto.String()}

	go func() {
		<-ctx.Done()
		zap.L().Debug("Shutting down")
		server.Shutdown()
	}()

	zap.L().Info(fmt.Sprintf("Starting server on %s/%s", t.listen.IP+":"+strconv.Itoa(t.listen.Port), t.listen.Proto.String()))

	return server.ListenAndServe()
}

func (t *Client) handleLocal(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:

		for _, q := range m.Question {

			switch q.Qtype {

			case dns.TypeA:
				lookup := t.aRecords[q.Name]
				if lookup != nil {
					record := fmt.Sprintf("%s A %s", q.Name, lookup.GetValue())

					if t.debug {
						zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
					}

					rr, err := dns.NewRR(record)

					if t.debug {
						zap.L().Debug(record)
					}

					if err == nil {
						m.Answer = append(m.Answer, rr)
					} else {
						zap.L().Error(err.Error())
					}
				} else {
					if t.debug {
						zap.L().Debug(fmt.Sprintf("fail -> %s has no A record", q.Name))
					}
				}

			case dns.TypePTR:
				lookup := t.ptrRecords[q.Name]
				if lookup != nil {
					record := fmt.Sprintf("%s PTR %s", q.Name, lookup.GetValue())

					if t.debug {
						zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
					}

					rr, err := dns.NewRR(record)

					if t.debug {
						zap.L().Debug(record)
					}

					if err == nil {
						m.Answer = append(m.Answer, rr)
					} else {
						zap.L().Error(err.Error())
					}
				} else {
					if t.debug {
						zap.L().Debug(fmt.Sprintf("fail -> %s has no PTR record", q.Name))
					}
				}

			case dns.TypeCNAME:
				lookup := t.cnameRecords[q.Name]
				if lookup != nil {
					record := fmt.Sprintf("%s CNAME %s", q.Name, lookup.GetValue())

					if t.debug {
						zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
					}

					rr, err := dns.NewRR(record)

					if t.debug {
						zap.L().Debug(record)
					}

					if err == nil {
						m.Answer = append(m.Answer, rr)
					} else {
						zap.L().Error(err.Error())
					}
				} else {
					if t.debug {
						zap.L().Debug(fmt.Sprintf("fail -> %s has no CNAME record", q.Name))
					}
				}

			}
		}

	}

	w.WriteMsg(m)
}

func (t *Client) handleRemote(w dns.ResponseWriter, r *dns.Msg) {

	dnsClient := t.tcpDnsClient

	for _, nameserver := range t.nameservers {

		switch nameserver.Proto {

		case types.ProtoTypeTCP:
			dnsClient = t.tcpDnsClient

		case types.ProtoTypeUDP:
			dnsClient = t.udpDnsClient

		}

		if r, _, err := dnsClient.Exchange(r, nameserver.GetIPColonPort()); err == nil {
			if r.Rcode == dns.RcodeSuccess {
				r.Compress = true
				w.WriteMsg(r)

				if t.debug {
					zap.L().Debug(fmt.Sprintf("Nameserver %s responded", nameserver.GetIPColonPort()))
				}

				return
			}
		}

	}

	zap.L().Error("failure to forward request")

	m := new(dns.Msg)
	m.SetReply(r)
	m.SetRcode(r, dns.RcodeServerFailure)
	w.WriteMsg(m)
}

type uniq struct {
	names []string
}

// func msgString(r *dns.Msg) string {
// 	for _, v := range r.Question {
// 		v.Name
// 		v.
// 	}
// }

// // Msg contains the layout of a DNS message.
// type Msg struct {
// 	MsgHdr
// 	Compress bool       `json:"-"` // If true, the message will be compressed when converted to wire format.
// 	Question []Question // Holds the RR(s) of the question section.
// 	Answer   []RR       // Holds the RR(s) of the answer section.
// 	Ns       []RR       // Holds the RR(s) of the authority section.
// 	Extra    []RR       // Holds the RR(s) of the additional section.
// }

func (t *uniq) add(input string) {

	for _, v := range t.names {
		if v == input {
			return
		}
	}

	t.names = append(t.names, input)
}
