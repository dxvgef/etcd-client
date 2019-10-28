package v2

import (
	"context"
	"time"

	etcd "go.etcd.io/etcd/client"
)

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

// 创建k/v，如果key已存在会报错
func (receiver *Client) Create(key, value string) error {
	_, err := receiver.EtcdAPI.Create(context.Background(), key, value)
	return err
}

func (receiver *Client) CreateInOrder(dir, value string, ttl time.Duration) error {
	_, err := receiver.EtcdAPI.CreateInOrder(context.Background(), dir, value, &etcd.CreateInOrderOptions{TTL: ttl})
	return err
}

// 设置key的value，如果key不存在则自动创建，否则会重新赋值
func (receiver *Client) Set(key, value string, ttl time.Duration) error {
	_, err := receiver.EtcdAPI.Set(context.Background(), key, value, &etcd.SetOptions{
		TTL:              ttl,
		Refresh:          false,
		NoValueOnSuccess: false,
	})
	return err
}

// 设置目录，如果目录不存在则自动创建
func (receiver *Client) SetDir(key string, ttl time.Duration) error {
	_, err := receiver.EtcdAPI.Set(context.Background(), key, "", &etcd.SetOptions{
		TTL: ttl,
		Dir: true,
	})
	return err
}

// 获得某个key的信息
func (receiver *Client) Get(key string) (*Key, error) {
	resp, err := receiver.EtcdAPI.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	return &Key{
		ClusterID:     resp.ClusterID,
		Name:          resp.Node.Key,
		Dir:           resp.Node.Dir,
		Value:         resp.Node.Value,
		CreatedIndex:  resp.Node.CreatedIndex,
		ModifiedIndex: resp.Node.ModifiedIndex,
		Expiration:    resp.Node.Expiration,
		TTL:           resp.Node.TTL,
	}, nil
}

// 更新某个key的value，如果key不存在则报错
func (receiver *Client) Update(key, value string) error {
	_, err := receiver.EtcdAPI.Update(context.Background(), key, value)
	return err
}

// 删除某个key
func (receiver *Client) Delete(key string) error {
	_, err := receiver.EtcdAPI.Delete(context.Background(), key, &etcd.DeleteOptions{Dir: false})
	return err
}

// 删除某个目录，如果force=true，则强制删除目录下的所有key，否则该目录下有key时不允许删除
func (receiver *Client) DeleteDir(key string, force bool) error {
	_, err := receiver.EtcdAPI.Delete(context.Background(), key, &etcd.DeleteOptions{Dir: true, Recursive: force})
	return err
}
