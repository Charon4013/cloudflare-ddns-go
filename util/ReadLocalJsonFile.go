package util

import (
	"encoding/json"
	"io"
	"os"

	"cloudflare-ddns-go/model"
)

func ReadLocalJsonFile() (localFileSettings model.MyDNSJSONStruct) {

	jsonContent, err := os.Open("myDNS.json")

	if err != nil {
		logger.Error("Can not open config.json", "Error", err)
		os.Exit(1)
	}

	defer jsonContent.Close()

	byteResult, _ := io.ReadAll(jsonContent)

	json.Unmarshal([]byte(byteResult), &localFileSettings)

	checkResult := CheckLocalFileSettings(localFileSettings)
	if !checkResult {
		logger.Error("Check local file failed!")
		os.Exit(1)
	}

	return
}

func CheckLocalFileSettings(localFileSettings model.MyDNSJSONStruct) bool {

	if localFileSettings.ApiEmail == "" || localFileSettings.ApiKey == "" || localFileSettings.ZoneId == "" {
		logger.Error("Api_key or api_email or zone_id is empty", "localFileSettings", localFileSettings)
		return false
	}

	if len(localFileSettings.MyDNSRecordStruct) == 0 {
		logger.Error("DNSRecords is empty", "localFileSettings", localFileSettings)
		return false
	}

	refName := localFileSettings.MyDNSRecordStruct[0].Name
	refType := localFileSettings.MyDNSRecordStruct[0].DnsType

	for _, localFileDNSRecordItem := range localFileSettings.MyDNSRecordStruct[1:] {
		if localFileDNSRecordItem.Name == refName && localFileDNSRecordItem.DnsType == refType {
			logger.Error("Duplicate DNS name and type", "localFileDNSRecordItem", localFileDNSRecordItem)
			return false
		}
	}
	return true
}
