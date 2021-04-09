package pingback

import "encoding/xml"

/* e.g. (see tests)
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/kristien.html</string></value>
		</param>
		<param>
			<value><string>https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/</string></value>
		</param>
	</params>
</methodCall>
*/
type XmlRPCMethodCall struct {
	XMLName    xml.Name     `xml:"methodCall"`
	MethodName string       `xml:"methodName"`
	Params     XmlRPCParams `xml:"params"`
}

func (rpc *XmlRPCMethodCall) Source() string {
	return rpc.Params.Parameters[0].Value.String
}

func (rpc *XmlRPCMethodCall) Target() string {
	return rpc.Params.Parameters[1].Value.String
}

type XmlRPCParams struct {
	XMLName    xml.Name      `xml:"params"`
	Parameters []XmlRPCParam `xml:"param"`
}

type XmlRPCParam struct {
	XMLName xml.Name    `xml:"param"`
	Value   XmlRPCValue `xml:"value"`
}

type XmlRPCValue struct {
	String string `xml:"string"`
}
