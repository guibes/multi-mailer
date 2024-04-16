package provider

import (
	"time"

	"github.com/google/uuid"
)

type DTO struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ImageURL      string `json:"image_url"`
	Description   string `json:"description"`
}

type Form struct {
	Title         string `json:"title" form:"required,max=255"`
	Author        string `json:"author" form:"required,alpha_space,max=255"`
	PublishedDate string `json:"published_date" form:"required,datetime=2006-01-02"`
	ImageURL      string `json:"image_url" form:"url"`
	Description   string `json:"description"`
}

type Provider struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // UUID as primary key
	ProviderName   string    `gorm:"size:255;not null"`                               // Name of the provider
	IsActive       bool      `gorm:"not null;default:true"`                           // Whether the provider is active
	IsDefault      bool      `gorm:"not null;default:false"`                          // Whether the provider is the default choice
	Configurations string    `gorm:"type:text"`                                       // Configuration as JSON or serialized string
	CreatedAt      time.Time `gorm:"autoCreateTime"`                                  // Automatically handle created timestamp
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	DeletedAt      time.Time `gorm:"autoDeleteTime"` // Automatically handle updated timestamp
}

type Providers []*Provider

func (b *Provider) ToDto() *DTO {
	return &DTO{
		ID: b.ID.String(),
	}
}

func (bs Providers) ToDto() []*DTO {
	dtos := make([]*DTO, len(bs))
	for i, v := range bs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() *Provider {

	return &Provider{}
}
