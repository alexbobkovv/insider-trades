package types

import (
	"encoding/base64"
	"fmt"
	"time"
)

type cursor struct {
	Encoded          string
	DecodedTimestamp *time.Time
}

func NewFromTime(timestampCursor *time.Time) *cursor {
	encodedTimestamp := encodeToString(timestampCursor)

	return &cursor{DecodedTimestamp: timestampCursor, Encoded: encodedTimestamp}
}

func NewFromEncodedString(encodedCursor string) (*cursor, error) {
	decodedCursor, err := decodeToTimestamp(encodedCursor)
	if err != nil {
		return nil, err
	}

	return &cursor{Encoded: encodedCursor, DecodedTimestamp: decodedCursor}, nil
}

func (c *cursor) GetEncoded() string {
	return c.Encoded
}

func (c *cursor) GetDecoded() *time.Time {
	return c.DecodedTimestamp
}

func encodeToString(decodedCursor *time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(decodedCursor.Format(time.RFC3339Nano)))
}

func decodeToTimestamp(encodedCursor string) (*time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, fmt.Errorf("decodeTimestampCursor: failed to decode cursor: %w", err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return nil, fmt.Errorf("decodeTimestampCursor: failed parse timestamp: %w", err)
	}

	return &timestamp, nil
}
