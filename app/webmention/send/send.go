package send

import (
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/pingback/send"
)

func mention() {
	send.SendPingbackToEndpoint("endpoint", mf.Mention{})
}
