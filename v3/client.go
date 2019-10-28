package v3

import (
	"context"
	"crypto/tls"
	"time"

	etcd "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// 客户端配置
type Config struct {
	// 定义一组端点，如果其中一个不可用，则客户端将尝试使用其它的端点进行操作
	// 如果曾经调用Client.Sync，则客户端可能会缓存备用
	Endpoints []string

	Username string

	Password string

	// 值为0表示不超时
	HeaderTimeoutPerRequest time.Duration

	// 使用其最新成员更新端点的时间间隔
	// 0禁用自动同步。 默认情况下，禁用自动同步
	AutoSyncInterval time.Duration

	// 与端点的连接超时时间
	DialTimeout time.Duration

	// 客户端保持连接的ping间隔时间
	DialKeepAliveTime time.Duration

	// 客户端等待持续活动探测响应的时间，如果此时未收到响应，则连接将关闭
	DialKeepAliveTimeout time.Duration

	// 客户端请求限制（以字节为单位）。
	// 如果为0，则默认为2.0 MiB（2 * 1024 * 1024）。
	// 确保“ MaxCallSendMsgSize” <服务器端默认的发送/接收限制
	// 设置etcd的"embed.Config.MaxRequestBytes"或"--max-request-bytes"参数可以达到同样效果
	MaxCallSendMsgSize int

	// 客户端响应的限制
	// 如果为0，则默认为“ math.MaxInt32”，因为范围响应可以轻松超过请求发送限制
	// 确保“ MaxCallRecvMsgSize”> =服务器端默认的发送/接收限制
	// 设置etcd的"embed.Config.MaxRequestBytes"或"--max-request-bytes"参数可以达到同样效果
	MaxCallRecvMsgSize int

	TLS *tls.Config

	// 拒绝针对过时的群集创建客户端
	RejectOldCluster bool

	// grpc连接参数
	// 例如，传递"grpc.WithBlock()"以阻塞直到基础连接建立
	// 如果不这样做，Dial将立即返回，并且连接服务器将在后台进行
	DialOptions []grpc.DialOption

	// 日志记录器配置，如果值为nil，则使用zap的默认配置
	LogConfig *zap.Config

	// 默认的客户端context，它可以用来取消grpc连接和其他没有明确上下文的操作
	Context context.Context

	// 允许客户端在没有任何活动流（RPC）的情况下向服务器发送保持活动的ping指令
	PermitWithoutStream bool
}

// 客户端
type Client struct {
	config     *Config
	EtcdConfig etcd.Config
	EtcdClient *etcd.Client
}

// 获得一个指定配置的客户端实例
func New(config *Config) (*Client, error) {
	var err error
	var client Client
	client.config = config
	client.EtcdConfig.Endpoints = config.Endpoints
	client.EtcdConfig.Username = config.Username
	client.EtcdConfig.Password = config.Password
	client.EtcdConfig.AutoSyncInterval = config.AutoSyncInterval
	client.EtcdConfig.DialTimeout = config.DialTimeout
	client.EtcdConfig.DialKeepAliveTime = config.DialKeepAliveTime
	client.EtcdConfig.DialKeepAliveTimeout = config.DialKeepAliveTimeout
	client.EtcdConfig.MaxCallSendMsgSize = config.MaxCallSendMsgSize
	client.EtcdConfig.MaxCallRecvMsgSize = config.MaxCallRecvMsgSize
	client.EtcdConfig.TLS = config.TLS
	client.EtcdConfig.RejectOldCluster = config.RejectOldCluster
	client.EtcdConfig.DialOptions = config.DialOptions
	client.EtcdConfig.LogConfig = config.LogConfig
	client.EtcdConfig.Context = config.Context
	client.EtcdConfig.PermitWithoutStream = config.PermitWithoutStream

	client.EtcdClient, err = etcd.New(client.EtcdConfig)
	if err != nil {
		return nil, err
	}
	return &client, nil
}
