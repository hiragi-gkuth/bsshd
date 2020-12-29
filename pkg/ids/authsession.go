package ids

import "sync"

// KVS is string-AuthInfo key-value pair
type KVS map[string]AuthInfo

// ReliableKVS is reliable of use across goroutines
type ReliableKVS struct {
	sync.Mutex
	dictionary KVS
}

var r *ReliableKVS = NewKVS()

// AuthSession is session store for authentication
type AuthSession struct {
	r *ReliableKVS
}

// GetAuthSession returns auth session
func GetAuthSession() *AuthSession {
	if r == nil {
		r = NewKVS()
	}
	return &AuthSession{
		r,
	}
}

// NewKVS returns new kvs
func NewKVS() *ReliableKVS {
	return &ReliableKVS{
		dictionary: make(KVS),
	}
}

// Set sets value for kvs
func (as *AuthSession) Set(key string, value AuthInfo) {
	as.r.Lock()
	as.r.dictionary[key] = value
	as.r.Unlock()
}

// Get gets value of kvs
func (as *AuthSession) Get(key string) (AuthInfo, bool) {
	as.r.Lock()
	v, ok := as.r.dictionary[key]
	as.r.Unlock()
	return v, ok
}
