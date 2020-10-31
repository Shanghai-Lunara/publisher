package runner

import (
	"github.com/gorilla/websocket"
	"github.com/nevercase/publisher/pkg/types"
	"k8s.io/klog"
	"net/url"
)

type client struct {
	conn *websocket.Conn
}

func newClient(addr string) (*client, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: types.WebsocketHandlerRunner}
	klog.Info("url:", u)
	a, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	c := &client{
		conn: a,
	}
	return c, nil
}
