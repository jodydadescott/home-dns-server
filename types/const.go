package types

import (
	"time"

	"github.com/jodydadescott/home-dns-server/types/proto"
)

const (
	CodeVersion      = "1.0.3"
	DefaultDomain    = "home"
	DefaultDnsProto  = proto.UDP
	DefaultDnsPort   = 53
	DefaultDnsDomain = "home"
	DefaultRefresh   = time.Hour
	DefaultHTTPPort  = 8080
)
