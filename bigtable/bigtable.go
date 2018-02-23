package bigtable

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"cloud.google.com/go/bigtable"
)

type BigtableRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func CreateTableWithColumnFamily(ctx context.Context, adminClient *bigtable.AdminClient, table string, columnFamily string) error {
	tables, err := adminClient.Tables(ctx)
	if err != nil {
		return fmt.Errorf("Could not fetch table list: instance=%s : %v", instance, err)
	}

	if !sliceContains(tables, table) {
		if err := adminClient.CreateTable(ctx, table); err != nil {
			return fmt.Errorf("Could not create table %s: %v", table, err)
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, table)
	if err != nil {
		return fmt.Errorf("Could not read info for table %s: %v", table, err)
	}

	if !sliceContains(tblInfo.Families, columnFamily) {
		if err := adminClient.CreateColumnFamily(ctx, table, columnFamily); err != nil {
			return fmt.Errorf("Could not create column family %s: %v", columnFamily, err)
		}
	}

	return nil
}

func GetRange(ctx context.Context, projectID string, instance string, table string, family string, column string) ([]BigtableRow, error) {
	client, err := bigtable.NewClient(ctx, projectID, instance)
	if err != nil {
		return nil, fmt.Errorf("failed Bigtable.NewClient(): projectID=%s, instance=%s", projectID, instance)
	}

	tbl := client.Open(table)

	var rows []BigtableRow
	err = tbl.ReadRows(ctx, bigtable.PrefixRange(column), func(row bigtable.Row) bool {
		item := row[family][0]
		rows = append(rows, BigtableRow{
			Key:   row.Key(),
			Value: string(item.Value),
		})
		return true
	}, bigtable.RowFilter(bigtable.ColumnFilter(column)))

	if err = client.Close(); err != nil {
		return nil, fmt.Errorf("Could not close data operations client: %v", err)
	}

	return rows, nil
}

func UpdateBigtable(ctx context.Context, client *bigtable.Client, table string, family string, column string) error {
	tbl := client.Open(table)
	mut := bigtable.NewMutation()
	mut.Set(family, column, bigtable.Now(), []byte("Hello Bigtable"))
	rowKey := fmt.Sprintf("%s%d", column, time.Now().UnixNano())

	err := tbl.Apply(ctx, rowKey, mut)
	if err != nil {
		return fmt.Errorf("Could not apply bulk row mutation: %v", err)
	}

	return nil
}

func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
