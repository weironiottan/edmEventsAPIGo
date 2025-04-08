package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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
	query := collRef.OrderBy("eventdate", firestore.Asc)

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

func (m *SnippetModel) FindAllWithFilter(field string, value string) ([]EdmEvent, error) {
	ctx := context.Background()

	// Get a reference to the collection
	collRef := m.Client.Collection(m.Collection)

	query := collRef.Where(field, "==", value).OrderBy("eventdate", firestore.Asc)

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

// Alternative initialization with credentials file
func NewSnippetModelWithCredentials(ctx context.Context, projectID, credentialsFile, collection string) (*SnippetModel, error) {
	// Create a Firestore client with credentials
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &SnippetModel{
		Client:     client,
		Collection: collection,
	}, nil
}
