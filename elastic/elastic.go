package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/guneyin/app-sdk/logger"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Elastic struct {
	client *elasticsearch.Client
}

func New(addr string) (*Elastic, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{addr}})
	if err != nil {
		return nil, err
	}

	return &Elastic{client: client}, nil
}

func (e *Elastic) StoreDocument(ctx context.Context, index, id string, doc any) error {
	if e == nil {
		logger.Warn("elasticsearch module not initialized")
		return nil
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
		Pretty:     true,
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	return nil
}
