package kafka

type Kafka struct {
}

func New() (*Kafka, error) {
	return &Kafka{}, nil
}
