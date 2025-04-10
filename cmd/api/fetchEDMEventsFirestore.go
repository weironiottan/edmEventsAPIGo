package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
)

type SnippetModeler interface {
	FindAll() ([]EdmEvent, error)
	FindAllWithFilter(string, string) ([]EdmEvent, error)
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
