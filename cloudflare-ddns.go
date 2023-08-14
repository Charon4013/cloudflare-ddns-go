package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// key in your information
const (
	auth_email      string = "" // "114514"
	auth_key        string = "" // "1919"
	zone_identifier string = "" // "810"
)

const (
	infoLogStringHead    string = "[INFO] "
	successLogStringHead string = "[SUCCESS] "
	failLogStringHead    string = "[FAILED] "
)

type dnsRecord struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	DnsType string `json:"type"`
	// Comment string `json:"comment"`
	// Tags    []string `json:"tags"`
	// Ttl int    `json:"ttl"`
	Id string `json:"id"`
}

func main() {
	currentTime := getTimeForLogs()
	fmt.Println(currentTime)

	dnsRecords := getAAndAAAADNSRecordsIdFromCF()
	checkLocalIPToCFAPI(dnsRecords)
}

func getLocalIPFromWeb(IPVersion uint) (localIP string) {
	var url string
	if IPVersion == 6 {
		url = "http://api-ipv6.ip.sb/ip"
	} else {
		url = "http://api-ipv4.ip.sb/ip"
	}
	localIPRes, localIPErr := http.Get(url)

	if localIPErr != nil || localIPRes.StatusCode != http.StatusOK {
		localIP = ""
		return
	} else {
		localIPv4Body, err := io.ReadAll(localIPRes.Body)
		if err != nil {
			fmt.Println(failLogStringHead+"Read localIPv4 response body error: ", err)
			localIP = ""
			return
		}
		localIP = strings.Replace(string(localIPv4Body), "\n", "", -1)
	}
	defer localIPRes.Body.Close()

	return
}

func getAAndAAAADNSRecordsIdFromCF() []dnsRecord {

	url := "https://api.cloudflare.com/client/v4/zones/" + zone_identifier + "/dns_records"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Key", auth_key)
	req.Header.Add("X-Auth-Email", auth_email)

	res, err := http.DefaultClient.Do(req)

	var body []byte

	if res != nil && err == nil {
		body, _ = io.ReadAll(res.Body)
	} else {
		fmt.Println(failLogStringHead+"Get dnsRecords failed! Error:", err)
		return nil
	}
	defer res.Body.Close()

	resBody := make(map[string]interface{})
	marshalErr := json.Unmarshal(body, &resBody)
	if marshalErr != nil {
		fmt.Println("Unmarshal error: ", marshalErr)
		return nil
	}

	if !resBody["success"].(bool) {
		fmt.Println("Request failed")
		return nil
	}

	var recordList []dnsRecord

	for _, recordValue := range resBody["result"].([]interface{}) {

		if recordValue.(map[string]interface{})["type"] == "A" || recordValue.(map[string]interface{})["type"] == "AAAA" {
			recordItem := dnsRecord{
				Content: recordValue.(map[string]interface{})["content"].(string),
				Name:    recordValue.(map[string]interface{})["name"].(string),
				DnsType: recordValue.(map[string]interface{})["type"].(string),
				Id:      recordValue.(map[string]interface{})["id"].(string),
			}
			recordList = append(recordList, recordItem)
		}
	}

	return recordList
}

func checkLocalIPToCFAPI(dnsRecords []dnsRecord) {

	if dnsRecords == nil {
		fmt.Println(failLogStringHead + "Get dnsRecords failed")
	}

	currentLocalIPv4 := getLocalIPFromWeb(4)
	currentLocalIPv6 := getLocalIPFromWeb(6)

	IPv4Parse := net.ParseIP(currentLocalIPv4)
	if IPv4 := IPv4Parse.To4(); IPv4 != nil {
		for _, recordValue := range dnsRecords {
			if recordValue.DnsType == "A" {
				if recordValue.Content == currentLocalIPv4 {
					fmt.Println(infoLogStringHead + recordValue.Name + "'s IP is up to date")
				} else {
					fmt.Println(infoLogStringHead + recordValue.Name + "'s IP is out of date")
					patchLocalIPToCFAPI(recordValue, currentLocalIPv4)
				}
			}
		}
		fmt.Println()
	} else {
		fmt.Println(failLogStringHead+"Get local IP failed", currentLocalIPv4)
	}

	IPv6Parse := net.ParseIP(currentLocalIPv4)
	if IPv6 := IPv6Parse.To16(); IPv6 != nil {
		for _, recordValue := range dnsRecords {
			if recordValue.DnsType == "AAAA" {
				if recordValue.Content == currentLocalIPv6 {
					fmt.Println(infoLogStringHead + recordValue.Name + "'s IPv6 is up to date")
				} else {
					fmt.Println(infoLogStringHead + recordValue.Name + "'s IPv6 is out of date")
					patchLocalIPToCFAPI(recordValue, currentLocalIPv6)
				}
			}
		}
		fmt.Println()
	} else {
		fmt.Println(failLogStringHead+"Get local IPv6 failed", currentLocalIPv6)
	}
}

func patchLocalIPToCFAPI(oldDnsRecord dnsRecord, currentIP string) {
	url := "https://api.cloudflare.com/client/v4/zones/" + zone_identifier + "/dns_records/" + oldDnsRecord.Id

	fmt.Println(oldDnsRecord)
	fmt.Println(url)

	payload := &dnsRecord{
		Content: currentIP,
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadJSON)

	req, _ := http.NewRequest("PATCH", url, payloadReader)

	req.Header.Add("X-Auth-Key", auth_key)
	req.Header.Add("X-Auth-Email", auth_email)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	defer res.Body.Close()

	resBody := make(map[string]interface{})
	err := json.Unmarshal(body, &resBody)

	if err != nil || !resBody["success"].(bool) {
		fmt.Println(failLogStringHead+"Update Record failed! Name:", oldDnsRecord.Name, "Type:", oldDnsRecord.DnsType, "Marshal error: ", err, "Response error: ", resBody["errors"])
	} else {
		fmt.Println(successLogStringHead+"Update Record success! Name:", oldDnsRecord.Name, "Type:", oldDnsRecord.DnsType)
	}
}

func getTimeForLogs() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}
