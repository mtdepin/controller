package mtlibp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p/core/network"
)

func CreateServer(ctx context.Context, ip string, port int, handleStream network.StreamHandler) (string, error) {
	port, err := FindFreePort("", 5)
	if err != nil {
		return "", err
	}

	host, err1 := MakeHost(port, rand.Reader)
	if err1 != nil {
		return "", err1
	}
	StartPeer(ctx, host, handleStream)

	return fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", port, host.ID().Pretty()), nil
}
