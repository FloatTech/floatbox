package web

import (
	"github.com/fumiama/terasu/ip"
)

// IsSupportIPv6 检查本机是否支持 ipv6
var IsSupportIPv6 = &ip.IsIPv6Available
