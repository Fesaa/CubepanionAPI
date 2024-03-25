package utils

import "sync"

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Remove(key K)

	Contains(key K) bool
	ForEach(func(key K, value V))
}

type DoubleMap[K, V comparable] interface {
	Map[K, V]
	RemoveByValue(value V)
}

type mapImpl[K comparable, V any] struct {
	_map  map[K]V
	mutex *sync.RWMutex
}

func NewMap[K comparable, V any]() Map[K, V] {
	return &mapImpl[K, V]{
		_map:  make(map[K]V),
		mutex: &sync.RWMutex{},
	}
}

type doubleMapImpl[K, V comparable] struct {
	mapImpl[K, V]
}

func NewDoubleMap[K, V comparable]() DoubleMap[K, V] {
	return &doubleMapImpl[K, V]{
		mapImpl: mapImpl[K, V]{
			_map:  make(map[K]V),
			mutex: &sync.RWMutex{},
		},
	}
}

func (m *mapImpl[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m._map[key]
	return value, ok
}

func (m *mapImpl[K, V]) Put(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m._map[key] = value
}

func (m *mapImpl[K, V]) Remove(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m._map, key)
}

func (m *mapImpl[K, V]) Contains(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, ok := m._map[key]
	return ok
}

func (m *doubleMapImpl[K, V]) RemoveByValue(value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key, v := range m._map {
		if value == v {
			delete(m._map, key)
		}
	}
}

func (m *mapImpl[K, V]) ForEach(f func(key K, value V)) {
	for key, value := range m._map {
		f(key, value)
	}
}
