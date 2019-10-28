package v2

import (
	"context"
	"time"

	etcd "go.etcd.io/etcd/client"
)

type Config struct {
	// 定义一组端点，如果其中一个不可用，则客户端将尝试使用其它的端点进行操作
	// 如果曾经调用Client.Sync，则客户端可能会缓存备用
	Endpoints []string

	// HTTP传输设置，如果不定义，则默认使用DefaultTransport
	Transport etcd.CancelableTransport

	// CheckRedirect指定用于处理HTTP重定向的策略。
	// 如果CheckRedirect不为nil，则客户端将在执行HTTP重定向之前调用它。
	// 唯一的参数是已经发出的请求数。 如果CheckRedirect返回错误，则Client.Do将不再发出任何其他请求，并将错误返回给调用方。
	// 如果CheckRedirect为nil，则客户端使用其默认策略，该策略将在连续10个请求后停止
	CheckRedirect etcd.CheckRedirectFunc

	Username string

	Password string

	// 值为0表示不超时
	HeaderTimeoutPerRequest time.Duration

	// SelectionMode是一个EndpointSelectionMode枚举，它指定用于选择向其发送请求的etcd群集节点的策略。
	SelectionMode etcd.EndpointSelectionMode
}

// 客户端
type Client struct {
	config     *Config
	EtcdConfig etcd.Config
	EtcdClient etcd.Client
	EtcdAPI    etcd.KeysAPI
}

// key的信息
type Key struct {
	Name string `json:"key"`
	// 如果key是个目录，则返回true
	Dir           bool   `json:"dir,omitempty"`
	Value         string `json:"value"`
	CreatedIndex  uint64 `json:"createdIndex"`
	ModifiedIndex uint64 `json:"modifiedIndex"`
	// key的到期时间
	Expiration *time.Time `json:"expiration,omitempty"`
	// key的生存时间(秒)
	TTL int64 `json:"ttl,omitempty"`
}

// 获得一个指定配置的客户端实例
func New(config *Config) (*Client, error) {
	var err error
	var client Client
	client.config = config
	client.EtcdConfig.Endpoints = config.Endpoints
	client.EtcdConfig.HeaderTimeoutPerRequest = config.HeaderTimeoutPerRequest
	client.EtcdConfig.Username = config.Username
	client.EtcdConfig.Password = config.Password
	client.EtcdConfig.CheckRedirect = config.CheckRedirect
	client.EtcdConfig.Transport = config.Transport
	client.EtcdConfig.SelectionMode = config.SelectionMode
	client.EtcdClient, err = etcd.New(client.EtcdConfig)
	if err != nil {
		return nil, err
	}
	client.EtcdAPI = etcd.NewKeysAPI(client.EtcdClient)
	return &client, nil
}

func (receiver *Client) Create(key, value string) error {
	_, err := receiver.EtcdAPI.Create(context.Background(), key, value)
	return err
}

func (receiver *Client) CreateInOrder(dir, value string, ttl time.Duration) error {
	_, err := receiver.EtcdAPI.CreateInOrder(context.Background(), dir, value, &etcd.CreateInOrderOptions{TTL: ttl})
	return err
}

func (receiver *Client) Set(key, value string, ttl time.Duration) error {
	_, err := receiver.EtcdAPI.Set(context.Background(), key, value, &etcd.SetOptions{
		TTL:              ttl,
		Refresh:          false,
		NoValueOnSuccess: false,
	})
	return err
}

func (receiver *Client) Get(key string) (*Key, error) {
	resp, err := receiver.EtcdAPI.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	return &Key{
		Name:          resp.Node.Key,
		Dir:           resp.Node.Dir,
		Value:         resp.Node.Value,
		CreatedIndex:  resp.Node.CreatedIndex,
		ModifiedIndex: resp.Node.ModifiedIndex,
		Expiration:    resp.Node.Expiration,
		TTL:           resp.Node.TTL,
	}, nil
}

func (receiver *Client) Update(key, value string) error {
	_, err := receiver.EtcdAPI.Update(context.Background(), key, value)
	return err
}

func (receiver *Client) Delete(key, dir bool) error {
	_, err := receiver.EtcdAPI.Delete(context.Background(), key, &etcd.DeleteOptions{Dir: false})
	return err
}
