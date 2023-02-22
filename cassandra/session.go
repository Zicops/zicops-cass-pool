package cassandra

import gocqlx "github.com/scylladb/gocqlx/v2"

var (
	manager *sessManager
)

func init() {
	// Initialize session manager with default configuration
	sm, err := NewSessionManager(NewCassandraConfig(""))
	if err != nil {
		panic(err)
	}

	manager = sm
}

// GetCassSession returns a session for the given keyspace.
func GetCassSession(keyspace string) (*gocqlx.Session, error) {
	return manager.GetCassSession(keyspace)
}
