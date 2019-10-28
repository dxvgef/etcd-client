package v3

import "time"

// key的信息
type Key struct {
	ClusterID string
	Name      string
	// 如果key是个目录，则返回true
	Dir           bool
	Value         string
	CreatedIndex  uint64
	ModifiedIndex uint64
	// key的到期时间
	Expiration *time.Time
	// key的生存时间(秒)
	TTL int64
}

