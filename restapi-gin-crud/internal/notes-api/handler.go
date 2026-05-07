package notes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo ) *Handler {
	return &Handler {
		repo: repo,
	}
}


func (h *Handler) CreateNote(c *gin.Context) {
	// bind request struct reference to the request body NoteReqBody
	var req NoteReqBody;

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": "Invalid JSON!",
		})
		return;
	}

	now := time.Now().UTC();

	note := Note{
		ID: primitive.NewObjectID(),
		Title: req.Title,
		Content: req.Content,
		Pinned: req.Pinned,
		CreatedAt: now,
		UpdatedAt: now,
	}

	note, err := h.repo.Create(c.Request.Context(),note )
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": "Failed to create note",
		})

		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"msg": "Note created successfully!",
		"note": note,
	})
}


func (h *Handler) ListAllNotes(c *gin.Context) {
	notes, err := h.repo.ListAll(c.Request.Context());
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": fmt.Errorf("Failed to fetch all notes: %w", err),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": notes,
		"msg": "Fetched all lists successfully!",
	})
}

func (h *Handler) GetNoteById(c *gin.Context) {
	id, exists := c.Params.Get("id");
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": "Param id is missing",
		})
		return;
	}

	// convert id
	noteId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Invalid ObjectID: %w", err),
		})
		return;
	}

	note, err := h.repo.GetById(c.Request.Context(), noteId);
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"ok": false,
			"error": fmt.Errorf("Note not Found: %w", err),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": note,
		"msg": "Note retrieved successfully!",
	})
}

func (h *Handler) UpdateNoteById(c *gin.Context) {
	id, exists := c.Params.Get("id");
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Param[id] is missing!"),
		})
		return;
	}

	// convert id 
	noteId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Invalid ObjectID: %w", err),
		})
		return;
	}

	// bind json to request body
	var reqBody UpdateReqBody;

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Invalid JSON from request body: %w", err),
		})
		return;
	}

	note, err := h.repo.FindAndUpdate(c.Request.Context(), noteId, reqBody);
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"ok": false,
			"error": fmt.Errorf("Failed to update note: %w", err),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": note,
		"msg": "Note updated successfully!",
	})
}

func (h *Handler) DeleteNote(c *gin.Context) {
	id, exists := c.Params.Get("id");
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Param[id] is missing!"),
		})
		return;
	}

	// convert id 
	noteId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Invalid ObjectID: %w", err),
		})
		return;
	}

	deleted, err := h.repo.DeleteNote(c.Request.Context(), noteId);
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": fmt.Errorf("Delete failed: %w", err),
		})
		return;
	}

	if !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": fmt.Errorf("Not found!: %w", err),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"msg": "Note deleted successfully!",
	})
}