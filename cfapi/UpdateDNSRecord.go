package cfapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"cloudflare-ddns-go/model"
	"cloudflare-ddns-go/util"
)

type UpdateDNSRecordResponseBody struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	DnsType string   `json:"type"`
	Comment string   `json:"comment"`
	Id      string   `json:"id"`
	Tags    []string `json:"tags"`
	Ttl     int      `json:"ttl"`
}

func UpdateDNSRecord(localFileSettingsAuthentification model.AuthenticationStruct, localFileSettingsItemDNSRecord model.MyDNSRecordStruct) {

	logger := util.GetLogger()

	url := "https://api.cloudflare.com/client/v4/zones/" + localFileSettingsAuthentification.ZoneId + "/dns_records/" + localFileSettingsItemDNSRecord.Id
	logger.Info("Request url", "url", url)

	payloadJSON, _ := json.Marshal(localFileSettingsItemDNSRecord)
	payloadReader := bytes.NewReader(payloadJSON)

	req, _ := http.NewRequest("PATCH", url, payloadReader)

	req.Header.Add("X-Auth-Key", localFileSettingsAuthentification.ApiKey)
	req.Header.Add("X-Auth-Email", localFileSettingsAuthentification.ApiEmail)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	defer res.Body.Close()

	resBody := make(map[string]interface{})
	err := json.Unmarshal(body, &resBody)

	if err != nil || !resBody["success"].(bool) {
		logger.Error("Update Record failed!", " Marshal error", err, "Response error: ", resBody["errors"])
	} else {
		logger.Error("Update Record success!", "localFileSettingsItemDNSRecord", localFileSettingsItemDNSRecord)
	}
}
