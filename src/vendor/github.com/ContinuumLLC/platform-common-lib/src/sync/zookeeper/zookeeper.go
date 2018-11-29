package zookeeper

import (
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/sync"
	"github.com/ContinuumLLC/platform-common-lib/src/web/rest"
	"github.com/samuel/go-zookeeper/zk"
)

var connection Connection

const ignoreME = "IGNORE-ME"

type zookeeper struct {
	config sync.Config
	log    logging.Logger
}

//Instance : Returns an instance of Zookeper implamentation for Sync Service
func Instance(config sync.Config) sync.Service {
	return zookeeper{
		config: config,
		log:    logging.GetLoggerFactory().Get(),
	}
}

func (z zookeeper) connect() error {
	if connection == nil {
		z.log.Logf(logging.DEBUG, "Creating Connection for Servers : %s ", z.config.Servers)
		conn, _, err := zk.Connect(z.config.Servers, (time.Duration(z.config.SessionTimeoutInSecond) * time.Second))

		if err != nil {
			return err
		}
		connection = conn
	}
	return nil
}

func (z zookeeper) Send(path string, data string) error {
	if err := z.connect(); err != nil {
		return err
	}
	return z.send(path, data, connection)
}

func (z zookeeper) send(path string, data string, conn Connection) error {
	s, err := z.createNode(path, data, conn)
	if err != nil {
		z.log.Logf(logging.ERROR, "Create Node Error : %v", err)
		return err
	}
	s, err = conn.Set(path, []byte(data), -1)
	if err != nil {
		return err
	}
	z.log.Logf(logging.TRACE, "Sending Data %s on Path : %s at Version : %d", data, path, s.Version)
	return nil
}

func (z zookeeper) createNode(path string, data string, conn Connection) (*zk.Stat, error) {
	flags := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)

	found, s, err := conn.Exists(path)

	if err != nil {
		z.log.Logf(logging.ERROR, "Listen Find Path Error : %v", err)
		return nil, err
	}

	if !found {
		_, err = conn.Create(path, []byte(data), flags, acl)
	}
	return s, err
}

func (z zookeeper) Listen(path string, c chan sync.Response) error {
	if err := z.connect(); err != nil {
		z.log.Logf(logging.ERROR, "Listen Connect Error : %v", err)
		return err
	}

	defer connection.Close()

	for {
		z.listen(path, connection, c)
	}
}

func (z zookeeper) listen(path string, conn Connection, c chan sync.Response) {
	data, _, ech, err := conn.GetW(path)
	if err == zk.ErrNoNode {
		_, err = z.createNode(path, ignoreME, conn)
		if err != nil {
			z.log.Logf(logging.ERROR, "Create Node Error : %v", err)
			c <- sync.Response{Error: err}
		}
		return
	}

	if err != nil {
		z.log.Logf(logging.ERROR, "Listen Get Error : %v", err)
		c <- sync.Response{Error: err}
		return
	}

	d := string(data)
	if d != ignoreME {
		z.log.Logf(logging.TRACE, "listen Data : %s", d)
		c <- sync.Response{Data: d}
	}
	<-ech
}

func (z zookeeper) connectionState() (string, error) {
	if err := z.connect(); err != nil {
		z.log.Logf(logging.ERROR, "Listen Connect Error : %v", err)
		return "", err
	}
	return connection.State().String(), nil
}

func (z zookeeper) Health() rest.Statuser {
	return zookeeperStatus{
		config:  z.config,
		service: z,
	}
}

type zookeeperStatus struct {
	config  sync.Config
	service zookeeper
}

func (s zookeeperStatus) Status(conn rest.OutboundConnectionStatus) *rest.OutboundConnectionStatus {
	conn.ConnectionType = "Zookeeper"
	conn.ConnectionURLs = s.config.Servers
	state, err := s.service.connectionState()
	if err != nil {
		conn.ConnectionStatus = rest.ConnectionStatusUnavailable
	} else {
		conn.ConnectionStatus = state
	}
	return &conn
}
