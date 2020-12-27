package ids

import "sync"

// KVS is string-AuthInfo key-value pair
type KVS map[string]AuthInfo

// ReliableKVS is reliable of use across goroutines
type ReliableKVS struct {
	sync.Mutex
	dictionary KVS
}

// NewKVS returns new kvs
func NewKVS() *ReliableKVS {
	return &ReliableKVS{
		dictionary: make(KVS),
	}
}

// Set sets value for kvs
func (kvs *ReliableKVS) Set(key string, value AuthInfo) {
	kvs.Lock()
	kvs.dictionary[key] = value
	kvs.Unlock()
}

// Get gets value of kvs
func (kvs *ReliableKVS) Get(key string) (AuthInfo, bool) {
	kvs.Lock()
	v, ok := kvs.dictionary[key]
	kvs.Unlock()
	return v, ok
}

// AuthSession is authinfo session to store authtime, ip, ...etc
var AuthSession *ReliableKVS

// InitKVS initialize KVS
func InitKVS() {
	AuthSession = NewKVS()
}
