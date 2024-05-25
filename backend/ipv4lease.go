package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Lease4 struct {
	IPAddress                string `json:"ip-address"`
	HWAddress                string `json:"hw-address"`
	ValidLifetime            int64  `json:"valid-lft"`
	ClientID                 string `json:"client-id"`
	SubnetID                 int    `json:"subnet-id"`
	ClientLastTransactioTime int64  `json:"cltt"`
}

type Lease4Tracker struct {
	lease   Lease4
	addedAt time.Time
}

func (lease Lease4) ExpiresAt() int64 {
	return lease.ClientLastTransactioTime + lease.ValidLifetime
}

type Lease4Dto struct {
	IPAddress string `json:"ipAddress"`
	HWAddress string `json:"hwAddress"`
	ExpiresAt int64  `json:"expiresAt"`
}

type Lease4KeaResponse struct {
	Arguments struct {
		Leases []Lease4 `json:"leases"`
	} `json:"arguments"`
	Result int `json:"result"`
}

type LeasesResponse struct {
	Leases        []Lease4Dto `json:"leases"`
	NewLeases     []Lease4Dto `json:"newLeases"`
	RemovedLeases []Lease4Dto `json:"removedLeases"`
}

func (ipv4Lease Lease4) ToDto() Lease4Dto {
	return Lease4Dto{
		IPAddress: ipv4Lease.IPAddress,
		HWAddress: ipv4Lease.HWAddress,
		ExpiresAt: ipv4Lease.ExpiresAt(),
	}
}

var (
	lastLeases    map[string]Lease4        = make(map[string]Lease4)
	newLeases     map[string]Lease4Tracker = make(map[string]Lease4Tracker)
	removedLeases map[string]Lease4Tracker = make(map[string]Lease4Tracker)
	removedTtlSec int64                    = 60
)

func getLeasesFromKea() []Lease4 {
	url := config.KeaControlAgentURL
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK response: %v", resp.StatusCode)
	}

	// Parse the JSON response
	var leaseResponses []Lease4KeaResponse
	if err := json.Unmarshal(body, &leaseResponses); err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		log.Fatalf("Raw response: %s", string(body))
	}

	var leaseResponse Lease4KeaResponse = leaseResponses[0]
	// Check the result
	if leaseResponse.Result != 0 {
		log.Fatalf("Error in response result: %v", leaseResponse.Result)
	}

	return leaseResponse.Arguments.Leases
}

func updateLeases() LeasesResponse {

	var keaLeases = getLeasesFromKea()

	/* keaLeases = append(keaLeases, Lease4{
		IPAddress:                "192.168.0.1" + strconv.Itoa(rand.IntN(9)),
		HWAddress:                "00:00:00:00:00:00",
		ValidLifetime:            0,
		ClientID:                 "asd",
		SubnetID:                 0,
		ClientLastTransactioTime: 0,
	})*/
	var resultLeases = make([]Lease4Dto, 0, len(keaLeases))

	var possibleNewLeases = make([]Lease4, 0, len(keaLeases))
	var possibleRemovedLeases = make([]Lease4, 0, len(keaLeases))

	// Find new leases
	for _, lease := range keaLeases {
		var _, ok = lastLeases[lease.IPAddress]
		if !ok {
			possibleNewLeases = append(possibleNewLeases, lease)
		}
	}

	// Find removed leases
	for ipAddress, lease := range lastLeases {
		var contained = false
		for _, keaLease := range keaLeases {
			if keaLease.IPAddress == ipAddress {
				contained = true
				break
			}
		}
		if !contained {
			possibleRemovedLeases = append(possibleRemovedLeases, lease)
		}
	}

	// Print all IPv4 leases
	for _, lease := range keaLeases {
		response := lease.ToDto()
		resultLeases = append(resultLeases, response)
	}

	var newLeaseDtos = make([]Lease4Dto, len(possibleNewLeases))
	var removedLeaseDtos = make([]Lease4Dto, len(possibleRemovedLeases))

	for i, lease := range possibleNewLeases {
		newLeaseDtos[i] = lease.ToDto()
	}

	// for i, lease := range possibleRemovedLeases {
	// 	removedLeaseDtos[i] = lease.ToDto()
	// }

	// Update lease stores
	for _, lease := range possibleRemovedLeases {
		delete(lastLeases, lease.IPAddress)
		delete(newLeases, lease.IPAddress)

		var _, existsInRemoved = removedLeases[lease.IPAddress]

		if !existsInRemoved {
			removedLeases[lease.IPAddress] = Lease4Tracker{lease, time.Now()}
		}
	}

	for _, lease := range possibleNewLeases {
		lastLeases[lease.IPAddress] = lease

		var _, exists = newLeases[lease.IPAddress]
		if !exists {
			newLeases[lease.IPAddress] = Lease4Tracker{lease, time.Now()}
		}
	}

	for _, lease := range removedLeases {
		removedLeaseDtos = append(removedLeaseDtos, lease.lease.ToDto())
	}

	// Clean up timed out removed leases
	timeNowUnix := time.Now().Unix()
	for _, lease := range removedLeases {
		if timeNowUnix-lease.addedAt.Unix() > removedTtlSec {
			delete(removedLeases, lease.lease.IPAddress)
		}
	}

	var result = LeasesResponse{
		Leases:        resultLeases,
		NewLeases:     newLeaseDtos,
		RemovedLeases: removedLeaseDtos,
	}

	return result
}
