package main

import (
	"context"
	"fmt"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"strings"
	"time"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {

	if err := validateSendRequest(req); err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()

	message := &Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}

	roomID, err := getRoomID(req.Message.GetChat())
	if err != nil {
		return nil, err
	}

	err = rdb.SaveMessage(ctx, roomID, message)
	if err != nil {
		return nil, err
	}

	// we only reach here if its a successful response
	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "Success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	first := req.GetCursor()
	last := first + int64(req.GetLimit())

	messages, err := rdb.GetMessagesByRoomID(ctx, roomID, first, last, req.GetReverse())
	if err != nil {
		return nil, err
	}

	responses := make([]*rpc.Message, 0)

	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false

	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			hasMore = true
			nextCursor = last
			break
		}
		tempMessage := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}

		responses = append(responses, tempMessage)
		counter++
	}

	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = 0, "Success"
	resp.Messages = responses
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor
	return resp, nil
}

func getRoomID(chat string) (string, error) {
	var roomID string
	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, ":")

	if len(senders) != 2 {
		err := fmt.Errorf("Chat ID Invalid format, should be user1:user2, is '%s'", chat)
		return "", err
	}

	sender1, sender2 := senders[0], senders[1]

	comparator := strings.Compare(sender1, sender2)

	if comparator == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	}

	return roomID, nil
}

func validateSendRequest(req *rpc.SendRequest) error {

	senders := strings.Split(req.Message.Chat, ":")

	if len(senders) != 2 {
		err := fmt.Errorf("Chat ID Invalid format, should be user1:user2, is '%s'", req)
		return err
	}

	sender1, sender2 := senders[0], senders[1]

	reqSender := req.Message.GetSender()

	if reqSender != sender1 && reqSender != sender2 {
		err := fmt.Errorf("Chat Request Invalid value, is '%s' , sender not found", reqSender)
		return err
	}

	return nil
}
