package spanner

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func readStaleData(ctx context.Context, client *spanner.Client, sec int) error {
	ro := client.ReadOnlyTransaction().WithTimestampBound(spanner.ExactStaleness(time.Duration(sec) * time.Second))
	defer ro.Close()

	iter := ro.Read(ctx, "Albums", spanner.AllKeys(), []string{"SingerId", "AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID string
		var albumID string
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", singerID, albumID, albumTitle)
	}
}

func readStaleData2(ctx context.Context, client *spanner.Client, sec int) error {
	now := time.Now()
	t := now.Add(time.Duration(-sec) * time.Second)
	ro := client.ReadOnlyTransaction().WithTimestampBound(spanner.ReadTimestamp(t))
	defer ro.Close()

	iter := ro.Read(ctx, "Albums", spanner.AllKeys(), []string{"SingerId", "AlbumId", "AlbumTitle"})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID string
		var albumID string
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", singerID, albumID, albumTitle)
	}
}