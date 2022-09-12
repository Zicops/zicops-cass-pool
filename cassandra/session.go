package cassandra

import gocqlx "github.com/scylladb/gocqlx/v2"

func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	cluster, err := New(NewCassandraConfig(keyspace))
	if err != nil {
		return nil, err
	}
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}
	return &session, nil
}
