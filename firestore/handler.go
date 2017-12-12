package firestore

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/trace"

	"github.com/google/uuid"

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
	span := tc.NewSpan("/firestore")
	defer span.FinishWait()
	ctx = trace.NewContext(ctx, span)
	do := grpc.WithUnaryInterceptor(tc.GRPCClientInterceptor())
	o := option.WithGRPCDialOption(do)

	client, err := firestore.NewClient(ctx, projectID, o)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create client: project=%s: %v", projectID, err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	id := uuid.New().String()
	row := Row{Value: "Hello Firestore"}
	doc, err := Put(ctx, client, id, row)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed put: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Done id = %s, doc.Path = %s", doc.ID, doc.Path)))
	return
}
