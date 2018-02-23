package datastore

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"cloud.google.com/go/trace"
	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var projectID string

func SetUp(p string) {
	projectID = p
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		doPost(w, r)
	} else if r.Method == "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
	} else if r.Method == "PUT" {
		http.Error(w, "", http.StatusMethodNotAllowed)
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
	span := tc.NewSpan("/datastore")
	defer span.FinishWait()
	ctx = trace.NewContext(ctx, span)
	do := grpc.WithUnaryInterceptor(tc.GRPCClientInterceptor())
	o := option.WithGRPCDialOption(do)

	client, err := datastore.NewClient(ctx, projectID, o)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create client: project=%s: %v", projectID, err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	row := Row{Value: "Hello Datastore"}
	key, err := Put(ctx, client, row)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed put: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Done id = %d, key.Encode = %s", key.ID, key.Encode())))
	return
}
