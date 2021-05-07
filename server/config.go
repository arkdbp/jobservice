package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"log"
)

var (
	defaultPort = 12000
	serverName  = "localhost"
	server2c    = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNQVENDQWVTZ0F3SUJBZ0lVZUdlblNNRk9K\nbEVlZ2N6UVRKQjg0R3k2YWJFd0NnWUlLb1pJemowRUF3TXcKZ1lBeEN6QUpCZ05WQkFZVEFrTkJN\nUXN3Q1FZRFZRUUlEQUpCUWpFUU1BNEdBMVVFQnd3SFEyRnNaMkZ5ZVRFTApNQWtHQTFVRUNnd0NS\nRkF4Q3pBSkJnTlZCQXNNQWtSUU1SQXdEZ1lEVlFRRERBZGpZUzVrWlcxdk1TWXdKQVlKCktvWklo\ndmNOQVFrQkZoZHdZVzVqYUdGc1pHRjRaWE5vUUdkdFlXbHNMbU52YlRBZUZ3MHlNVEF6TURReE9U\nRXgKTkRCYUZ3MHlNakF6TURReE9URXhOREJhTUdFeEN6QUpCZ05WQkFZVEFrRlZNUk13RVFZRFZR\nUUlEQXBUYjIxbApMVk4wWVhSbE1TRXdId1lEVlFRS0RCaEpiblJsY201bGRDQlhhV1JuYVhSeklG\nQjBlU0JNZEdReEdqQVlCZ05WCkJBTU1FV0Z3YVdSbGJXOHNiRzlqWVd4b2IzTjBNRmt3RXdZSEtv\nWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUUKTWpCb00rMXY4NUFHK2pxMlFveXN3ZWxyYjVTeFN3\nSzVMcnlJSFQ5MUd5aHZ2T3greHpCcU9yS25NU2FaZDMyNgpyM0QzNEZLRkR3L0NvS0o4Zm5MMU5h\nTmFNRmd3SHdZRFZSMGpCQmd3Rm9BVXZMMFRLZFFxYkgvNlorcS9zK2ozCjE5V3VVaEF3Q1FZRFZS\nMFRCQUl3QURBTEJnTlZIUThFQkFNQ0JQQXdIUVlEVlIwUkJCWXdGSUlKYkc5allXeG8KYjNOMGdn\nZGhjR2xrWlcxdk1Bb0dDQ3FHU000OUJBTURBMGNBTUVRQ0lDWW1mMldWZGIwTXcxOTB5MkJzYTZD\nNgpjWjJyV05MYnVsSkl1TFpTMDdTN0FpQnJITVZURmNNc2VMT1I0c3M3YXlRMCszL2JIMXFHZUxI\najc0amdzWU43CnZBPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	server2k    = "LS0tLS1CRUdJTiBFQyBQQVJBTUVURVJTLS0tLS0KQmdncWhrak9QUU1CQnc9PQotLS0tLUVORCBF\nQyBQQVJBTUVURVJTLS0tLS0KLS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVF\nSUtVb3kyaDhYMlJqaUlROXJraFFMT3hhNEliQkdmczFmODk4YUJPR0gwdmFvQW9HQ0NxR1NNNDkK\nQXdFSG9VUURRZ0FFTWpCb00rMXY4NUFHK2pxMlFveXN3ZWxyYjVTeFN3SzVMcnlJSFQ5MUd5aHZ2\nT3greHpCcQpPcktuTVNhWmQzMjZyM0QzNEZLRkR3L0NvS0o4Zm5MMU5RPT0KLS0tLS1FTkQgRUMg\nUFJJVkFURSBLRVktLS0tLQo="
	caCertStr   = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNWakNDQWYyZ0F3SUJBZ0lVV25OWElTcDIy\nYTZ4S0w3ak95bTFqbS80VkNvd0NnWUlLb1pJemowRUF3TXcKZ1lBeEN6QUpCZ05WQkFZVEFrTkJN\nUXN3Q1FZRFZRUUlEQUpCUWpFUU1BNEdBMVVFQnd3SFEyRnNaMkZ5ZVRFTApNQWtHQTFVRUNnd0NS\nRkF4Q3pBSkJnTlZCQXNNQWtSUU1SQXdEZ1lEVlFRRERBZGpZUzVrWlcxdk1TWXdKQVlKCktvWklo\ndmNOQVFrQkZoZHdZVzVqYUdGc1pHRjRaWE5vUUdkdFlXbHNMbU52YlRBZUZ3MHlNVEF6TURReE9E\nTTIKTkRCYUZ3MHpNVEF6TURJeE9ETTJOREJhTUlHQU1Rc3dDUVlEVlFRR0V3SkRRVEVMTUFrR0Ex\nVUVDQXdDUVVJeApFREFPQmdOVkJBY01CME5oYkdkaGNua3hDekFKQmdOVkJBb01Ba1JRTVFzd0NR\nWURWUVFMREFKRVVERVFNQTRHCkExVUVBd3dIWTJFdVpHVnRiekVtTUNRR0NTcUdTSWIzRFFFSkFS\nWVhjR0Z1WTJoaGJHUmhlR1Z6YUVCbmJXRnAKYkM1amIyMHdXVEFUQmdjcWhrak9QUUlCQmdncWhr\nak9QUU1CQndOQ0FBUlhySCtHYVZieTg1MitnR2FkMnNwTApVaTNuQjd0dGNMa0pUTGVUNkI0RDNN\nOThQaGNNL3A4Q1RtNHVJVGdqQ2lXcXArLzZwSU8zQk42cTltN2h1TEw4Cm8xTXdVVEFkQmdOVkhR\nNEVGZ1FVdkwwVEtkUXFiSC82WitxL3MrajMxOVd1VWhBd0h3WURWUjBqQkJnd0ZvQVUKdkwwVEtk\nUXFiSC82WitxL3MrajMxOVd1VWhBd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBS0JnZ3Foa2pPUFFR\nRApBd05IQURCRUFpQjdGbW1MZGFOQ2J3YlprUU9ZNWEzTWl4T0pYMTJsL0J3SktNSWdIRDc4SXdJ\nZ09GMW1GYTNMClN5ZExuMHNqZHdQWnFGU0ZZLzVGdmJ3bTBSeEFoSFAxTXZZPQotLS0tLUVORCBD\nRVJUSUZJQ0FURS0tLS0tCg=="
)

