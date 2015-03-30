package sharedmap

type KeyType interface{}
type ValueType interface{}

// A concurrency-safe Map implementation based on
// synchronization through channel I/O instead of
// mutexes
type SharedMap struct {
	cGetItem chan itemYield
	cAddItem chan itemAdd
	cDelItem chan itemDel
	cMapSize chan chan int
}

type itemYield struct {
	key    KeyType
	cValue chan ValueType
}

type itemAdd struct {
	key       KeyType
	val       ValueType
	cModified chan bool
}

type itemDel struct {
	key      KeyType
	cSuccess chan bool
}

// Gets the value corresponding to the given key from the SharedMap
func (s *SharedMap) Get(key KeyType) ValueType {
	g := itemYield{key, make(chan ValueType, 1)}
	s.cGetItem <- g
	return <-g.cValue
}

// Adds/updates the given key-value pair to the SharedMap and returns
// true if the key already existed
func (s *SharedMap) Add(key KeyType, value ValueType) bool {
	a := itemAdd{key, value, make(chan bool, 1)}
	s.cAddItem <- a
	return <-a.cModified
}

// Deletes the given key from the SharedMap and returns true if the
// key was found and deleted
func (s *SharedMap) Delete(key KeyType) bool {
	d := itemDel{key, make(chan bool, 1)}
	s.cDelItem <- d
	return <-d.cSuccess
}

// Returns the number of elements stored in the map
func (s *SharedMap) Size() int {
	rchan := make(chan int, 1)
	s.cMapSize <- rchan
	return <-rchan
}

// Goroutine to handle the internal state of SharedMap
func (s *SharedMap) MapManager() {
	m := map[KeyType]ValueType{}
	for {
		select {
		case info := <-s.cAddItem:
			_, ok := m[info.key]
			m[info.key] = info.val
			info.cModified <- ok
		case info := <-s.cGetItem:
			info.cValue <- m[info.key]
		case info := <-s.cDelItem:
			_, ok := m[info.key]
			delete(m, info.key)
			info.cSuccess <- ok
		case rchan := <-s.cMapSize:
			rchan <- len(m)
		}
	}
}

// Creates a new SharedMap
func NewSharedMap() *SharedMap {
	s := &SharedMap{
		make(chan itemYield),
		make(chan itemAdd),
		make(chan itemDel),
		make(chan chan int),
	}
	go s.MapManager()
	return s
}
