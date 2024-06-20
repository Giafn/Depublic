package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"gorm.io/gorm"
	// "github.com/google/uuid"
)

type EventRepository interface {
    CreateEvent(event *entity.Event,  pricings []entity.Pricing) (*entity.Event, error)
    CreatePricing(pricing *entity.Pricing) error
}

type eventRepository struct {
    db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
    return &eventRepository{db}
}

func (r *eventRepository) CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
    

    err := r.db.Transaction(func(tx *gorm.DB) error {
       
        if err := tx.Create(event).Error; err != nil {
            tx.Rollback()
            return err
        }
        
        for _, p := range pricings{
            p.EventID = event.ID
    
            if err := tx.Create(&p).Error; err != nil {
                tx.Rollback()
                return err
            }
        }
    
        return nil
      })

      if err != nil { 
        return nil, err
       }
      

	if err := r.db.Where("id = ?", event.ID).Preload("Pricings").Find(&event).Error; err != nil {
        return nil, err
    }

    return event, nil
}


func (r *eventRepository) CreatePricing(pricing *entity.Pricing) error {
    return r.db.Create(pricing).Error
}