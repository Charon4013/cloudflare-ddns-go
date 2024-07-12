package model

type MyDNSRecordStruct struct {
	Content string   `json:"content"` // if auto, then will get your host IP from website
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	DnsType string   `json:"type"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
	Ttl     int      `json:"ttl"`
	Id      string   `json:"id"`
}

type AuthenticationStruct struct {
	ApiEmail string `json:"api_email"`
	ApiKey   string `json:"api_key"`
	ZoneId   string `json:"zone_id"`
}

type MyDNSJSONStruct struct {
	AuthenticationStruct
	MyDNSRecordStruct []MyDNSRecordStruct `json:"myDNSRecords"`
}
