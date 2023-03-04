package cassandra

import (
	"context"
	"errors"
	"sync"
	"time"

	gocql "github.com/gocql/gocql"
	"github.com/hailocab/go-hostpool"
	gocqlx "github.com/scylladb/gocqlx/v2"
)

const (
	idleTimeout = 10 * time.Minute
	pingPeriod  = 1 * time.Second
)

// cassandraSession is a wrapper around gocqlx.Session with additional fields
// for managing the session's state.
type cassandraSession struct {
	session  *gocqlx.Session
	keyspace string
	lastUsed time.Time
	mu       sync.Mutex
}

// CassandraPool is a thread-safe connection pool for Cassandra sessions.
type CassandraPool struct {
	sessions   map[string]*cassandraSession
	poolSize   int
	poolConfig *gocql.PoolConfig
	mu         sync.Mutex
}

// Singleton instance of the CassandraPool
var pool *CassandraPool
var once sync.Once

// GetCassandraPoolInstance returns the singleton instance of the CassandraPool.
func GetCassandraPoolInstance() *CassandraPool {
	once.Do(func() {
		pool = &CassandraPool{
			sessions: make(map[string]*cassandraSession),
			poolSize: 5,
			poolConfig: &gocql.PoolConfig{
				HostSelectionPolicy: gocql.HostPoolHostPolicy(
					hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),
				),
			},
		}
		go pool.pinger()
	})
	return pool
}

// GetSession returns a session for the specified keyspace.
func (p *CassandraPool) GetSession(ctx context.Context, keyspace string) (*gocqlx.Session, error) {
	// Try to get an existing session from the pool.
	session, err := p.getSession(keyspace)
	if err != nil {
		// If no session is available, create a new one and add it to the pool.
		session, err = p.createSession(ctx, keyspace)
		if err != nil {
			p.removeSession(keyspace)
			return nil, err
		}
		p.addSession(session)
	}
	session.markUsed()
	return session.session, nil
}

// getSession returns an existing session from the pool if one is available.
func (p *CassandraPool) getSession(keyspace string) (*cassandraSession, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	session, ok := p.sessions[keyspace]
	if !ok {
		return nil, errors.New("no session available")
	}
	return session, nil
}

// createSession creates a new session for the specified keyspace.
func (p *CassandraPool) createSession(ctx context.Context, keyspace string) (*cassandraSession, error) {
	cluster, err := New(NewCassandraConfig(keyspace))
	if err != nil {
		return nil, err
	}
	cluster.NumConns = p.poolSize

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		session.Close()
		return nil, err
	}
	if err := session.Query("SELECT now() FROM system.local", nil).Exec(); err != nil {
		session.Close()
		return nil, err
	}
	return &cassandraSession{
		session:  &session,
		keyspace: keyspace,
		lastUsed: time.Now(),
	}, nil
}

// addSession adds a session to the pool.
func (p *CassandraPool) addSession(session *cassandraSession) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sessions[session.keyspace] = session
}

// removeSession removes a session from the pool.
func (p *CassandraPool) removeSession(keyspace string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.sessions, keyspace)
}

// markUsed updates the last used time for the session.
func (s *cassandraSession) markUsed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastUsed = time.Now()
}

func (s *cassandraSession) isLocked() bool {
	return s.mu.TryLock()
}

// pinger periodically checks for idle sessions and closes them.
func (p *CassandraPool) pinger() {
	for {
		p.mu.Lock()
		for keyspace, session := range p.sessions {
			session.mu.Lock()
			if time.Since(session.lastUsed) > idleTimeout {
				session.session.Close()
				delete(p.sessions, keyspace)
			}
			session.mu.Unlock()
		}
		p.mu.Unlock()
		time.Sleep(pingPeriod)
	}
}
