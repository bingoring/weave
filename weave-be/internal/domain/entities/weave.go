package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Weave domain entity - represents the core content unit
type Weave struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	ChannelID      uuid.UUID
	Title          string
	CoverImage     *string
	Content        WeaveContent
	Version        int
	ParentWeaveID  *uuid.UUID
	IsPublished    bool
	IsFeatured     bool
	ViewCount      int
	LikeCount      int
	ForkCount      int
	CommentCount   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// WeaveContent represents the structured content of a weave
type WeaveContent struct {
	Type string                 `json:"type"` // recipe, travel-plan, workout, etc.
	Data map[string]interface{} `json:"data"`
}

// Weave business methods
func (w *Weave) IsValidForPublication() bool {
	return w.Title != "" && w.Content.Type != "" && len(w.Content.Data) > 0
}

func (w *Weave) CanBeEditedBy(userID uuid.UUID) bool {
	return w.UserID == userID
}

func (w *Weave) CanBeFeatured() bool {
	return w.IsPublished && w.LikeCount >= 10 // Example criteria
}

func (w *Weave) IsForked() bool {
	return w.ParentWeaveID != nil
}

func (w *Weave) Publish() {
	w.IsPublished = true
	w.UpdatedAt = time.Now()
}

func (w *Weave) Unpublish() {
	w.IsPublished = false
	w.UpdatedAt = time.Now()
}

func (w *Weave) IncrementView() {
	w.ViewCount++
}

func (w *Weave) IncrementLike() {
	w.LikeCount++
}

func (w *Weave) DecrementLike() {
	if w.LikeCount > 0 {
		w.LikeCount--
	}
}

func (w *Weave) IncrementFork() {
	w.ForkCount++
}

func (w *Weave) UpdateContent(content WeaveContent) {
	w.Content = content
	w.Version++
	w.UpdatedAt = time.Now()
}

func (w *Weave) ToJSON() (string, error) {
	data, err := json.Marshal(w.Content)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewWeave(userID, channelID uuid.UUID, title string, content WeaveContent) *Weave {
	return &Weave{
		ID:          uuid.New(),
		UserID:      userID,
		ChannelID:   channelID,
		Title:       title,
		Content:     content,
		Version:     1,
		IsPublished: false,
		IsFeatured:  false,
		ViewCount:   0,
		LikeCount:   0,
		ForkCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func ForkWeave(originalWeave *Weave, newUserID uuid.UUID) *Weave {
	forkedWeave := &Weave{
		ID:            uuid.New(),
		UserID:        newUserID,
		ChannelID:     originalWeave.ChannelID,
		Title:         originalWeave.Title + " (Forked)",
		CoverImage:    originalWeave.CoverImage,
		Content:       originalWeave.Content,
		Version:       1, // Reset version for forked weave
		ParentWeaveID: &originalWeave.ID,
		IsPublished:   false, // Forked weaves start as drafts
		IsFeatured:    false,
		ViewCount:     0,
		LikeCount:     0,
		ForkCount:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Increment fork count on original
	originalWeave.IncrementFork()

	return forkedWeave
}