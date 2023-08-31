package publisher

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
	"strings"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-publisher"
	channel   = "order-notification"
)

type Publisher struct {
	clusterId,
	clientId,
	channel string
}

func New() *Publisher {
	return &Publisher{
		clusterId: clusterID,
		clientId:  clientID,
		channel:   channel,
	}
}

func (p *Publisher) publicMessage() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Write absolut path to json file")
	path, err := reader.ReadString('\n')
	if err != nil {
		return errors.New("not correct path to file")
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return errors.New("can not read the file")
	}
	con, err := stan.Connect(p.clusterId, p.clientId)
	defer con.Close()
	if err != nil {
		return errors.New("can not connect to nats streaming")
	}
	err = con.Publish(p.channel, file)
	if err != nil {
		return errors.New("can not publish data")
	}
	return nil
}

func (p *Publisher) CirclePublicMessage() {
	for {
		err := p.publicMessage()
		if err != nil {
			if err.Error() == "" {
				os.Exit(1)
			}
			fmt.Println(err)
		} else {
			fmt.Println("file was loaded")
		}
	}
}
