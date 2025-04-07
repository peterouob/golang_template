package im

import (
	"github.com/peterouob/golang_template/api/protobuf"
	"github.com/peterouob/golang_template/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

//broadcast.go REF:https://github.com/rodaine/grpc-chat/blob/main/server.go

type BroadCastServer struct {
	protobuf.UnimplementedChatServer
	Broadcast     chan *protobuf.BroadCastResponse
	ClientTokens  map[string]string
	ClientStreams map[string]chan *protobuf.BroadCastResponse
	clientCache   sync.Map
	sync.RWMutex
}

func NewBroadCastServer() *BroadCastServer {
	return &BroadCastServer{
		Broadcast:     make(chan *protobuf.BroadCastResponse, 1024),
		ClientTokens:  make(map[string]string),
		ClientStreams: make(map[string]chan *protobuf.BroadCastResponse),
	}
}

func (b *BroadCastServer) BroadCast(srv protobuf.Chat_BroadCastServer) error {
	utils.Log("Hello this is BroadCast")
	ctx := srv.Context()
	id := ctx.Value("id")
	log.Println(id)
	return nil
}

func (b *BroadCastServer) openStream(token string) (stream chan *protobuf.BroadCastResponse) {
	stream = make(chan *protobuf.BroadCastResponse, 1024)
	b.Lock()
	b.ClientStreams[token] = stream
	b.clientCache.Store(token, stream)
	b.Unlock()
	return
}

func (b *BroadCastServer) sendBroadcast(srv protobuf.Chat_BroadCastServer, token string) {
	stream := b.openStream(token)
	defer b.closeStream(token)

	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			if s, ok := status.FromError(srv.Send(res)); ok {
				switch s.Code() {
				//log.Println(s.Message())
				case codes.OK:
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					utils.ErrorMsgF("client (%s) terminated connection", token)
				default:
					utils.ErrorMsgF("failed to send to client (%s): %v", token, s.Err())
				}
			}
		}
	}
}

func (b *BroadCastServer) closeStream(token string) {
	b.Lock()
	if stream, ok := b.ClientStreams[token]; ok {
		delete(b.ClientStreams, token)
		close(stream)
	}
	//TODO:感覺可以把原本的login cache record紀錄redis,等確定一定時間後再刪除
	b.Unlock()
}
