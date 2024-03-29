package cassandra

import (
	"strconv"
	"time"

	gocql "github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

// New cassandra session and return Cassandra struct
func New(conf *CassandraConfig) (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(conf.Host, conf.Host, conf.Host)
	cluster.Keyspace = conf.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.Username,
		Password: conf.Password,
	}
	if conf.Port != "" {
		port, err := strconv.Atoi(conf.Port)
		if err != nil {
			log.Error("Error parsing cassandra port")
			return nil, err
		}
		cluster.Port = port
	}
	if conf.TlsConfig != nil {
		cluster.SslOpts = &gocql.SslOptions{
			Config: conf.TlsConfig,
		}
	}
	// Set timeout to 30 seconds
	cluster.Timeout = 30 * time.Second
	// Use a retry policy
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	return cluster, nil
}
