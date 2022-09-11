package cassandra

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"os"
)

// struct for cassandra config
type CassandraConfig struct {
	Host      string      `yaml:"host"`
	Port      string      `yaml:"port"`
	Username  string      `yaml:"username"`
	Password  string      `yaml:"password"`
	Keyspace  string      `yaml:"keyspace"`
	TlsConfig *tls.Config `yaml:"tls_config"`
}

// initialize cassandra config struct using env variables
func NewCassandraConfig() *CassandraConfig {
	currentConfig := &CassandraConfig{
		Host:     getEnv("CASSANDRA_HOST", "127.0.0.1"),
		Port:     getEnv("CASSANDRA_PORT", "9042"),
		Username: getEnv("CASSANDRA_USERNAME", "cassandra"),
		Password: getEnv("CASSANDRA_PASSWORD", "cassandra"),
		Keyspace: getEnv("CASSANDRA_KEYSPACE", "userz"),
	}
	cert := getEnv("CASSANDRA_CERT", "")
	certPEMBlock, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		panic(err)
	}
	key := getEnv("CASSANDRA_KEY", "")
	keyPEMBlock, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(err)
	}
	ca := getEnv("CASSANDRA_CA", "")
	caPEMBlock, err := base64.StdEncoding.DecodeString(ca)
	if err != nil {
		panic(err)
	}
	certPair, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caPEMBlock)
	currentConfig.TlsConfig = &tls.Config{
		Certificates: []tls.Certificate{certPair},
		RootCAs:      caCertPool,
		ServerName:   currentConfig.Host,
	}
	return currentConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
