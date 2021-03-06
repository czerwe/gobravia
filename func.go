package gobravia

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/Sirupsen/logrus"
	wol "github.com/ghthor/gowol"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetBravia(address string, pin string, mac string) *BraviaTV {
	retVal := BraviaTV{Address: address, Pin: pin, Mac: mac}
	return &retVal
}

func (envelope *Envelope) GetRequestXML(code string) ([]byte, error) {
	header := []byte("<?xml version=\"1.0\"?>")

	envelope.EncodingStyle = "http://schemas.xmlsoap.org/soap/encoding/"
	envelope.Xmlns = "http://schemas.xmlsoap.org/soap/envelope/"
	envelope.SendIRCC.Xmlns = "urn:schemas-sony-com:service:IRCC:1"
	envelope.SendIRCC.IRCCCode = code
	retVal, ok := xml.Marshal(envelope)

	return append(header, retVal...), ok
}

type ComGet struct {
	Id      int      `json:"id"`
	Method  string   `json:"method"`
	Version string   `json:"version"`
	Params  []string `json:"params"`
}

func (tv *BraviaTV) Poweron(bcast string) {
	wol.MagicWake(tv.Mac, bcast)
}

func (tv *BraviaTV) GetCommands() bool {
	commands := make(map[string]string)

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{Timeout: timeout}
	url := fmt.Sprintf("http://%v/sony/system", tv.Address)
	// url := fmt.Sprintf("https://czerwe.no-ip.org/testserver/bravia")

	var comstr = ComGet{
		Id:      10,
		Method:  "getRemoteControllerInfo",
		Version: "1.0",
		Params:  make([]string, 0),
	}

	bytestring, err := json.Marshal(comstr)

	if err != nil {
		log.Error(err)
		return false
	}

	jsonreader := bytes.NewReader(bytestring)

	req, err := http.NewRequest("POST", url, jsonreader)

	if err != nil {
		log.Error(err)
		return false
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		tv.Connected = false
		log.Error(err)
		return false
	} else {
		tv.Connected = true
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	coderesp := &CodeResponse{}

	err = json.Unmarshal(body, &coderesp)
	if err != nil {
		log.Error("Cannot unmarshal the codelist")
		return false
	}

	coderesp.Header = &Header{}
	err = json.Unmarshal(coderesp.RawResult[0], &coderesp.Header)
	if err != nil {
		log.Error("Cannot unmarshal the header")
		return false
	}

	coderesp.Values = []*Value{}
	json.Unmarshal(coderesp.RawResult[1], &coderesp.Values)
	if err != nil {
		log.Error("Cannot unmarshal the values")
		return false
	}

	for _, k := range coderesp.Values {
		curVal := *k
		strings.ToLower(curVal.Name)
		commands[strings.ToLower(curVal.Name)] = curVal.Value
	}

	tv.Commands = commands
	return true
}

func (tv *BraviaTV) SendCommand(code string) {

	client := &http.Client{}
	url := fmt.Sprintf("http://%v/sony/IRCC", tv.Address)

	bytestring, _ := tv.Envelope.GetRequestXML(code)
	jsonreader := bytes.NewReader(bytestring)

	req, err := http.NewRequest("POST", url, jsonreader)

	if err != nil {
		log.Error(err)
		return
	}

	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("X-Auth-PSK", tv.Pin)
	req.Header.Add("SOAPACTION", "urn:schemas-sony-com:service:IRCC:1#X_SendIRCC")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(body))

}

func (tv *BraviaTV) SearchCode(code string) (string, bool) {
	code, ok := tv.Commands[code]
	return code, ok
}

func (tv *BraviaTV) PrintCodes() {

	var count int32

	count = 0
	for k, _ := range tv.Commands {
		count++

		if count%3 == 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("%-23v", k)
	}
	fmt.Printf("\n")
}

func (tv *BraviaTV) SendAlias(alias string) bool {
	code, ok := tv.SearchCode(alias)
	if ok {
		tv.SendCommand(code)
	}
	return ok
}
