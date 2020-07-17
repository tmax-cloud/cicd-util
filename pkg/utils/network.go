package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func PostReq(uri string, header, params map[string]string) (string, error) {
	c := http.Client{}

	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if math.Floor(float64(resp.StatusCode)/100.) == 2 {
		return string(body), nil
	} else {
		return "", fmt.Errorf("error code: %d, message: %s", resp.StatusCode, string(body))
	}
}

func PodIp() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile("([^ \t]*)[ \t]+" + hostName)
	if err != nil {
		return "", err
	}

	file, err := os.Open("/etc/hosts")
	if err != nil {
		return "", err
	}
	defer file.Close()

	ip := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		match := reg.FindStringSubmatch(text)
		if len(match) > 1 {
			ip = match[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	if ip == "" {
		return "", fmt.Errorf("cannot find current ip")
	} else {
		return ip, nil
	}
}
