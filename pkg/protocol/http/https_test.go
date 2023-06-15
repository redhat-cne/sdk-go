package http_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/redhat-cne/sdk-go/pkg/channel"
	cneevent "github.com/redhat-cne/sdk-go/pkg/event"
	ceHttp "github.com/redhat-cne/sdk-go/pkg/protocol/http"
	"github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	storePathTLS      = "./tls"
	subscriptionTwoID = "123e4567-e89b-12d3-a456-426614174002"
	serverAddressTLS  = types.ParseURI("https://127.0.0.1:8049")
	clientAddressTLS  = types.ParseURI("https://127.0.0.1:8047")
	hostPortTLS       = 8049
	clientPortTLS     = 8047

	clientClientIDTLS = func(serviceName string) uuid.UUID {
		var namespace = uuid.NameSpaceURL
		var url = []byte(serviceName)
		return uuid.NewMD5(namespace, url)
	}(clientAddressTLS.String())

	subscriptionTwo = &pubsub.PubSub{
		ID:       subscriptionTwoID,
		Resource: "/test/test/2",
	}
)

const (
	caFile         = "/tmp/ca.crt"
	serverCertFile = "/tmp/server.crt"
	serverKeyFile  = "/tmp/server.key"
	clientCertFile = "/tmp/client.crt"
	clientKeyFile  = "/tmp/client.key"
)

func TestTLSSubscribeCreated(t *testing.T) {
	in := make(chan *channel.DataChan, 10)
	out := make(chan *channel.DataChan, 10)
	closeCh := make(chan struct{})
	eventChannel := make(chan *channel.DataChan, 10)
	serverTLS, clientTLS, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, nil, nil, true)
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	err = server.Start(&wg)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	var clientS *ceHttp.Server
	tls, err = InitializeTLSConfig(clientTLS.TLSCertFile, clientTLS.TLSKeyFile, clientTLS.CABundleFile)
	assert.Nil(t, err)
	go createTLSClient(t, clientS, tls, closeCh, eventChannel)
	time.Sleep(500 * time.Millisecond)
	<-out
	assert.Equal(t, 1, len(server.Sender))
	d := <-eventChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)
	close(closeCh)
}
func TestTLSSendEvent(t *testing.T) {
	time.Sleep(2 * time.Second)
	e := CloudEvents()
	in := make(chan *channel.DataChan, 10)
	out := make(chan *channel.DataChan, 10)
	clientOutChannel := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	serverTLS, clientTLS, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, nil, nil, true)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	err = server.Start(&wg)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	time.Sleep(500 * time.Millisecond)
	var clientS *ceHttp.Server
	tls, err = InitializeTLSConfig(clientTLS.TLSCertFile, clientTLS.TLSKeyFile, clientTLS.CABundleFile)
	assert.Nil(t, err)
	go createTLSClient(t, clientS, tls, closeCh, clientOutChannel)
	//  read what server has in outChannel
	<-out
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, len(server.Sender))
	// read what client put in out channel
	d := <-clientOutChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	// send event
	in <- &channel.DataChan{
		Address: subscriptionTwo.Resource,
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}
	// read event
	log.Info("waiting for event channel from the client when it received the event")
	d = <-clientOutChannel // a client needs to break out or else it will be holding it forever
	assert.Equal(t, channel.EVENT, d.Type)
	dd := cneevent.Data{}
	err = json.Unmarshal(e.Data(), &dd)
	assert.Nil(t, err)
	assert.Equal(t, dd.Version, "1.0")

	log.Info("waiting for event response")
	d = <-out
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	time.Sleep(250 * time.Millisecond)
	close(closeCh)
}

func TestTLSSendSuccess(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	clientOutChannel := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	serverTLS, clientTLS, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, func(e cloudevents.Event, dataChan *channel.DataChan) error {
		dataChan.Address = clientAddress.String()
		e.SetType(channel.EVENT.String())
		if err = ceHttp.Post(fmt.Sprintf("%s/event", clientAddress), e); err != nil {
			log.Errorf("error %s sending event %v at  %s", err, e, clientAddress)
			return err
		}
		return nil
	}, nil, true)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	// Start the server and channel processor
	err = server.Start(&wg)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)

	// create a sender
	var clientS *ceHttp.Server
	tls, err = InitializeTLSConfig(clientTLS.TLSCertFile, clientTLS.TLSKeyFile, clientTLS.CABundleFile)
	assert.Nil(t, err)
	go createTLSClient(t, clientS, tls, closeCh, clientOutChannel)
	time.Sleep(500 * time.Millisecond)
	<-out
	assert.Equal(t, 1, len(server.Sender))
	close(closeCh)
	//waitTimeout(&wg, timeout)
}

