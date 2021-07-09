package domain

import (
	"context"
)

type Document struct {
	Word     string   `json:"word"`
	Examples []string `json:"examples"`
}
type DocumentResponse struct {
	Word     string   `json:"word"`
	Examples []string `json:"examples"`
	Score    float64  `json:"score"`
}

type DocumentUsecase interface {
	SearchWord(ctx context.Context, word string) ([]DocumentResponse, error)
	SearchSentence(ctx context.Context, sentence string) ([]DocumentResponse, error)
	Create(ctx context.Context, document *Document) error
}

type DocumentRepository interface {
	SearchWord(ctx context.Context, word string) ([]DocumentResponse, error)
	SearchSentence(ctx context.Context, sentence string) ([]DocumentResponse, error)
	Create(ctx context.Context, document *Document) error
}
