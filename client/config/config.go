package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"log"
)

var (
	caCertStr  = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNWakNDQWYyZ0F3SUJBZ0lVV25OWElTcDIy\nYTZ4S0w3ak95bTFqbS80VkNvd0NnWUlLb1pJemowRUF3TXcKZ1lBeEN6QUpCZ05WQkFZVEFrTkJN\nUXN3Q1FZRFZRUUlEQUpCUWpFUU1BNEdBMVVFQnd3SFEyRnNaMkZ5ZVRFTApNQWtHQTFVRUNnd0NS\nRkF4Q3pBSkJnTlZCQXNNQWtSUU1SQXdEZ1lEVlFRRERBZGpZUzVrWlcxdk1TWXdKQVlKCktvWklo\ndmNOQVFrQkZoZHdZVzVqYUdGc1pHRjRaWE5vUUdkdFlXbHNMbU52YlRBZUZ3MHlNVEF6TURReE9E\nTTIKTkRCYUZ3MHpNVEF6TURJeE9ETTJOREJhTUlHQU1Rc3dDUVlEVlFRR0V3SkRRVEVMTUFrR0Ex\nVUVDQXdDUVVJeApFREFPQmdOVkJBY01CME5oYkdkaGNua3hDekFKQmdOVkJBb01Ba1JRTVFzd0NR\nWURWUVFMREFKRVVERVFNQTRHCkExVUVBd3dIWTJFdVpHVnRiekVtTUNRR0NTcUdTSWIzRFFFSkFS\nWVhjR0Z1WTJoaGJHUmhlR1Z6YUVCbmJXRnAKYkM1amIyMHdXVEFUQmdjcWhrak9QUUlCQmdncWhr\nak9QUU1CQndOQ0FBUlhySCtHYVZieTg1MitnR2FkMnNwTApVaTNuQjd0dGNMa0pUTGVUNkI0RDNN\nOThQaGNNL3A4Q1RtNHVJVGdqQ2lXcXArLzZwSU8zQk42cTltN2h1TEw4Cm8xTXdVVEFkQmdOVkhR\nNEVGZ1FVdkwwVEtkUXFiSC82WitxL3MrajMxOVd1VWhBd0h3WURWUjBqQkJnd0ZvQVUKdkwwVEtk\nUXFiSC82WitxL3MrajMxOVd1VWhBd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBS0JnZ3Foa2pPUFFR\nRApBd05IQURCRUFpQjdGbW1MZGFOQ2J3YlprUU9ZNWEzTWl4T0pYMTJsL0J3SktNSWdIRDc4SXdJ\nZ09GMW1GYTNMClN5ZExuMHNqZHdQWnFGU0ZZLzVGdmJ3bTBSeEFoSFAxTXZZPQotLS0tLUVORCBD\nRVJUSUZJQ0FURS0tLS0tCg=="
	serverName = "localhost"
	clientKey  = "LS0tLS1CRUdJTiBFQyBQQVJBTUVURVJTLS0tLS0KQmdncWhrak9QUU1CQnc9PQotLS0tLUVORCBF\nQyBQQVJBTUVURVJTLS0tLS0KLS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVF\nSU04OStBVGdpQytwck9rRThJeTRXRTJnQjVKM1A5emlHYlVxZHR6TUhJRzVvQW9HQ0NxR1NNNDkK\nQXdFSG9VUURRZ0FFbkVVUVZaY3lvaWtRSG5LaUFra3dlbjRXUmd5dityL2RZNGQ4ZFRRSnRKSzJG\nazRMT204RApRcVhVNmRVQVc3UStsTVF1T2JxVGxtZElET0w4TURteVpBPT0KLS0tLS1FTkQgRUMg\nUFJJVkFURSBLRVktLS0tLQo="
	clientCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLRENDQWMyZ0F3SUJBZ0lVZUdlblNNRk9K\nbEVlZ2N6UVRKQjg0R3k2YWJJd0NnWUlLb1pJemowRUF3TXcKZ1lBeEN6QUpCZ05WQkFZVEFrTkJN\nUXN3Q1FZRFZRUUlEQUpCUWpFUU1BNEdBMVVFQnd3SFEyRnNaMkZ5ZVRFTApNQWtHQTFVRUNnd0NS\nRkF4Q3pBSkJnTlZCQXNNQWtSUU1SQXdEZ1lEVlFRRERBZGpZUzVrWlcxdk1TWXdKQVlKCktvWklo\ndmNOQVFrQkZoZHdZVzVqYUdGc1pHRjRaWE5vUUdkdFlXbHNMbU52YlRBZUZ3MHlNVEF6TURrd01q\nVXgKTVRkYUZ3MHlNakF6TURrd01qVXhNVGRhTUVveEN6QUpCZ05WQkFZVEFrTkJNUXN3Q1FZRFZR\nUUlEQUpCUWpFTQpNQW9HQTFVRUJ3d0RRMEZNTVFzd0NRWURWUVFLREFKRVVERVRNQkVHQTFVRUF3\nd0tjbVZoWkdOc2FXVnVkREJaCk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkp4RkVG\nV1hNcUlwRUI1eW9nSkpNSHArRmtZTXIvcS8KM1dPSGZIVTBDYlNTdGhaT0N6cHZBMEtsMU9uVkFG\ndTBQcFRFTGptNms1Wm5TQXppL0RBNXNtU2pXakJZTUI4RwpBMVVkSXdRWU1CYUFGTHk5RXluVUtt\neC8rbWZxdjdQbzk5ZlZybElRTUFrR0ExVWRFd1FDTUFBd0N3WURWUjBQCkJBUURBZ1R3TUIwR0Ex\nVWRFUVFXTUJTQ0NXeHZZMkZzYUc5emRJSUhZWEJwWkdWdGJ6QUtCZ2dxaGtqT1BRUUQKQXdOSkFE\nQkdBaUVBcHpVYm5BYmthV2x4cHZhTTNzWUNva1llWUpJeHFJSU5QQ1YwUzkwbzhIb0NJUUNESXFJ\nOQpBU0VlQ0RoZlg5cDJidUgycUpPVTQwdXovbVN3cVZ5czBCR2dqZz09Ci0tLS0tRU5EIENFUlRJ\nRklDQVRFLS0tLS0K"
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
	config := &Config{
		ServerName: serverName,
		TLSKey:     clientKey,
		TLSCert:    clientCert,
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
		Certificates: []tls.Certificate{c.TLSCertificate},
		RootCAs:      c.RootCAs,
		ServerName:   c.ServerName,
	}
}