func TestTLSHealth(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	var status int
	var urlErr error
	serverTLS, _, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, nil, nil, true)
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	err = server.Start(&wg)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	time.Sleep(2 * time.Second)
	status, urlErr = ceHttp.Get(fmt.Sprintf("%s/health", serverAddressTLS.String()))
	assert.Nil(t, urlErr)
	assert.Equal(t, http.StatusOK, status)
	close(closeCh)
}

func TestTLSSender(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})

	serverTLS, _, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, nil, nil, true)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	// Start the server and channel processor
	err = server.Start(&wg)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	time.Sleep(2 * time.Second)
	err = server.NewSender(serverClientID, serverAddressTLS.String())
	assert.Nil(t, err)
	sender := server.GetSender(serverClientID, ceHttp.HEALTH)
	assert.NotNil(t, sender)
	e := CloudEvents()
	err = sender.Send(e)
	assert.Nil(t, err)
	close(closeCh)
}

func TestTLSStatusWithSubscription(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	eventChannel := make(chan *channel.DataChan, 10)
	closeCh := make(chan struct{})
	onStatusReceiveOverrideFn := func(e event.Event, d *channel.DataChan) error {
		ce := CloudEvents()
		d.Data = &ce
		return nil
	}
	serverTLS, clientTLS, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, onStatusReceiveOverrideFn, nil, true)
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	// Start the server and channel processor
	err = server.Start(&wg)
	server.RegisterPublishers(serverAddress)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	time.Sleep(2 * time.Second)
	// create client and create subscription
	var clientS *ceHttp.Server
	tls, err = InitializeTLSConfig(clientTLS.TLSCertFile, clientTLS.TLSKeyFile, clientTLS.CABundleFile)
	assert.Nil(t, err)
	go createTLSClient(t, clientS, tls, closeCh, eventChannel)

	time.Sleep(500 * time.Millisecond)
	<-out
	assert.Equal(t, 1, len(server.Sender))
	d := <-eventChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	transport := http.Transport{
		TLSClientConfig:     tls.Clone(),
		MaxIdleConnsPerHost: 20,
	}
	// send status ping
	hClient := &http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}
	requestURL := fmt.Sprintf("%s/%s/%s/CurrentState", serverAddressTLS.String(), subscriptionTwo.Resource, clientClientIDTLS)
	log.Printf(requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Nil(t, err)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	ce := cloudevents.Event{}
	err = json.Unmarshal(bodyBytes, &ce)
	log.Info(string(bodyBytes))
	if e, ok := err.(*json.SyntaxError); ok {
		log.Infof("syntax error at byte offset %d", e.Offset)
	}
	assert.Nil(t, err)

	close(closeCh)
}

func TestTLSStatusWithOutSubscription(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	onStatusReceiveOverrideFn := func(e event.Event, d *channel.DataChan) error {
		ce := CloudEvents()
		d.Data = &ce
		return nil
	}
	serverTLS, clientTLS, err := makeTLSConfig()
	assert.Nil(t, err)
	tls, err := InitializeTLSConfig(serverTLS.TLSCertFile, serverTLS.TLSKeyFile, serverTLS.CABundleFile)
	assert.Nil(t, err)
	server, err := ceHttp.InitServer(serverAddressTLS.String(), hostPortTLS, tls, storePathTLS, in, out, closeCh, onStatusReceiveOverrideFn, nil, true)
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	// Start the server and channel processor
	err = server.Start(&wg)
	server.RegisterPublishers(serverAddressTLS)
	assert.Nil(t, err)
	server.HTTPProcessor(&wg)
	time.Sleep(2 * time.Second)
	tls, err = InitializeTLSConfig(clientTLS.TLSCertFile, clientTLS.TLSKeyFile, clientTLS.CABundleFile)
	assert.Nil(t, err)
	transport := http.Transport{
		TLSClientConfig:     tls.Clone(),
		MaxIdleConnsPerHost: 20,
	}
	// send status ping
	hClient := &http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	requestURL := fmt.Sprintf("%s/%s/%s/CurrentState", serverAddressTLS.String(), subscriptionNotFound.Resource, clientClientID)
	log.Printf(requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := hClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	log.Info(string(bodyBytes))
	close(closeCh)
}

func createTLSClient(t *testing.T, clientS *ceHttp.Server, tls *tls.Config, closeCh chan struct{}, clientOutChannel chan *channel.DataChan) {
	in := make(chan *channel.DataChan, 10)
	var err error
	assert.Nil(t, clientS)
	clientS, err = ceHttp.InitServer(clientAddressTLS.String(), clientPortTLS, tls, storePathTLS, in, clientOutChannel, closeCh, nil, nil, true)
	assert.Nil(t, err)
	clientS.RegisterPublishers(serverAddressTLS)
	wg := sync.WaitGroup{}
	time.Sleep(250 * time.Millisecond)
	// Start the server and channel processor
	err = clientS.Start(&wg)
	assert.Nil(t, err)
	clientS.HTTPProcessor(&wg)
	time.Sleep(250 * time.Millisecond)
	// create a subscription
	in <- &channel.DataChan{
		ID:      subscriptionTwoID,
		Address: subscriptionTwo.Resource,
		Type:    channel.SUBSCRIBER,
	}
	time.Sleep(250 * time.Millisecond)

	<-closeCh
}

func InitializeTLSConfig(tlsCert, tlsKey, caBundle string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, err
	}
	caCert, err := os.ReadFile(caBundle)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:                caCertPool,
		RootCAs:                  caCertPool,
		Certificates:             []tls.Certificate{cert},
		ClientAuth:               tls.RequireAndVerifyClientCert, // mutual TLS
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			//tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			//tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	return tlsConfig, nil
}

