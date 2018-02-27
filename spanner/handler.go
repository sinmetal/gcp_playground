package spanner

import (
	"net/http"
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc"
	"cloud.google.com/go/trace"
	"google.golang.org/api/option"
)

var projectID string

func SetUp(p string) {
	projectID = p
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
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

func doGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	tc, err := trace.NewClient(ctx, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span := tc.NewSpan("/spanner")
	defer span.FinishWait()
	 ctx = trace.NewContext(ctx, span)
	do := grpc.WithUnaryInterceptor(tc.GRPCClientInterceptor())
	o := option.WithGRPCDialOption(do)

	client, err := spanner.NewClient(ctx, "projects/souzoh-spanner-dev/instances/souzoh-shared-instance/databases/sinmetal", o)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create client: project=%s: %v", projectID, err), http.StatusInternalServerError)
		return
	}

	secParam := r.FormValue("sec")
	sec, err := strconv.Atoi(secParam)
	if err != nil {
		sec = 15
	}
	if err := readStaleData2(ctx, client, sec); err != nil {
		http.Error(w, fmt.Sprintf("readStaleData: project=%s: %v", projectID, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
	return
}