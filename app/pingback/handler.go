
package pingback

import (
	"encoding/xml"
	"github.com/rs/zerolog/log"
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/webmention/receive"
	"github.com/wgroeneveld/go-jamming/rest"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/wgroeneveld/go-jamming/common"
)

func HandlePost(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			pingbackError(w, "Unable to read request body")
		}
		rpc := &XmlRPCMethodCall{}
		err = xml.Unmarshal(body, rpc)
		if err != nil {
			pingbackError(w, "Unable to unmarshal XMLRPC request body")
		}

		if !validate(rpc, conf) {
			pingbackError(w, "malformed pingback request")
			return
		}

		wm := mf.Mention{
			Source: rpc.Source(),
			Target: rpc.Target(),
		}
		receiver := receive.Receiver{
			RestClient: &rest.HttpClient{},
			Conf:       conf,
		}
		go receiver.Receive(wm)
		pingbackSuccess(w, "Thanks, bro. Will process this soon, pinky swear!")
	}
}

var successXml = `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param>
            <value>
                <string>
                    {{ . }}
                </string>
            </value>
        </param>
    </params>
</methodResponse>
`
// compile once, execute as many times as needed.
var successTpl, _ = template.New("success").Parse(successXml)

func pingbackSuccess(w http.ResponseWriter, msg string) {
	w.WriteHeader(200)
	successTpl.Execute(w, msg)
}

// according to the XML-RPC spec, always return a 200, but encode it into the XML.
func pingbackError(w http.ResponseWriter, msg string) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <fault>
        <value>
            <struct>
                <member>
                    <name>
                        faultCode
                    </name>
                    <value>
                        <int>
                            0
                        </int>
                    </value>
                </member>
                <member>
                    <name>
                        faultString
                    </name>
                    <value>
                        <string>
                        	Sorry pal. Malformed request? Or something else, who knows...
                        </string>
                    </value>
                </member>
            </struct>
        </value>
    </fault>
</methodResponse>`
	log.Error().Str("msg", msg).Msg("Pingback receive went wrong")
	w.WriteHeader(200)
	w.Write([]byte(xml))
}


