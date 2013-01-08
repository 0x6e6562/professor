package professor

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Connection struct {
	host           string
	conn           net.Conn
	stream         uint8
	shouldClose    chan bool
	pingInterval   time.Duration
	recvTimeout    time.Duration
	connectTimeout time.Duration
}

func Connect(host string) (*Connection, error) {

	connection := Connection{
		host:           host + ":9042",
		conn:           nil,
		stream:         uint8(0),
		shouldClose:    make(chan bool),
		recvTimeout:    time.Second,
		pingInterval:   10 * time.Second,
		connectTimeout: 1 * time.Second,
	}

	conn, err := net.DialTimeout("tcp", connection.host, connection.connectTimeout)

	if err != nil {
		return nil, err
	} else {

		connection.conn = conn
		log.Printf("Connected to %s", connection.host)
	}

	frame := Options(connection.nextStream())
	send(conn, frame)

	frame, err = recv(conn)
	if err != nil {
		return nil, err
	}

	if supported, ok := frame.body.(map[string][]string); ok {

		log.Printf("Remote peer supports following options %+v", supported)

		options := make(map[string]string)

		options["CQL_VERSION"] = "3.0.0" // TODO hardcoded - need to extract this value from the supported options

		frame = Startup(connection.nextStream(), options)
		send(conn, frame)

		frame, err = recv(conn)
		if err != nil {
			return nil, err
		}

		if frame.header.opcode != READY {
			return nil, fmt.Errorf("Expected %s opcode but received %s", READY, frame.header.opcode)
		}

	} else {
		return nil, errors.New("Could not get any supported options")
	}

	return &connection, nil
}

func (c *Connection) nextStream() uint8 {
	return c.stream
}

func (c *Connection) Query(cql LongString) (string, error) {
	query := &Query{cql: cql, consistency: ANY}
	sendFrame := Cql(c.nextStream(), query)
	send(c.conn, sendFrame)
	responseFrame, err := recv(c.conn)

	if err != nil {
		return "", err
	}

	switch response := responseFrame.body.(type) {
	case *Result:
		switch response.kind {
		case SET_KEYSPACE:
			return response.body.(string), nil
		}
	}

	return "", fmt.Errorf("Response processing not yet implemented")
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
