package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

type Response struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Params []IcpResponse `json:"params"`
	Status int           `json:"status"`
	Time   int           `json:"time"`
}

type IcpResponse struct {
	ContentTypeName  string `json:"contentTypeName"`
	Domain           string `json:"domain"`
	LeaderName       string `json:"leaderName"`
	LimitAccess      string `json:"limitAccess"`
	MainLicence      string `json:"mainLicence"`
	NatureName       string `json:"natureName"`
	IceLicence       string `json:"iceLicence"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
}

func sign(params map[string]string, secret string) string {
	var keys []string
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		if signStr != "" {
			signStr += "&"
		}
		signStr += fmt.Sprintf("%s=%s", k, params[k])
	}
	signStr += secret

	h := md5.New()
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

func getdomain(param string, star string) {
	timestamp := time.Now().Unix()
	timestampStr := strconv.Itoa(int(timestamp))
	//param = "中国人民财产保险股份有限公司"
	params := map[string]string{
		"appid":     "yourid",
		"params":    param,
		"timestamp": timestampStr,
	}
	secret := "yoursecret"
	//fmt.Println(sign(params, secret))
	sign := sign(params, secret)
	urlstr := "https://www.icpapi.com/api/v1"
	paramss := url.Values{
		"appid":     {"yourid"},
		"params":    {param},
		"timestamp": {timestampStr},
		"sign":      {sign},
	}
	//fmt.Println(timestampStr)
	resp, err := http.PostForm(urlstr, paramss)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodystr := string(body)
	//fmt.Println(bodystr)
	var responses Response
	err = json.Unmarshal([]byte(bodystr), &responses)
	if err != nil {
		log.Fatal(err)
	}

	for _, response := range responses.Params {
		if star != "" {
			tmp := "*." + response.Domain
			fmt.Println(tmp)
		} else {
			fmt.Println(response.Domain)
		}

	}
}

func main() {
	var name string
	var filename string
	var param string
	var outfile string
	var star string
	flag.StringVar(&name, "n", "", "-n 单位名称或域名或备案号")
	flag.StringVar(&filename, "f", "", "-f input file")
	flag.StringVar(&outfile, "o", "", "-o output file")
	flag.StringVar(&star, "s", "", "")
	flag.Parse()
	if name != "" && filename == "" {
		param = name
		getdomain(param, "")
	}
	if name == "" && filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		for _, v := range lines {
			getdomain(v, star)
		}
	}
}
