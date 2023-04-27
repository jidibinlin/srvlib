package natsq

import (
	"github.com/nats-io/nats.go"
)

var (
	Url  = nats.DefaultURL
	conn *nats.Conn
)

func SetConnUrl(url string) error {
	if url != Url {
		Url = url
		if conn == nil {
			return nil
		}
		var err error
		conn, err = nats.Connect(Url)
		if err != nil {
			return err
		}
	}
	return nil
}

func Conn() (*nats.Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = nats.Connect(Url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func MustConn() *nats.Conn {
	var err error
	conn, err = Conn()
	if err != nil {
		panic(err)
	}
	return conn
}

func FlushConn() error {
	if conn == nil {
		return nil
	}

	return conn.Drain()
}
