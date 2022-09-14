package cassandra

import gocqlx "github.com/scylladb/gocqlx/v2"

// GlobalSession is a map of cassandra sessions
var GlobalSession = make(map[string]*gocqlx.Session)

func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	if GlobalSession[keyspace] == nil || GlobalSession[keyspace].Closed() {
		cluster, err := New(NewCassandraConfig(keyspace))
		if err != nil {
			return nil, err
		}
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return nil, err
		}
		GlobalSession[keyspace] = &session
	}
	return GlobalSession[keyspace], nil
}
