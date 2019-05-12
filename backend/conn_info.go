package backend

import (
	"net"
	"strconv"
	"sync"
)

type ConnInfo struct {
	ClientAddress     net.Addr
	ServerAddress     net.Addr
	BackendProcessID  uint32
	BackendSecretKey  uint32
	StartupParameters map[string]string
}

type ConnInfoStore interface {
	Find(clientAddress net.Addr, backendProcessID, backendSecretKey uint32) (*ConnInfo, error)
	Save(*ConnInfo) error
	Delete(*ConnInfo) error
}

type InMemoryConnInfoStore struct {
	store sync.Map
}

func (s *InMemoryConnInfoStore) Find(clientAddress net.Addr, backendProcessID, backendSecretKey uint32) (*ConnInfo, error) {
	key := s.key(backendProcessID, backendSecretKey)
	if v, ok := s.store.Load(key); ok {
		if v, ok := v.(*ConnInfo); ok {
			return v, nil
		}
	}
	return nil, nil
}

func (s *InMemoryConnInfoStore) Save(i *ConnInfo) error {
	key := s.key(i.BackendProcessID, i.BackendSecretKey)
	s.store.Store(key, i)
	return nil
}

func (s *InMemoryConnInfoStore) Delete(i *ConnInfo) error {
	key := s.key(i.BackendProcessID, i.BackendSecretKey)
	s.store.Delete(key)
	return nil
}

func (s *InMemoryConnInfoStore) key(processID, secretKey uint32) string {
	return strconv.Itoa(int(processID)) + ":" + strconv.Itoa(int(secretKey))
}

func NewInMemoryConnInfoStore() *InMemoryConnInfoStore {
	return &InMemoryConnInfoStore{}
}
