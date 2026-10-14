# edmEventsAPIGo

[![CI](https://github.com/weironiottan/edmEventsAPIGo/actions/workflows/ci.yml/badge.svg)](https://github.com/weironiottan/edmEventsAPIGo/actions/workflows/ci.yml)

A small Go HTTP API that serves EDM events (club, artist, date, ticket link) stored in Google Cloud Firestore.

## Requirements

- Go 1.27.2 or newer (see `go.mod`)
- A Google Cloud project with a Firestore database and an events collection

## Configuration

The app is configured entirely through environment variables:

| Variable               | Required | Description                                  |
| ---------------------- | -------- | -------------------------------------------- |
| `GOOGLE_CLOUD_PROJECT` | yes      | GCP project ID                               |
| `DATABASE_ID`          | yes      | Firestore database ID                        |
| `COLLECTION_NAME`      | yes      | Firestore collection holding the events      |
| `PORT`                 | no       | Port to listen on (defaults to `8080`)       |

Firestore credentials are picked up via [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials). Locally, run `gcloud auth application-default login`.

## Running locally

```sh
export GOOGLE_CLOUD_PROJECT=my-project
export DATABASE_ID=my-database
export COLLECTION_NAME=edm-events

go run ./cmd/api
```

Run it from the repository root so the HTML templates in `ui/html` can be found.

## Endpoints

| Method | Path                                    | Description                                 |
| ------ | --------------------------------------- | ------------------------------------------- |
| GET    | `/home`                                 | HTML landing page                           |
| GET    | `/v1/healthcheck`                       | Status and version                          |
| GET    | `/find-edm/all-events`                  | All events, sorted by event date            |
| GET    | `/find-edm/venue?clubname=<name>`       | Events at a given club                      |
| GET    | `/find-edm/artist?artistname=<name>`    | Events for a given artist                   |

Events are returned as JSON:

```json
{
  "Id": "...",
  "ClubName": "...",
  "ArtistName": "...",
  "EventDate": "...",
  "TicketUrl": "...",
  "ArtistImageUrl": "..."
}
```

## Development

These are the checks CI runs on every push and pull request to `main`:

```sh
go mod tidy -diff
gofmt -l .
go vet ./...
go build ./...
go test -race ./...
```
