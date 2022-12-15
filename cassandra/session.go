package cassandra

import (
	"time"

	gocql "github.com/gocql/gocql"
	gocqlx "github.com/scylladb/gocqlx/v2"
)

// GlobalSession is a map of cassandra sessions
var GlobalSession = make(map[string]*gocqlx.Session)

func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	if GlobalSession[keyspace] == nil || GlobalSession[keyspace].Closed() {
		cluster, err := New(NewCassandraConfig(keyspace))
		if err != nil {
			return nil, err
		}
		cluster.ReconnectionPolicy = &gocql.ConstantReconnectionPolicy{MaxRetries: 10, Interval: 5 * time.Second}
		cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 10}
		cluster.NumConns = 2
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return nil, err
		}
		GlobalSession[keyspace] = &session
	}
	return GlobalSession[keyspace], nil
}
