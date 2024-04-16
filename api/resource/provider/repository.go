package provider

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Providers, error) {
	providers := make([]*Provider, 0)
	if err := r.db.Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (r *Repository) Create(provider *Provider) (*Provider, error) {
	if err := r.db.Create(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (r *Repository) Read(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	if err := r.db.Where("id = ?", id).First(&provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (r *Repository) Update(provider *Provider) (int64, error) {
	result := r.db.Model(&Provider{}).
		Select("Title", "Author", "PublishedDate", "ImageURL", "Description", "UpdatedAt").
		Where("id = ?", provider.ID).
		Updates(provider)

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&Provider{})

	return result.RowsAffected, result.Error
}
