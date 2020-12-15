package ids

// KVS は，実際の実装
var KVS map[string]AuthInfo

// InitKVS initialize KVS
func InitKVS() {
	KVS = make(map[string]AuthInfo)
}
