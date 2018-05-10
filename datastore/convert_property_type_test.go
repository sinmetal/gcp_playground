package datastore

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
)

// HogeV1 is Old Entity struct
type HogeV1 struct {
	Value int
}

// HogeV2 is New Entity struct
type HogeV2 struct {
	Value float64
}

var _ datastore.PropertyLoadSaver = &HogeV2{}

// Load is datastore.PropertyLoadSaver を満たすためのfunc
// Entityをstructに変換する処理
func (h *HogeV2) Load(ps []datastore.Property) error {
	for idx, v := range ps {
		if v.Name == "Value" {
			switch i := v.Value.(type) {
			case int64:
				v.Value = float64(i)
			default:
				// noop
			}
			ps[idx] = v
		}
	}

	return datastore.LoadStruct(h, ps)
}

// Save is datastore.PropertyLoadSaver を満たすためのfunc
// structをEntityに変換する処理
func (h *HogeV2) Save() ([]datastore.Property, error) {
	return datastore.SaveStruct(h)
}

// TestConvertPropertyType is Sample Test
func TestConvertPropertyType(t *testing.T) {
	ctx := context.Background()

	ds, err := datastore.NewClient(ctx, "testproject")
	if err != nil {
		t.Fatal(err)
	}

	key := datastore.NameKey("Hoge", "key1", nil)
	_, err = ds.Put(ctx, key, &HogeV1{
		Value: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	var h HogeV2
	if err := ds.Get(ctx, key, &h); err != nil {
		t.Fatal(err)
	}
	if e, g := 10.0, h.Value; e != g {
		t.Fatalf("expected Value is %f, got %f", e, g)
	}
}
