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

func HandlerBigbang(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bigbang(w, r)
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

func bigbang(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	tc, err := trace.NewClient(ctx, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span := tc.NewSpan("/bigtable/bigbang")
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

	for i := 0; i < 1000; i++ {
		err = UpdateBigtable(ctx, client, table, family, "mycolumn")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done Bigbang!"))
	return
}
