package notes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repo struct {
	coll *mongo.Collection
}


func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
			coll: db.Collection("notes"),
		}
	
}

func (r *Repo) Create(ctx context.Context, note Note)(Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second);

	defer cancel();

	_, err := r.coll.InsertOne(opCtx, note);
	if err != nil {
		return Note{}, fmt.Errorf("Failed to insert note: %w", err)
	}

	return note, nil;
}

func (r *Repo) ListAll(ctx context.Context)([]Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second);
	defer cancel();

	filter := bson.M{}; // match all

	// Find returns a cursor -> like an iterator over all matching elements
	cursor, err := r.coll.Find(opCtx, filter);
	if err != nil {
		return nil, fmt.Errorf("Failed to list all notes: %w", err)
	}

	// cursor must be closed after usage -> mean to avoid data leaks;
	defer cursor.Close(opCtx);

	var notes []Note;

	// decode cursor and set the notes slice pointer
	if err := cursor.All(opCtx, &notes); err != nil {
		return nil, fmt.Errorf("Notes list decode failed: %w", err);
	}


	return notes, nil;
}

func (r *Repo) GetById(ctx context.Context, id primitive.ObjectID) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second);
	defer cancel()


	filter := bson.M{"_id": id};

	var note Note;

	if err := r.coll.FindOne(opCtx, filter, options.FindOne()).Decode(&note); err != nil {
		return Note{}, fmt.Errorf("Note failed to find note: %w", err)
	}

	return note, nil
}

func (r *Repo) FindAndUpdate(ctx context.Context, id primitive.ObjectID, note UpdateReqBody) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second);
	defer cancel();

	filter := bson.M{"_id": id};

	update := bson.M{
		"$set": bson.M{
			"title": note.Title,
			"content": note.Content,
			"pinned": note.Pinned,
			"updatedAt": time.Now().UTC(),
		},
	}

	after := options.After;
	opts := options.FindOneAndUpdate().SetReturnDocument(after)
	
	var updatedNote Note;

	err := r.coll.FindOneAndUpdate(opCtx, filter, update, opts).Decode(&updatedNote); 
	if err != nil {
		return Note{}, fmt.Errorf("Failed to update Note: %w", err)
	}


	return updatedNote, nil
}

func (r *Repo) DeleteNote(ctx context.Context, id primitive.ObjectID) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second);
	defer cancel();


	filter := bson.M{"_id": id};

	res, err := r.coll.DeleteOne(opCtx, filter);

	if err != nil {
		return false, fmt.Errorf("Failed to delete note: %w", err)
	}

	if(res.DeletedCount == 0){
		return false, nil;
	}

	return true, nil;
}