package repository

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingRepository interface {
	CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error)
	FindPricingByID(id uuid.UUID) (*entity.Pricing, error)
	UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	DeletePricing(pricing *entity.Pricing) (bool, error)
	GetPricingByEventID(eventID uuid.UUID, pricingID uuid.UUID) (*entity.Pricing, error)
	UpdatePricingRemaining(pricingId uuid.UUID, remaining int) (*entity.Pricing, error)
}

type pricingRepository struct {
	db *gorm.DB
}

func NewPricingRepository(db *gorm.DB) PricingRepository {
	return &pricingRepository{db}
}

func (r *pricingRepository) CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
	if err := r.db.Create(pricing).Error; err != nil {
		return nil, err
	}
	return pricing, nil
}

func (r *pricingRepository) FindPricingByID(id uuid.UUID) (*entity.Pricing, error) {
	var pricing entity.Pricing
	if err := r.db.First(&pricing, "pricing_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *pricingRepository) FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error) {
	var priceList []entity.Pricing
	if err := r.db.Where("event_id = ?", eventID).Find(&priceList).Error; err != nil {
		return nil, err
	}
	return priceList, nil
}

func (r *pricingRepository) UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
	fields := setUpdateFieldPricing(pricing)

	if err := r.db.Model(&pricing).Updates(fields).Error; err != nil {
		return nil, err
	}

	return pricing, nil
}

func (r *pricingRepository) DeletePricing(pricing *entity.Pricing) (bool, error) {
	if err := r.db.Where("pricing_id = ?", pricing.PricingId).Delete(&pricing).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *pricingRepository) GetPricingByEventID(eventID uuid.UUID, pricingID uuid.UUID) (*entity.Pricing, error) {
	var pricing entity.Pricing
	if err := r.db.Where("event_id = ?", eventID).
		Where("pricing_id = ?", pricingID).
		First(&pricing).Error; err != nil {
		return nil, errors.New("pricing not found")
	}
	return &pricing, nil
}

func (r *pricingRepository) UpdatePricingRemaining(pricingId uuid.UUID, remaining int) (*entity.Pricing, error) {
	pricing := new(entity.Pricing)
	if err := r.db.Model(pricing).Where("pricing_id = ?", pricingId).Update("remaining", remaining).Error; err != nil {
		return nil, err
	}
	return pricing, nil
}