// Instance TLS configuration
type InstanceTLSConfig struct {
	TLSCertFile  string
	TLSKeyFile   string
	CABundleFile string
}

// Makes a CA and generates server and client certificates
func makeTLSConfig() (server *InstanceTLSConfig, client *InstanceTLSConfig, err error) {
	if err := checkFilesExist(caFile, serverCertFile, serverKeyFile,
		clientCertFile, clientKeyFile); err != nil {
		log.Info("generating certificates")
		certSubject := pkix.Name{
			Organization:  []string{"Red Hat"},
			Country:       []string{""},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		}
		ca := &x509.Certificate{
			SerialNumber:          big.NewInt(1),
			Subject:               certSubject,
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		pub, pr, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		// CA is self-signed:
		caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, pr)
		if err != nil {
			return nil, nil, err
		}

		// pem encode
		caPEM := new(bytes.Buffer)
		err = pem.Encode(caPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		})
		if err != nil {
			return nil, nil, err
		}
		err = writePEM(caFile, caPEM)
		if err != nil {
			return nil, nil, err
		}

		caPrivKeyPEM := new(bytes.Buffer)
		byt, err := x509.MarshalPKCS8PrivateKey(pr)
		if err != nil {
			return nil, nil, err
		}

		err = pem.Encode(caPrivKeyPEM, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: byt,
		})
		if err != nil {
			return nil, nil, err
		}
		// set up server certificate
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(2),
			Subject:      certSubject,
			IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
			NotBefore:    time.Now(),
			NotAfter:     time.Now().AddDate(1, 0, 0),
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		}

		certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, pr)
		if err != nil {
			return nil, nil, err
		}

		certPEM := new(bytes.Buffer)
		err = pem.Encode(certPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		})
		if err != nil {
			return nil, nil, err
		}
		err = writePEM(serverCertFile, certPEM)
		if err != nil {
			return nil, nil, err
		}

		certPrivKeyPEM := new(bytes.Buffer)
		err = pem.Encode(certPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
		})
		if err != nil {
			return nil, nil, err
		}
		err = writePEM(serverKeyFile, certPrivKeyPEM)
		if err != nil {
			return nil, nil, err
		}
		// set up client certificate
		cert.SerialNumber = big.NewInt(3)
		certPrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, pr)
		if err != nil {
			return nil, nil, err
		}

		certPEM = new(bytes.Buffer)
		err = pem.Encode(certPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		})
		if err != nil {
			return nil, nil, err
		}
		err = writePEM(clientCertFile, certPEM)
		if err != nil {
			return nil, nil, err
		}

		certPrivKeyPEM = new(bytes.Buffer)
		err = pem.Encode(certPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
		})
		if err != nil {
			return nil, nil, err
		}
		err = writePEM(clientKeyFile, certPrivKeyPEM)
		if err != nil {
			return nil, nil, err
		}
	}
	server = &InstanceTLSConfig{

		TLSCertFile:  serverCertFile,
		TLSKeyFile:   serverKeyFile,
		CABundleFile: caFile,
	}
	client = &InstanceTLSConfig{

		TLSCertFile:  clientCertFile,
		TLSKeyFile:   clientKeyFile,
		CABundleFile: caFile,
	}
	return server, client, nil
}

func writePEM(fn string, data *bytes.Buffer) error {
	// Open the file for writing
	file, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func checkFilesExist(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			return err
		}
	}
	return nil
}

func removeFiles(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}

func TestTeardownTls(t *testing.T) {
	removeFiles(
		fmt.Sprintf("./%s.json", path.Join(storePathTLS, clientClientIDTLS.String())),
		caFile, serverCertFile, serverKeyFile,
		clientCertFile, clientKeyFile,
	)
}
