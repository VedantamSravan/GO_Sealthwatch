package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	SMC_USER      = "admin"
	SMC_PASSWORD  = "P@ck3t08.."
	SMC_HOST      = "10.91.170.206"
	SMC_TENANT_ID = "102"
)

func GetClient() (*http.Client, *http.Response) {
	URL := "https://" + SMC_HOST + "/token/v2/authenticate"
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "NextCom",
		Value:  "rb",
		Path:   "/",
		Domain: ".nextcom",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(URL)
	jar.SetCookies(u, cookies)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Jar: jar,
	}
	data := url.Values{
		"username": {SMC_USER},
		"password": {SMC_PASSWORD},
	}
	resp, err := client.PostForm(URL, data)
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	return client, resp
}

func get_top_ports() {
	var client, resp = GetClient()

	if resp.StatusCode == 200 {
		var URL = "https://" + SMC_HOST + "/sw-reporting/v1/tenants/" + SMC_TENANT_ID + "/flow-reports/top-ports/queries"
		end_timestamp := time.Now().UTC().Truncate(time.Millisecond)
		start_timestamp := end_timestamp.Add(time.Duration(-60) * time.Minute)
		values := map[string]string{"startTime": start_timestamp.Format("2006-01-02T15:04:05.000"), "endTime": end_timestamp.Format("2006-01-02T15:04:05.000"), "maxRows": "50"}
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Post(URL, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		var res = map[string]map[string]string{}
		if resp.StatusCode == 200 {
			responseData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal([]byte(responseData), &res)
			queryid := " "
			searchstatus := " "
			for _, element := range res {
				queryid = element["queryId"]
				searchstatus = element["status"]
			}

			URL = "https://" + SMC_HOST + "/sw-reporting/v1/tenants/" + SMC_TENANT_ID + "/flow-reports/top-ports/queries/" + queryid
			for searchstatus != "COMPLETED" {
				searchresponseData, err := client.Get(URL)
				if err != nil {
					log.Println(err)
				}
				defer searchresponseData.Body.Close()
				search_response_Data_result, err := ioutil.ReadAll(searchresponseData.Body)
				if err != nil {
					log.Println(err)
				}
				defer searchresponseData.Body.Close()
				json.Unmarshal([]byte(search_response_Data_result), &res)
				for _, element := range res {
					queryid = element["queryId"]
					searchstatus = element["status"]
				}
				time.Sleep(1 * time.Second)

			}

			URL = "https://" + SMC_HOST + "/sw-reporting/v1/tenants/" + SMC_TENANT_ID + "/flow-reports/top-ports/results/" + queryid
			top_ports_response_Data, err := client.Get(URL)
			if err != nil {
				log.Println(err)
			}
			defer top_ports_response_Data.Body.Close()
			top_ports_response_data_results, err := ioutil.ReadAll(top_ports_response_Data.Body)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(top_ports_response_data_results))

		}
	}

}
func main() {
	get_top_ports()

}
