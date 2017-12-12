package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type Row struct {
	Value string `json:"value"`
}

func Put(ctx context.Context, client *firestore.Client, id string, row Row) (*firestore.DocumentRef, error) {
	doc := client.Doc(fmt.Sprintf("Fire/%s", id))
	_, err := doc.Create(ctx, row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}
