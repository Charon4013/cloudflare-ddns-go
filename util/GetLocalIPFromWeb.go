package util

import (
	"io"
	"net/http"
	"strings"
)

func GetLocalIPFromWeb(IPVersion uint) (localIP string) {

	logger = GetLogger()

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
			logger.Error("Read localIPv4 response body error!", "Error", err)
			localIP = ""
			return
		}
		localIP = strings.Replace(string(localIPv4Body), "\n", "", -1)
	}
	defer localIPRes.Body.Close()

	return
}
