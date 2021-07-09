package esservice

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"strings"
	"wordsapi/domain"
)

type examples struct {
	Fuzziness string `json:"fuzziness"`
	Value     string `json:"value"`
}

type fuzzy struct {
	Examples examples `json:"examples"`
}

type match struct {
	Fuzzy fuzzy `json:"fuzzy"`
}

type spanMulti struct {
	Match match `json:"match"`
}

type clause struct {
	SpanMulti spanMulti `json:"span_multi"`
}

type esDocumentRepository struct {
	ES *elastic.Client
}

func NewESDocumentRepository(es *elastic.Client) domain.DocumentRepository {
	return &esDocumentRepository{
		ES: es,
	}
}

func (rx *esDocumentRepository) Create(ctx context.Context, document *domain.Document) error {
	var index = viper.GetString(`elasticsearch.index`)
	body, err := json.Marshal(document)
	if err != nil {
		return err
	}
	contentString := string(body)
	doc := strings.Replace(document.Word, " ", "_", -1)
	_, err = rx.ES.Index().
		Index(index).
		OpType("create").
		BodyString(contentString).
		Id(doc).
		Do(ctx)
	return err
}

func (rx *esDocumentRepository) SearchWord(ctx context.Context, word string) ([]domain.DocumentResponse, error) {
	var index = viper.GetString(`elasticsearch.index`)
	esQuery := elastic.NewMultiMatchQuery(word, "word").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := rx.ES.Search().
		Index(index).
		Query(esQuery).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	var res []domain.DocumentResponse
	for _, data := range result.Hits.Hits {
		doc := &domain.DocumentResponse{}
		if err = json.Unmarshal(data.Source, doc); err == nil {
			doc.Score = *data.Score
			res = append(res, *doc)
		}
	}
	return res, err
}

func (rx *esDocumentRepository) SearchSentence(ctx context.Context, sentence string) ([]domain.DocumentResponse, error) {
	var index = viper.GetString(`elasticsearch.index`)
	var clauses []elastic.Query

	words := strings.Split(sentence, " ")
	for _, word := range words {
		clause := &clause{
			SpanMulti: spanMulti{
				Match: match{
					Fuzzy: fuzzy{
						Examples: examples{
							Fuzziness: "1",
							Value:     word,
						},
					},
				},
			},
		}
		marshal, err := json.Marshal(clause)

		if err != nil {
			return nil, err
		}
		clauses = append(clauses, elastic.NewRawStringQuery(string(marshal)))
	}

	esQuery := elastic.NewSpanNearQuery(clauses...).Slop(len(words)).InOrder(false)
	result, err := rx.ES.Search().
		Index(index).
		Query(esQuery).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	res := []domain.DocumentResponse{}
	for _, data := range result.Hits.Hits {
		doc := &domain.DocumentResponse{}
		err = json.Unmarshal(data.Source, doc)
		if err == nil {
			doc.Score = *data.Score
			res = append(res, *doc)
		}
	}
	return res, nil
}
