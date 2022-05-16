package pingback

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/webmention/recv"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HandlePost(conf *common.Config, db db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			pingbackError(w, fmt.Errorf("pingback POST: Unable to read body: %v", err))
			return
		}
		rpc := &XmlRPCMethodCall{}
		err = xml.Unmarshal(body, rpc)
		if err != nil {
			pingbackError(w, fmt.Errorf("pingback POST: Unable to unmarshal XMLRPC %s: %v", body, err))
			return
		}

		if !validate(rpc, conf) {
			pingbackError(w, fmt.Errorf("pingback POST: malformed pingback request: %s", body))
			return
		}

		wm := mf.Mention{
			Source: rpc.Source(),
			Target: rpc.Target(),
		}
		receiver := &recv.Receiver{
			RestClient: &rest.HttpClient{},
			Conf:       conf,
			Repo:       db,
		}
		go receiver.Receive(wm)
		pingbackSuccess(w)
	}
}

func pingbackSuccess(w http.ResponseWriter) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param>
            <value>
                <string>
                    Thanks, bro. Will process this soon, pinky swear!
                </string>
            </value>
        </param>
    </params>
</methodResponse>`
	w.WriteHeader(200)
	w.Write([]byte(xml))
}

// according to the XML-RPC spec, always return a 200, but encode it into the XML.
func pingbackError(w http.ResponseWriter, err error) {
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
	// No longer interested in pingback errors, these are 99.9% badly formatted spam that clog up syslog
	// log.Error().Err(err).Msg("Pingback receive went wrong")
	w.WriteHeader(200)
	w.Write([]byte(xml))
}
