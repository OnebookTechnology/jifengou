package etcdsdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/robfig/config"
	"strconv"
	"strings"
	"time"
)

type EtcdService struct {
	cluster []string
	client  *clientv3.Client
}

var (
	DialTimeout    = 10 * time.Second
	RequestTimeout = 10 * time.Second
)

func NewEtcdService(cfgPath string) (*EtcdService, error) {
	c, err := config.ReadDefault(cfgPath)
	if err != nil {
		return nil, err
	}
	etcdService := new(EtcdService)

	cluster, err := c.String("OneBookConsistence", "cluster")
	if err != nil {
		return nil, err
	}
	etcdService.cluster = strings.Split(cluster, ",")

	client, err := clientv3.New(clientv3.Config{Endpoints: etcdService.cluster, DialTimeout: DialTimeout})
	if err != nil {
		return nil, err
	}
	etcdService.client = client
	return etcdService, nil
}

func (es *EtcdService) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	resp, err := es.client.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}
	for _, ev := range resp.Kvs {
		fmt.Println(string(ev.Key), string(ev.Value))
	}
	for _, ev := range resp.Kvs {
		return string(ev.Value), nil
	}
	return "", nil
}

func (es *EtcdService) GetWithPrefix(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	resp, err := es.client.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	var values map[string]string
	values = make(map[string]string)
	for _, ev := range resp.Kvs {
		fmt.Println(string(ev.Key), string(ev.Value))
		values[string(ev.Key)] = string(ev.Value)
	}
	return values, nil
}

func (es *EtcdService) Put(key, value string, ttl time.Duration) error {
	seconds := int64(ttl.Seconds())
	if ttl < 0 {
		return errors.New("invalid ttl value: " + strconv.FormatInt(seconds, 10))
	}

	//Duration Storage
	if ttl == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
		_, err := es.client.Put(ctx, key, value)
		cancel()
		if err != nil {
			return err
		}
		return nil
	} else {
		//With ttl
		// minimum lease ttl is 5-second
		resp, err := es.client.Grant(context.TODO(), seconds)
		if err != nil {
			return err
		}
		_, err = es.client.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
		if err != nil {
			return err
		}
		return nil
	}
}

func (es *EtcdService) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()
	_, err := es.client.Delete(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	return nil
}
