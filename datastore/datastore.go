package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Row struct {
	Value string `json:"value"`
}

func Put(ctx context.Context, client *datastore.Client, row Row) (*datastore.Key, error) {
	key, err := client.Put(ctx, datastore.IncompleteKey("Item", nil), &row)
	if err != nil {
		return nil, err
	}

	return key, nil
}
