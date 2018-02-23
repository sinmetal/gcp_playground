package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
		doGet(w, r)
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

	chErr := make(chan error)
	id := uuid.New().String()
	row := Row{Value: "Hello Firestore", Number1: rand.Int(), Number2: rand.Int(), CreatedAt: time.Now()}
	go func() {
		_, err := Put(ctx, client, id, row)
		if err != nil {
			fmt.Printf("failed put: %v", err, http.StatusInternalServerError)
		}
		chErr <- err
	}()
	go func() {
		err := PutTxMulti(ctx, client, id, row)
		if err != nil {
			fmt.Printf("failed putTxMulti: %v", err, http.StatusInternalServerError)
		}
		chErr <- err
	}()
	<-chErr
	<-chErr

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
	return
}

func doGet(w http.ResponseWriter, r *http.Request) {
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

	l, err := List(ctx, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var rows []*Row
	for _, v := range l {
		r := Row{}
		err := v.DataTo(&r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows = append(rows, &r)
	}

	b, err := json.Marshal(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	w.Write(b)
	return
}
