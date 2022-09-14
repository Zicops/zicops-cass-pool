package cassandra

import (
	"crypto/tls"
	"crypto/x509"
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
func NewCassandraConfig(keyspace string) *CassandraConfig {
	currentConfig := &CassandraConfig{
		Host:     getEnv("CASSANDRA_HOST", "127.0.0.1"),
		Port:     getEnv("CASSANDRA_PORT", "9042"),
		Username: getEnv("CASSANDRA_USERNAME", "cassandra"),
		Password: getEnv("CASSANDRA_PASSWORD", "cassandra"),
		Keyspace: keyspace,
	}
	cert := getEnv("CASSANDRA_CERT", "")
	if cert != "" {
		certPEMBlock := []byte(cert)
		key := getEnv("CASSANDRA_KEY", "")
		keyPEMBlock := []byte(key)
		ca := getEnv("CASSANDRA_CA", "")
		caPEMBlock := []byte(ca)
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
	}
	return currentConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
