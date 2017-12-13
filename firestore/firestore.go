package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type Row struct {
	Value     string    `json:"value"`
	Number1   int       `json:"number1"`
	Number2   int       `json:"number2"`
	CreatedAt time.Time `json:"createdAt"`
}

func Put(ctx context.Context, client *firestore.Client, id string, row Row) (*firestore.DocumentRef, error) {
	doc := client.Doc(fmt.Sprintf("Fire/%s", id))
	_, err := doc.Create(ctx, row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}

func PutTxMulti(ctx context.Context, client *firestore.Client, id string, row Row) error {
	return client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for i := 0; i < 30; i++ {
			doc := client.Doc(fmt.Sprintf("Fire/%s-%d", id, i))
			err := tx.Create(doc, row)
			if err != nil {
				return errors.WithStack(err)
			}
			//if i > 25 {
			//	return errors.New("EGが２５超えたぞ")
			//}
		}
		return nil
	})
}

func List(ctx context.Context, client *firestore.Client) ([]*firestore.DocumentSnapshot, error) {
	fires := client.Collection("Fire")
	q := fires.OrderBy("CreatedAt", firestore.Desc).Limit(100)
	return q.Documents(ctx).GetAll()
}
