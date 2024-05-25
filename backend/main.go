package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Lease4 struct {
	IPAddress                string `json:"ip-address"`
	HWAddress                string `json:"hw-address"`
	ValidLifetime            int64  `json:"valid-lft"`
	ClientID                 string `json:"client-id"`
	SubnetID                 int    `json:"subnet-id"`
	ClientLastTransactioTime int64  `json:"cltt"`
}

func (lease Lease4) ExpiresAt() int64 {
	return lease.ClientLastTransactioTime + lease.ValidLifetime
}

type Lease4Dto struct {
	IPAddress string `json:"ipAddress"`
	HWAddress string `json:"hwAddress"`
	ExpiresAt int64  `json:"expiresAt"`
}

type Lease4Response struct {
	Arguments struct {
		Leases []Lease4 `json:"leases"`
	} `json:"arguments"`
	Result int `json:"result"`
}

func hello(w http.ResponseWriter, req *http.Request) {

	url := "http://172.16.0.1:8010/"

	// JSON request payload to get all IPv4 leases
	requestPayload := []byte(`{
		"command": "lease4-get-all",
		"service": [ "dhcp4" ]
	}`)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK response: %v", resp.StatusCode)
	}

	// Parse the JSON response
	var leaseResponses []Lease4Response
	if err := json.Unmarshal(body, &leaseResponses); err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		log.Fatalf("Raw response: %s", string(body))
	}

	var leaseResponse Lease4Response = leaseResponses[0]
	// Check the result
	if leaseResponse.Result != 0 {
		log.Fatalf("Error in response result: %v", leaseResponse.Result)
	}

	var responses = make([]Lease4Dto, 0, len(leaseResponse.Arguments.Leases))

	// Print all IPv4 leases
	for _, lease := range leaseResponse.Arguments.Leases {
		var response Lease4Dto

		response.IPAddress = lease.IPAddress
		response.HWAddress = lease.HWAddress
		response.ExpiresAt = lease.ExpiresAt()

		responses = append(responses, response)

		// fmt.Fprintf(w, "IP Address: %s, HW Address: %s, Client ID: %s, Subnet ID: %d, Valid Lifetime: %d, Cltt: %d, Expires: %s\n",
		// 	lease.IPAddress,
		// 	lease.HWAddress,
		// 	lease.ClientID,
		// 	lease.SubnetID,
		// 	lease.ValidLifetime,
		// 	lease.ClientLastTransactioTime,
		// 	time.Unix(lease.ExpiresAt(), 0).String())
	}

	response, err := json.Marshal(responses)
	if err != nil {
		log.Fatalf("Error marshalling JSON response: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // URL of the Kea control agent API endpoint
	w.Write(response)

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/api/ipv4-leases", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}
