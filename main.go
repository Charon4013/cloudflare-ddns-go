package main

import (
	"cloudflare-ddns-go/cfapi"
	"cloudflare-ddns-go/util"
)

func main() {

	util.InitLogger()
	defer util.CloseLogger()

	logger := util.GetLogger()

	logger.Info("#1 Script start")

	// Get my DNS settings from MyDns.json
	localFileSettings := util.ReadLocalJsonFile()
	logger.Info("#2 Check local file settings")

	// Get local IP from website
	localIPv4 := util.GetLocalIPFromWeb(4)
	localIPv6 := util.GetLocalIPFromWeb(6)
	logger.Info("#3 Get local IP from website")

	// Get all the DNS records from Cloudflare
	remoteDnsRecordsList := cfapi.ListDNSRecords(localFileSettings.AuthenticationStruct)
	logger.Info("#4 Get remote DNS records from Cloudflare")

	// Check if remote IP address is up to date
	for _, localFileSettingsItem := range localFileSettings.MyDNSRecordStruct {
		logger.Info("#5 Check IP address and update DNS record")
		for _, remoteDnsRecordsListItem := range remoteDnsRecordsList {

			if localFileSettingsItem.Name == remoteDnsRecordsListItem.Name && localFileSettingsItem.DnsType == remoteDnsRecordsListItem.DnsType {

				if localFileSettingsItem.Content == "DDNS" || localFileSettingsItem.Content == "ddns" {

					if remoteDnsRecordsListItem.DnsType == "A" {
						localFileSettingsItem.Content = localIPv4
					} else if remoteDnsRecordsListItem.DnsType == "AAAA" {
						localFileSettingsItem.Content = localIPv6
					}
				}

				localFileSettingsItem.Id = remoteDnsRecordsListItem.Id

				// Update DNS record to Cloudflare
				cfapi.UpdateDNSRecord(localFileSettings.AuthenticationStruct, localFileSettingsItem)
			}
		}
	}

	logger.Info("==================================Done!==================================")
}
