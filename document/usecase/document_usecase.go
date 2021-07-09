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
		documentRepo:   repo,
		contextTimeout: timout,
	}
}

func (rx *documentUsecase) SearchWord(c context.Context, word string) ([]domain.DocumentResponse, error) {
	ctx, cancel := context.WithTimeout(c, rx.contextTimeout)
	defer cancel()

	return rx.documentRepo.SearchWord(ctx, word)
}

func (rx *documentUsecase) SearchSentence(c context.Context, sentence string) ([]domain.DocumentResponse, error) {
	ctx, cancel := context.WithTimeout(c, rx.contextTimeout)
	defer cancel()
	return rx.documentRepo.SearchSentence(ctx, sentence)
}

func (rx *documentUsecase) Create(c context.Context, document *domain.Document) error {
	ctx, cancel := context.WithTimeout(c, rx.contextTimeout)
	defer cancel()
	return rx.documentRepo.Create(ctx, document)
}
