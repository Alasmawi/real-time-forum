package database

import (
"context"
"time"
)

type Message struct {
SenderID   int       `json:"sender_id"`
Sender     string    `json:"sender"`
ReceiverID int       `json:"receiver_id"`
Message    string    `json:"message"`
CreatedAt  time.Time `json:"created_at"`
}

func (db *DB) InsertMessage(senderid, recieverid int, message string) (int, error) {
ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
defer cancel()

query := `
INSERT INTO message (sender_id, reciever_id, message, created_at)
VALUES ($1, $2, $3, $4)`

result, err := db.ExecContext(ctx, query, senderid, recieverid, message, time.Now())
if err != nil {
return 0, err
}

id, err := result.LastInsertId()
if err != nil {
return 0, err
}
return int(id), err
}

// func (db *DB) GetMessagesForUser(userid int) ([]Message, error) {
// ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
// defer cancel()

// var messages []Message

// query := `SELECT m.id, m.sender_id, m.reciever_id, m.message, m.created_at
//   FROM message m
//   WHERE m.reciever_id = $1`

// err := db.GetContext(ctx, &messages, query, userid)
// if err != nil {
// return nil, err
// }

// return messages, nil
// }

func (db *DB) GetMessagesBetweenUsers(senderid, recieverid int) ([]Message, error) {
ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
defer cancel()

var messages []Message

query := `SELECT m.id, m.sender_id, m.reciever_id, m.message, m.created_at
  FROM message m
  WHERE m.sender_id = $1 AND m.reciever_id = $2`

err := db.GetContext(ctx, &messages, query, senderid, recieverid)
if err != nil {
return nil, err
}

return messages, nil
}
