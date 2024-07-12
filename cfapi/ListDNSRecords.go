package cfapi

import (
	"encoding/json"
	"io"
	"net/http"

	"cloudflare-ddns-go/model"
	"cloudflare-ddns-go/util"
)

type ListDNSRecordsResponse struct {
	Result   interface{}   `json:"result"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Success  bool          `json:"success"`
}

func ListDNSRecords(myDNSJSONStruct model.AuthenticationStruct) (myCloudflareDNSRecordList []model.MyDNSRecordStruct) {

	logger := util.GetLogger()

	url := "https://api.cloudflare.com/client/v4/zones/" + myDNSJSONStruct.ZoneId + "/dns_records"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", myDNSJSONStruct.ApiEmail)
	req.Header.Add("X-Auth-Key", myDNSJSONStruct.ApiKey)

	res, err := http.DefaultClient.Do(req)

	var body []byte

	if res != nil && err == nil {
		body, _ = io.ReadAll(res.Body)
	} else {
		logger.Error("Get dnsRecords failed!", "error", err)
		return nil
	}
	defer res.Body.Close()

	resBody := ListDNSRecordsResponse{}
	marshalErr := json.Unmarshal(body, &resBody)
	if marshalErr != nil {
		logger.Error("Unmarshal error!", "error", marshalErr)
		return nil
	}

	if !resBody.Success {
		logger.Error("Request failed")
		return nil
	}

	for _, recordValue := range resBody.Result.([]interface{}) {

		recordItem := model.MyDNSRecordStruct{
			Content: recordValue.(map[string]interface{})["content"].(string),
			Name:    recordValue.(map[string]interface{})["name"].(string),
			DnsType: recordValue.(map[string]interface{})["type"].(string),
			Id:      recordValue.(map[string]interface{})["id"].(string),
			Proxied: recordValue.(map[string]interface{})["proxied"].(bool),
			Ttl:     int(recordValue.(map[string]interface{})["ttl"].(float64)),
		}
		myCloudflareDNSRecordList = append(myCloudflareDNSRecordList, recordItem)
	}

	return

}
