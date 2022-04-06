package elastic

import (
	"context"
	"log"

	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
	"github.com/elastic/go-elasticsearch/v8"
)

type EsRepoImpl struct {
	esClient *elasticsearch.Client
}

func NewElasticRepository(esClient *elasticsearch.Client) domain.ESRepository {
	return &EsRepoImpl{
		esClient: esClient,
	}
}

func (e *EsRepoImpl) Create(ctx context.Context, user web.UserCreateEs) error {
	log.Println("create")
	log.Println(user)
	return nil
}

func (e *EsRepoImpl) Update(ctx context.Context, user web.UserCreateEs) error {
	log.Println("update")
	log.Println(user)
	return nil
}

func (e *EsRepoImpl) Search(ctx context.Context, name string, jurusan string, angkatan int) (*domain.SearchResponse, error) {
	return nil, nil
}
