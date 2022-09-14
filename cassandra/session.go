package cassandra

import gocqlx "github.com/scylladb/gocqlx/v2"

var GlobalSession *gocqlx.Session

func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	if GlobalSession == nil || GlobalSession.Closed() {
		cluster, err := New(NewCassandraConfig(keyspace))
		if err != nil {
			return nil, err
		}
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return nil, err
		}
		GlobalSession = &session
	}
	return GlobalSession, nil
}
