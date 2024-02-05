package clients

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gorcon/rcon"
)

type RconClientOptions struct {
	Host     string
	Port     string
	Password string
}

type rconClient struct {
	connection *rcon.Conn
}

type RconClient interface {
	Close()
	Broadcast(message string)
}

func Rcon(options RconClientOptions) (RconClient, error) {
	if options.Host == "" {
		return nil, errors.New("missing rcon host")
	}

	if options.Port == "" {
		return nil, errors.New("missing rcon port")
	}

	if options.Password == "" {
		return nil, errors.New("missing rcon password")
	}

	conn, err := rcon.Dial(fmt.Sprintf("%s:%s", options.Host, options.Port), options.Password)

	if err != nil {
		return nil, err
	}

	return &rconClient{
		connection: conn,
	}, nil
}

func (m *rconClient) Close() {
	m.connection.Close()
}

func (m *rconClient) Broadcast(message string) {
	words := strings.Split(message, " ")

	// determined this through experimentation, palworld will
	// automatically trim strings
	maxLineSize := 54

	// it can't handle whitespace in broadcast message so fake it with a control char
	spaceChar := '\x1f'
	spaceWidth := 2 // for purposes of calculating width

	toBroadcast := [][]string{}

	running := 0
	cur := []string{}

	for i, v := range words {
		running += len(v)

		if i != 0 {
			running += spaceWidth
		}

		if running < maxLineSize {
			cur = append(cur, v)
		} else {
			toBroadcast = append(toBroadcast, cur)
			running = len(v)

			cur = []string{}
			cur = append(cur, v)
		}
	}

	toBroadcast = append(toBroadcast, cur)

	for _, line := range toBroadcast {
		joinedLine := strings.Join(line, string(spaceChar))
		execute(m.connection, fmt.Sprintf("Broadcast %s", joinedLine))
	}
}

func execute(connection *rcon.Conn, command string) (string, bool) {
	result, err := connection.Execute(command)
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	} else {
		return strings.TrimRight(result, string('\n')), true
	}

}
