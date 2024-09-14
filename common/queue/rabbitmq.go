package queue

import (
	"fmt"
	"ginProject/config"
	"github.com/streadway/amqp"
	"sync"
)

var AMQP = new(AMQPConnectionPool)

type AMQPConnectionPool struct {
	mu    sync.Mutex
	conns chan *amqp.Connection
	options
}
type options struct {
	maxOpen     int
	maxIdle     int
	maxAttempts int
	url         string
}

func Init() {
	AMQP = NewAMQPConnectionPool(&options{
		maxOpen:     config.RabbitMQMaxOpen,
		maxIdle:     config.RabbitMQMaxIdle,
		maxAttempts: config.RabbitMQMaxAttemps,
		url:         config.RabbitMQ,
	})
}

func NewAMQPConnectionPool(o *options) *AMQPConnectionPool {
	return &AMQPConnectionPool{
		options: *o,
		conns:   make(chan *amqp.Connection, o.maxOpen),
	}
}

func (p *AMQPConnectionPool) Get() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	attempts := 0

	for attempts < p.maxAttempts {
		conn, err = p.GetOne()
		if err == nil && conn != nil {
			return conn, nil
		}
		attempts++
	}

	return nil, fmt.Errorf("failed to get a valid connection after %d attempts", p.maxAttempts)
}

func (p *AMQPConnectionPool) GetOne() (*amqp.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		if len(p.conns) < p.maxOpen {
			conn, err := amqp.Dial(p.url)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		return nil, fmt.Errorf("no available connections")
	}
}

func (p *AMQPConnectionPool) Put(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns) >= p.maxIdle {
		conn.Close()
		return
	}
	p.conns <- conn
}

func (p *AMQPConnectionPool) DeclareQueue(name string) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(name, true, false, false, false, nil)
	return err
}

func (p *AMQPConnectionPool) Publish(exchange, key string, msg []byte) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	return err
}

func (p *AMQPConnectionPool) Consume(queueName string, handler func(delivery amqp.Delivery)) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		handler(msg)
	}

	return nil
}

func (p *AMQPConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.conns) > 0 {
		conn := <-p.conns
		err := conn.Close()
		if err != nil {
			//logger.Logger.Errorf("Error closing connection pool: %v", err)
		}
		fmt.Println("close", conn)
	}
}
