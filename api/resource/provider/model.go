package provider

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DTO struct {
	ID             string `json:"id"`
	ProviderName   string `json:"provider_name"`
	IsActive       bool   `json:"is_active"`
	IsDefault      bool   `json:"is_default"`
	Configurations string `json:"configurations"`
}

type Form struct {
	ProviderName   string `json:"provider_name" form:"required,max=255"`
	IsActive       bool   `json:"is_active" form:"required,boolean"`
	IsDefault      bool   `json:"is_default" form:"required,boolean"`
	Configurations string `json:"configurations" validate:"required,json"`
}

type Provider struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProviderName   string         `gorm:"size:255;not null"`
	IsActive       bool           `gorm:"not null;default:true"`
	IsDefault      bool           `gorm:"not null;default:false"`
	Configurations string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Providers []*Provider

func (b *Provider) ToDto() *DTO {
	return &DTO{
		ID:             b.ID.String(),
		ProviderName:   b.ProviderName,
		IsActive:       b.IsActive,
		IsDefault:      b.IsDefault,
		Configurations: b.Configurations,
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

	return &Provider{
		ProviderName:   f.ProviderName,
		IsActive:       f.IsActive,
		IsDefault:      f.IsDefault,
		Configurations: f.Configurations,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
