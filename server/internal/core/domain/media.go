package domain

import "time"

type MediaAsset struct {
	ID        int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string                 `gorm:"unique;not null" json:"uuid"`
	Disk      string                 `gorm:"not null;default:'s3_main';size:50" json:"disk"`
	Path      string                 `gorm:"not null;size:512" json:"path"`
	Filename  string                 `gorm:"not null;size:255" json:"filename"`
	MimeType  string                 `gorm:"not null;size:100" json:"mime_type"`
	SizeBytes int64                  `gorm:"not null" json:"size_bytes"`
	Width     *int                   `json:"width"`
	Height    *int                   `json:"height"`
	Blurhash  *string                `gorm:"size:100" json:"blurhash"`
	AltText   *string                `gorm:"size:255" json:"alt_text"`
	AITags    map[string]interface{} `gorm:"type:jsonb" json:"ai_tags"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type MediaLink struct {
	ID         int         `gorm:"primaryKey;autoIncrement" json:"id"`
	MediaID    int         `gorm:"not null" json:"media_id"`
	Media      *MediaAsset `gorm:"foreignKey:MediaID" json:"media"`
	EntityType string      `gorm:"not null;size:50" json:"entity_type"`
	EntityID   int         `gorm:"not null" json:"entity_id"`
	Zone       string      `gorm:"not null;default:'gallery';size:50" json:"zone"`
	SortOrder  int         `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time   `json:"created_at"`
}
