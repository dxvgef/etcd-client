package v2

import (
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	// 创建客讔端实例
	client, err := New(&Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// 设置目录，如果目录不存在则自动创建
	if err = client.SetDir("/dxvgef", 30*time.Second); err != nil {
		t.Fatal(err)
	}

	// 创建k/v，如果key存在则报错
	if err = client.Create("/dxvgef/test", "Hello ETCD by Create"); err != nil {
		t.Error(err)
	}

	// 设置k/v，如果key存在则更新
	if err = client.Set("/dxvgef/test", "Hello ETCD by Set", 30*time.Second); err != nil {
		t.Error(err)
	}

	// 读取key的信息，如果key不存在则报错
	key, err := client.Get("/dxvgef/test")
	if err != nil {
		t.Error(err)
	}
	t.Log(key.Value, key.TTL, key.Expiration.Local())

	// 设置key的信息，如果key不存在则创建
	if err = client.Set("/dxvgef/test2", "Hello ETCD by Set", 30*time.Second); err != nil {
		t.Error(err)
	}

	// 更新k/v，如果key不存在则报错
	if err = client.Update("/dxvgef/test2", "Hello ETCD by Update"); err != nil {
		t.Error(err)
	}

	// 读取key的信息，如果key不存在则报错
	key, err = client.Get("/dxvgef/test2")
	if err != nil {
		t.Error(err)
	}
	t.Log(key.Value, key.TTL, key.Expiration.Local())

	// 删除一个key
	if err = client.Delete("/dxvgef/test2"); err != nil {
		t.Error(err)
	}

	// 删除一个目录，包含目录下的所有key(force: true)
	if err = client.DeleteDir("/dxvgef", true); err != nil {
		t.Error(err)
	}
}
