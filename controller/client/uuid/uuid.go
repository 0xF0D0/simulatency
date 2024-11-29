package uuid

import "github.com/google/uuid"

type UUIDGenerator interface {
	GenerateString() string
}

type uuidGenerator struct{}

func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}

func (u *uuidGenerator) GenerateString() string {
	return uuid.New().String()
}
