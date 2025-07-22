package database

import (
	"context"
	"time"
)

type Message struct {
	SenderID   int       `db:"sender_id" json:"sender_id"`
	Sender     string    `db:"sender" json:"sender"`
	ReceiverID int       `db:"receiver_id" json:"receiver_id"`
	Message    string    `db:"message" json:"message"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (db *DB) InsertMessage(senderid, receiverid int, message string, currentDateTime string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
    INSERT INTO message (sender_id, receiver_id, message, created_at)
    VALUES ($1, $2, $3, $4)`

	result, err := db.ExecContext(ctx, query, senderid, receiverid, message, currentDateTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

// GetPaginatedMessageHistory returns message history with pagination
func (db *DB) GetPaginatedMessageHistory(senderid, receiverid, offset, limit int) ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var messages []Message

	query := `
    SELECT m.sender_id, m.receiver_id, m.message, m.created_at,
           u.username as sender
    FROM message m
    JOIN user u ON m.sender_id = u.id
    WHERE (m.sender_id = $1 AND m.receiver_id = $2) 
       OR (m.sender_id = $2 AND m.receiver_id = $1)
    ORDER BY m.created_at DESC
    LIMIT $3 OFFSET $4`

	err := db.SelectContext(ctx, &messages, query, senderid, receiverid, limit, offset)
	if err != nil {
		return nil, err
	}

	// Reverse the order to show oldest first (for chat display)
	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetMessageCount returns total number of messages between two users
func (db *DB) GetMessageCount(senderid, receiverid int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var count int
	query := `
    SELECT COUNT(*) 
    FROM message 
    WHERE (sender_id = $1 AND receiver_id = $2) 
       OR (sender_id = $2 AND receiver_id = $1)`

	err := db.GetContext(ctx, &count, query, senderid, receiverid)
	return count, err
}
