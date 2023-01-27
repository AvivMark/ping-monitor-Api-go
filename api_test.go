package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Function creates 100 demo servers for tests only
func Get100Servers() []Host {
	hostsOverload := []Host{}
	for i := 1; i < 101; i++ {
		numStr := strconv.Itoa(i)
		name := "server-" + numStr
		groupNum := strconv.Itoa(i / 10)
		newHost := Host{
			ID:       numStr,
			Hostname: name,
			HostIP:   "10.0.0." + numStr,
			IsAlive:  false,
			Group:    "Group" + groupNum,
			PingData: "",
		}
		hostsOverload = append(hostsOverload, newHost)
	}
	return hostsOverload
}

// Routes tests
func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homePage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Welcome to the HomePage!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateHost(t *testing.T) {
	host := Host{ID: "3", Hostname: "server 3", HostIP: "192.168.1.13", IsAlive: false}
	hostJson, _ := json.Marshal(host)

	req, err := http.NewRequest("POST", "/host", bytes.NewBuffer(hostJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createHost)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var returnedHost Host
	json.Unmarshal(rr.Body.Bytes(), &returnedHost)

	if returnedHost != host {
		t.Errorf("handler returned unexpected host: got %v want %v",
			returnedHost, host)
	}

	if len(Hosts) != 3 {
		t.Errorf("Hosts slice not updated with new host: got %v want %v",
			len(Hosts), 3)
	}
}