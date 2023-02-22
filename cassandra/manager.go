package cassandra

import (
	"strconv"

	"github.com/gocql/gocql"
	"github.com/hailocab/go-hostpool"
	gocqlx "github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
)

// sessionManager is a simple session manager that stores sessions in a map.
type sessManager struct {
	sessions map[string]*gocqlx.Session
	cluster  *gocql.ClusterConfig
}

// GetCassSession returns a session for the given keyspace.
func (sm *sessManager) GetCassSession(keyspace string) (*gocqlx.Session, error) {
	if session, ok := sm.sessions[keyspace]; ok {
		return session, nil
	}
	if keyspace == "" {
		keyspace = "userz"
	}
	sm.cluster.Keyspace = keyspace
	session, err := gocqlx.WrapSession(sm.cluster.CreateSession())
	if err != nil {
		return nil, err
	}

	sm.sessions[keyspace] = &session

	return &session, nil
}

// NewSessionManager creates a new session manager with the given config.
func NewSessionManager(config *CassandraConfig) (*sessManager, error) {
	cluster := gocql.NewCluster(config.Host, config.Host, config.Host)
	cluster.Keyspace = config.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	if config.Port != "" {
		port, err := strconv.Atoi(config.Port)
		if err != nil {
			log.Error("Error parsing cassandra port")
			return nil, err
		}
		cluster.Port = port
	}
	if config.TlsConfig != nil {
		cluster.SslOpts = &gocql.SslOptions{
			Config: config.TlsConfig,
		}
	}
	cluster.NumConns = 10
	cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(
		hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),
	)

	return &sessManager{
		sessions: make(map[string]*gocqlx.Session),
		cluster:  cluster,
	}, nil
}
