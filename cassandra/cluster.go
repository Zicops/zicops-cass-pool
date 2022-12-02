package cassandra

import (
	"strconv"
	"time"

	gocql "github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

// New cassandra session and return Cassandra struct
func New(conf *CassandraConfig) (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(conf.Host)
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
	cluster.Timeout = time.Hour
	cluster.ConnectTimeout = time.Hour
	return cluster, nil
}
