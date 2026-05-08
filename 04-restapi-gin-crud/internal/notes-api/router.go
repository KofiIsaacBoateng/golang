package notes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	// create repo and handler once
	repo := NewRepo(db);
	h := NewHandler(repo);


	// create routes group
	notesGroup := r.Group("/notes");

	// create routes
	{
		notesGroup.POST("", h.CreateNote)
		notesGroup.GET("", h.ListAllNotes)
		notesGroup.GET("/:id", h.GetNoteById)
		notesGroup.PUT("/:id", h.UpdateNoteById)
		notesGroup.DELETE("/:id", h.DeleteNote)
	}
}