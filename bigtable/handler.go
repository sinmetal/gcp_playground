package bigtable

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/trace"
)

var projectID string
var instance string

const table = "Item"
const family = "myfamily"

func SetUp(p string, i string) {
	projectID = p
	instance = i
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

	tc, err := trace.NewClient(ctx, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span := tc.NewSpan("/bigtable")
	defer span.FinishWait()
	ctx = trace.NewContext(ctx, span)
	do := grpc.WithUnaryInterceptor(tc.GRPCClientInterceptor())
	o := option.WithGRPCDialOption(do)
	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instance, o)
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

	tc, err := trace.NewClient(ctx, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span := tc.NewSpan("/bigtable")
	defer span.FinishWait()
	ctx = trace.NewContext(ctx, span)
	do := grpc.WithUnaryInterceptor(tc.GRPCClientInterceptor())
	o := option.WithGRPCDialOption(do)
	client, err := bigtable.NewClient(ctx, projectID, instance, o)
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
