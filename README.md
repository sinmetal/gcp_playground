# gcp_playground

## Cloud Bigtable

### Local Run

#### Bigtable Emulator Start

```
gcloud beta emulators bigtable start
$(gcloud beta emulators bigtable env-init)
```

#### Datastore Emulator Start

```
gcloud beta emulators datastore start
$(gcloud beta emulators datastore env-init)
```

#### My Application Run

```
go run *.go --project=myproject --bigtableInstance=sample
```