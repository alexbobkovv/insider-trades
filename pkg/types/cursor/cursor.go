package cursor

import (
	"encoding/base64"
	"fmt"
	"time"
)

type Cursor struct {
	Encoded          string
	DecodedTimestamp *time.Time
	Empty            bool
}

func NewFromTime(timestampCursor *time.Time) *Cursor {
	encodedTimestamp := encodeToString(timestampCursor)

	return &Cursor{DecodedTimestamp: timestampCursor, Encoded: encodedTimestamp}
}

func NewFromEncodedString(encodedCursor string) (*Cursor, error) {
	if encodedCursor == "" {
		return NewEmpty(), nil
	}

	decodedCursor, err := decodeToTimestamp(encodedCursor)
	if err != nil {
		return nil, err
	}

	return &Cursor{Encoded: encodedCursor, DecodedTimestamp: decodedCursor}, nil
}

func NewEmpty() *Cursor {
	return &Cursor{
		Encoded:          "",
		DecodedTimestamp: nil,
		Empty:            true,
	}
}

func (c *Cursor) GetEncoded() string {
	return c.Encoded
}

func (c *Cursor) GetDecoded() *time.Time {
	return c.DecodedTimestamp
}

func (c *Cursor) GetUNIXTime() int64 {
	return c.DecodedTimestamp.Unix()
}

func (c *Cursor) IsEmpty() bool {
	return c.Empty
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
