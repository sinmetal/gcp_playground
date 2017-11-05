# gcp_playground

## Cloud Bigtable

### Local Run

```
gcloud beta emulators bigtable start
$(gcloud beta emulators bigtable env-init)

go run *.go --project=myproject --bigtableInstance=sample
```