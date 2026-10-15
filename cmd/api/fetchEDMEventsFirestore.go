package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// The scraper records each run in meta/{COLLECTION_NAME}
const metaCollection = "meta"

type SnippetModeler interface {
	FindAll() ([]EdmEvent, error)
	FindAllWithFilter(string, string) ([]EdmEvent, error)
	// LastUpdated returns when the scraper last stored events; found is false if it never has
	LastUpdated() (lastUpdated time.Time, found bool, err error)
}

type ScrapeMeta struct {
	LastUpdated time.Time `firestore:"lastUpdated"`
}

// SnippetModel Define a SnippetModel type which wraps a Firestore client.
type SnippetModel struct {
	Client     *firestore.Client
	Collection string
}

func (m *SnippetModel) FindAll() ([]EdmEvent, error) {
	ctx := context.Background()

	// Get a reference to the collection
	collRef := m.Client.Collection(m.Collection)

	// Create a query that sorts by eventdate in ascending order
	query := collRef.OrderBy("EventDate", firestore.Asc)

	// Execute the query
	iter := query.Documents(ctx)
	defer iter.Stop()

	var edmEvents []EdmEvent

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var edmEvent EdmEvent
		if err := doc.DataTo(&edmEvent); err != nil {
			return nil, err
		}

		edmEvents = append(edmEvents, edmEvent)
	}

	return edmEvents, nil
}

// This needs to be reworked later
func (m *SnippetModel) FindAllWithFilter(field string, value string) ([]EdmEvent, error) {
	ctx := context.Background()

	// Get a reference to the collection
	collRef := m.Client.Collection(m.Collection)

	query := collRef.Where(field, "==", value).OrderBy("EventDate", firestore.Asc)

	// Execute the query
	iter := query.Documents(ctx)
	defer iter.Stop()

	var edmEvents []EdmEvent

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var edmEvent EdmEvent
		if err := doc.DataTo(&edmEvent); err != nil {
			return nil, err
		}

		edmEvents = append(edmEvents, edmEvent)
	}

	return edmEvents, nil
}

func (m *SnippetModel) LastUpdated() (time.Time, bool, error) {
	doc, err := m.Client.Collection(metaCollection).Doc(m.Collection).Get(context.Background())
	if status.Code(err) == codes.NotFound {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}

	var meta ScrapeMeta
	if err := doc.DataTo(&meta); err != nil {
		return time.Time{}, false, err
	}
	if meta.LastUpdated.IsZero() {
		return time.Time{}, false, nil
	}
	return meta.LastUpdated, true, nil
}
