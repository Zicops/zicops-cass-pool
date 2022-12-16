package cassandra

import (
	gocql "github.com/gocql/gocql"
	"github.com/hailocab/go-hostpool"
	gocqlx "github.com/scylladb/gocqlx/v2"
)

// GlobalSession is a map of cassandra sessions
var GlobalSession = make(map[string]*gocqlx.Session)

func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	cluster, err := New(NewCassandraConfig(keyspace))
	if err != nil {
		return nil, err
	}
	cluster.NumConns = 5
	cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(
		hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),
	)
	if GlobalSession[keyspace] == nil || GlobalSession[keyspace].Closed() {

		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return nil, err
		}
		GlobalSession[keyspace] = &session
	} else if GlobalSession[keyspace].Query("SELECT now() FROM system.local", nil).Exec() != nil {
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return nil, err
		}
		return &session, nil
	}
	return GlobalSession[keyspace], nil
}