type Config struct {
	Port           int
	ServerName     string
	TLSKey         string
	TLSCert        string
	CACert         string
	TLSCertificate tls.Certificate
	RootCAs        *x509.CertPool
}

// ProvideConfig  todo read from envs and override defaults
func ProvideConfig() *Config {
	config := &Config{Port: defaultPort,
		ServerName: serverName,
		TLSKey:     server2k,
		TLSCert:    server2c,
		CACert:     caCertStr,
	}
	err := config.SetCertPool()
	if err != nil {
		log.Fatal("failed to set certPool with error", err)
	}
	err = config.SetTlsCertificate()
	if err != nil {
		log.Fatal("failed to set TLSCertificate with error", err)
	}
	return config
}

func (c *Config) SetTlsCertificate() error {
	tlsKey, err := base64.StdEncoding.DecodeString(c.TLSKey)
	if err != nil {
		return err
	}
	tlsCert, err := base64.StdEncoding.DecodeString(c.TLSCert)
	if err != nil {
		return err
	}
	keyPair, err := tls.X509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return err
	}
	c.TLSCertificate = keyPair
	return nil
}

func (c *Config) SetCertPool() error {
	caCert, err := base64.StdEncoding.DecodeString(c.CACert)
	if err != nil {
		return err
	}
	RootCAs := x509.NewCertPool()
	if ok := RootCAs.AppendCertsFromPEM(caCert); !ok {
		return err
	}
	c.RootCAs = RootCAs
	return nil
}

func (c *Config) TLSConfig() *tls.Config {
	return &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		Certificates:             []tls.Certificate{c.TLSCertificate},
		ClientCAs:                c.RootCAs,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS13,
	}
}
