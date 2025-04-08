package main

//
// Archiving this code
//import (
//	"context"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//type SnippetModelInterface interface {
//	FindAll() ([]EdmEvent, error)
//	FindAllWithFilter(string, string) ([]EdmEvent, error)
//}
//
//// SnippetModel Define a SnippetModel type which wraps a MongoDB connection pool.
//type SnippetModel struct {
//	DB         *mongo.Client
//	collection *mongo.Collection
//}
//
//func (m *SnippetModel) FindAll() ([]EdmEvent, error) {
//	filter := bson.D{}
//	// Sorts eventdate in Ascending order
//	opts := options.Find().SetSort(bson.D{{"eventdate", 1}})
//
//	cursor, err := m.collection.Find(context.TODO(), filter, opts)
//	if err != nil {
//		return nil, err
//	}
//	var edmEvents []EdmEvent
//
//	err = cursor.All(context.TODO(), &edmEvents)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return edmEvents, nil
//}
//
//func (m *SnippetModel) FindAllWithFilter(query string, value string) ([]EdmEvent, error) {
//
//	// creates an index so that we can search the field via text
//	model := mongo.IndexModel{Keys: bson.D{{query, "text"}}}
//	_, err := m.collection.Indexes().CreateOne(context.TODO(), model)
//
//	// search the given value in the given field provided by query
//	filter := bson.D{{"$text", bson.D{{"$search", value}}}}
//
//	// Sorts eventdate in Ascending order
//	opts := options.Find().SetSort(bson.D{{"eventdate", 1}})
//
//	cursor, err := m.collection.Find(context.TODO(), filter, opts)
//	if err != nil {
//		return nil, err
//	}
//	_, err = m.collection.Indexes().DropAll(context.TODO())
//	if err != nil {
//		return nil, err
//	}
//
//	var edmEvents []EdmEvent
//
//	err = cursor.All(context.TODO(), &edmEvents)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return edmEvents, nil
//}
