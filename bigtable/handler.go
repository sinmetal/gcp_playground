package bigtable

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"cloud.google.com/go/bigtable"
)

var projectID string
var instance string

const table = "Item"
const family = "myfamily"

func SetUp(projectID string, instance string) {
	projectID = projectID
	instance = instance
}

func HandlerBigtable(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		doPost(w, r)
	} else if r.Method == "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
	} else if r.Method == "PUT" {
		doPut(w, r)
	} else if r.Method == "DELETE" {
		http.Error(w, "", http.StatusMethodNotAllowed)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func doPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create admin client: project=%s, instance=%s : %v", projectID, instance, err), http.StatusInternalServerError)
		return
	}
	defer adminClient.Close()

	err = CreateTableWithColumnFamily(ctx, adminClient, table, family)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done Create Table with ColumnFamily!"))
	return
}

func doPut(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	client, err := bigtable.NewClient(ctx, projectID, instance)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed Bigtable.NewClient(): projectID=%s, instance=%s", projectID, instance), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	err = UpdateBigtable(ctx, client, table, family, "mycolumn")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done Update!"))
	return
}
