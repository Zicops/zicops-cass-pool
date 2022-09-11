package cassandra

import gocql "github.com/gocql/gocql"

func GetCassSession() (*gocql.Session, error) {
	cluster, err := New(NewCassandraConfig())
	if err != nil {
		return nil, err
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
