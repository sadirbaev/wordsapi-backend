package usecase

import (
	"context"
	"time"
	"wordsapi/domain"
)

type documentUsecase struct {
	documentRepo   domain.DocumentRepository
	contextTimeout time.Duration
}



func NewDocumentUsecase(repo domain.DocumentRepository, timout time.Duration) domain.DocumentUsecase {
	return &documentUsecase{
		documentRepo: repo,
		contextTimeout: timout,
	}
}

func (rx *documentUsecase) SearchWord(ctx context.Context, word string) ([]domain.DocumentResponse, error) {
	res, err := rx.documentRepo.SearchWord(ctx, word)
	return res, err
}

func (rx *documentUsecase) SearchSentence(ctx context.Context, sentence string) ([]domain.DocumentResponse, error) {
	res, err := rx.documentRepo.SearchSentence(ctx, sentence)
	return res, err
}