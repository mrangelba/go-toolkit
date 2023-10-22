package rabbitmq

import (
	"fmt"

	"sync"

	"github.com/mrangelba/go-toolkit/config"
	"github.com/mrangelba/go-toolkit/logger"
	gorabbitmq "github.com/wagslane/go-rabbitmq"
)

var once sync.Once
var instance *gorabbitmq.Conn

func GetRabbitMQConn() *gorabbitmq.Conn {
	once.Do(func() {
		instance = newConn()
	})

	return instance
}

func newConn() *gorabbitmq.Conn {
	cfg := config.Get()

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	oplog := logger.Get()

	conn, err := gorabbitmq.NewConn(
		url,
		gorabbitmq.WithConnectionOptionsLogger(oplog),
	)

	if err != nil {
		oplog.Panic().Msgf("failed to connect to RabbitMQ: %s", err)
	}

	return conn
}
