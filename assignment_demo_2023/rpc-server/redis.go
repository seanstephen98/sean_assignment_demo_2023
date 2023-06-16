package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (client *RedisClient) InitClient(contxt context.Context, addr, password string) error {
	newClient := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})

	anError := newClient.Ping(contxt).Err()

	if anError != nil {
		return anError
	}

	client.cli = newClient

	return nil
}

func (client *RedisClient) SaveMessage(contxt context.Context, roomID string, message *Message) error {

	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	member := &redis.Z{
		Score:  (float64)(message.Timestamp),
		Member: text,
	}

	_, err = client.cli.ZAdd(contxt, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessagesByRoomID(ctx context.Context, roomID string, start, end int64, sortReverse bool) ([]*Message, error) {
	var (
		rawMessages []string
		messages    []*Message
		err         error
	)

	if sortReverse {
		// Desc order with time -> first message is the latest message
		rawMessages, err = c.cli.ZRevRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// Asc order with time -> first message is the earliest message
		rawMessages, err = c.cli.ZRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msg := range rawMessages {
		temp := &Message{}
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}
