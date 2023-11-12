package main

import "go.mongodb.org/mongo-driver/mongo"

type SnippetModelInterface interface {
}

// SnippetModel Define a SnippetModel type which wraps a MongoDB connection pool.
type SnippetModel struct {
	DB         *mongo.Client
	collection *mongo.Collection
}
