package gobravia

import (
	"encoding/json"
	"encoding/xml"
)

type Tv interface {
	GetCommands() bool
	SendCommand(code string)
}

type BraviaTV struct {
	Address   string
	Envelope  Envelope
	Pin       string
	Mac       string
	Commands  map[string]string
	Connected bool
}

type Envelope struct {
	XMLName       xml.Name `xml:"s:Envelope"`
	EncodingStyle string   `xml:"s:encodingStyle,attr"`
	Xmlns         string   `xml:"xmlns:s,attr"`
	SendIRCC      SendIRCC `xml:"s:Body>u:X_SendIRCC"`
}

type SendIRCC struct {
	XMLName  xml.Name `xml:"u:X_SendIRCC"`
	Xmlns    string   `xml:"xmlns:u,attr"`
	IRCCCode string   `xml:"IRCCCode"`
}

type CodeResponse struct {
	ID        int               `json:"id"`
	RawResult []json.RawMessage `json:"result"`
	Header    *Header           `json:"-"`
	Values    []*Value          `json:"-"`
}

type Header struct {
	Bundled bool   `json:"bundled"`
	Type    string `json:"type"`
}

type Value struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
