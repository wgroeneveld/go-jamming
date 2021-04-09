package pingback

import (
	"brainbaking.com/go-jamming/common"
	"strings"
)

func validate(rpc *XmlRPCMethodCall, conf *common.Config) bool {
	if rpc.MethodName != "pingback.ping" {
		return false
	}
	if len(rpc.Params.Parameters) != 2 {
		return false
	}

	target := rpc.Target()
	if !strings.HasPrefix(target, "http") {
		return false
	}
	_, err := conf.FetchDomain(target)
	if err != nil {
		return false
	}

	source := rpc.Source()
	if !strings.HasPrefix(source, "http") {
		return false
	}

	if source == target {
		return false
	}
	return true
}
