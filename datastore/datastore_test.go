package datastore

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
)

type HogeV1 struct {
	Value int
}

type HogeV2 struct {
	Value float64
}

var _ datastore.PropertyLoadSaver = &HogeV2{}

func (h *HogeV2) Load(ps []datastore.Property) error {
	var nps []datastore.Property
	for _, v := range ps {
		if v.Name == "Value" {
			switch i := v.Value.(type) {
			case int:
				v.Value = float64(i)
			case int8:
				v.Value = float64(i)
			case int16:
				v.Value = float64(i)
			case int32:
				v.Value = float64(i)
			case int64:
				v.Value = float64(i)
			default:
				// noop
			}
		}
		nps = append(nps, v)
	}

	err := datastore.LoadStruct(h, nps)
	if err != nil {
		return err
	}

	return nil
}

func (h *HogeV2) Save() ([]datastore.Property, error) {
	return datastore.SaveStruct(h)
}

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
		t.Fatalf("expected Value is %d, got %d", e, g)
	}
}
