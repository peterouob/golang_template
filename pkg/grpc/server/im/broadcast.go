package im

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/tools"
	"sync"
)

type BroadCastServer struct {
	protobuf.UnimplementedChatServer
	Broadcast     chan *protobuf.BroadCastResponse
	ClientStreams map[string]chan *protobuf.BroadCastResponse

	sync.RWMutex
}

func NewBroadCastServer() *BroadCastServer {
	return &BroadCastServer{
		Broadcast:     make(chan *protobuf.BroadCastResponse, 1024),
		ClientStreams: make(map[string]chan *protobuf.BroadCastResponse),
	}
}

func (b *BroadCastServer) BroadCast(srv protobuf.Chat_BroadCastServer) error {
	tools.Log("Hello this is BroadCast")
	return nil
}
