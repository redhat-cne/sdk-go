package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/redhat-cne/sdk-go/channel"
	"github.com/redhat-cne/sdk-go/protocol/rest-api/rest"
	"github.com/redhat-cne/sdk-go/pubsub"
	"github.com/redhat-cne/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

var (
	server     *rest.Server
	eventInCh  chan channel.DataEvent
	eventOutCh chan channel.DataEvent
	closeCh    chan bool
	wg         sync.WaitGroup
	port       int    = 8080
	apPath     string = "/api/cne/v1/"
	storePath  string = "."
)

func init() {
	eventInCh = make(chan channel.DataEvent, 10)
	eventOutCh = make(chan channel.DataEvent, 10)
	closeCh = make(chan bool)

}

func TestServer_New(t *testing.T) {

	server = rest.InitServer(port, apPath, storePath, eventOutCh)
	//start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start()
	}()

	time.Sleep(3 * time.Second)
	// this should actually send an event

	// CHECK URL IS UP
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "health"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// create subscription
	sub := pubsub.PubSub{
		ID:          "",
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: fmt.Sprintf("%s%s", apPath, "dummy")}},
		Resource:    "test/test1",
	}

	data, err := json.Marshal(&sub)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	resp.Body.Close()
	/// create new subscription
	req, err = http.NewRequest("POST", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "subscriptions"), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Print(bodyString)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	err = json.Unmarshal(bodyBytes, &sub)
	assert.Nil(t, err)

	// Get Just Created Subscription
	req, err = http.NewRequest("GET", fmt.Sprintf("%s%s%s/%s", "http://localhost:8080", apPath, "subscriptions", sub.ID), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var rSub pubsub.PubSub
	err = json.Unmarshal(bodyBytes, &rSub)
	if e, ok := err.(*json.SyntaxError); ok {
		log.Printf("syntax error at byte offset %d", e.Offset)
	}
	bodyString = string(bodyBytes)
	log.Print(bodyString)
	assert.Nil(t, err)
	assert.Equal(t, sub.ID, rSub.ID)
	resp.Body.Close()

	// Get All Subscriptions
	req, err = http.NewRequest("GET", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "subscriptions"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // Close body only if response non-nil
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var subList []pubsub.PubSub
	log.Println(string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &subList)
	assert.Nil(t, err)
	assert.Greater(t, len(subList), 0)

	//********************Publisher

	pub := pubsub.PubSub{
		ID:          "",
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: fmt.Sprintf("%s%s", apPath, "dummy")}},
		Resource:    "test/test",
	}
	pubData, err := json.Marshal(&pub)
	assert.Nil(t, err)
	assert.NotNil(t, pubData)

	req, err = http.NewRequest("POST", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "publishers"), bytes.NewBuffer(pubData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	pubBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(pubBodyBytes, &pub)
	assert.Nil(t, err)

	pubBodyString := string(pubBodyBytes)
	log.Print(pubBodyString)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Get Just created Publisher
	req, err = http.NewRequest("GET", fmt.Sprintf("%s%s%s/%s", "http://localhost:8080", apPath, "publishers", pub.ID), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	pubBodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var rPub pubsub.PubSub
	log.Printf("the data %s", string(pubBodyBytes))
	err = json.Unmarshal(pubBodyBytes, &rPub)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Nil(t, err)
	assert.Equal(t, pub.ID, rPub.ID)
	resp.Body.Close()

	// Get All Publisher
	req, err = http.NewRequest("GET", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "publishers"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	pubBodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var pubList []pubsub.PubSub
	err = json.Unmarshal(pubBodyBytes, &pubList)
	assert.Nil(t, err)
	assert.Greater(t, len(pubList), 0)

	/*
		r, _ = http.NewRequest("GET", "http://localhost:8080/api/v1/addresses", nil)
		resp, err = client.Do(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Print(bodyString)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

	*/

	// Delete All Publisher
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "publishers"), nil)
	req.Header.Set("Content-Type", "application/json")
	_, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Delete All Subscriptions
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", "http://localhost:8080", apPath, "subscriptions"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = server.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	close(eventOutCh)
	close(eventInCh)

	//wg.Wait()

}
