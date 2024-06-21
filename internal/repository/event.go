package repository

import (
	"time"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
    CreateEvent(event *entity.Event,  pricings []entity.Pricing) (*entity.Event, error)
    CreatePricing(pricing *entity.Pricing)  (*entity.Pricing, error)
    FindEventByID(id uuid.UUID) (*entity.Event, error)
    FindPricingByID(id uuid.UUID) (*entity.Pricing, error)
    FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error)
    UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
    UpdateEvent(event *entity.Event) (*entity.Event, error)
    UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
    DeleteEvent(event *entity.Event) (bool, error)
    DeletePricing(pricing *entity.Pricing) (bool, error)
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


func (r *eventRepository) CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
    if err := r.db.Create(pricing).Error; err != nil {
        return nil, err
    }
    return pricing, nil
}

func (r *eventRepository) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Preload("Pricings").First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
func (r *eventRepository) FindPricingByID(id uuid.UUID) (*entity.Pricing, error) {
	var pricing entity.Pricing
	if err := r.db.First(&pricing, "pricing_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *eventRepository) FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error) {
	var priceList []entity.Pricing
	if err := r.db.Where("event_id = ?",eventID).Find(&priceList).Error; err != nil {
		return nil, err
	}
	return priceList, nil
}


func (r *eventRepository) UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {

    fieldsEvent := setUpdateFieldEvent(event)

    
    r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&event).Updates(fieldsEvent).Error; err != nil {
            tx.Rollback()
            return err
        }
        for _, p := range pricings{
            fieldsPricing := setUpdateFieldPricing(&p)

            if err := tx.Model(&p).Where("pricing_id = ? AND event_id = ?", p.PricingId, event.ID).Updates(fieldsPricing).Error; err != nil {
                tx.Rollback()
                return err
            }
        }

        return nil
    })
    
    if err := r.db.Where("id = ?", event.ID).Preload("Pricings").Find(&event).Error; err != nil {
        return nil, err
    }

    return event, nil
}


func (r *eventRepository) UpdateEvent(event *entity.Event) (*entity.Event, error) {
    fields := setUpdateFieldEvent(event)

    if err := r.db.Model(&event).Updates(fields).Error; err != nil {
        return nil, err
    }

    return event, nil
}


func (r *eventRepository) UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
    fields := setUpdateFieldPricing(pricing)

    if err := r.db.Model(&pricing).Updates(fields).Error; err != nil {
        return nil, err
    }

    return pricing, nil
}

func (r *eventRepository) DeleteEvent(event *entity.Event) (bool, error) {
    if err := r.db.Delete(&event).Error; err != nil {
		return false, err
	}

    return true, nil
}
func (r *eventRepository) DeletePricing(pricing *entity.Pricing) (bool, error) {
    if err := r.db.Where("pricing_id = ?", pricing.PricingId).Delete(&pricing).Error; err != nil {
		return false, err
	}

    return true, nil
}



//function tambahan
func setUpdateFieldEvent(event *entity.Event) map[string]interface{} {
    fieldsEvent := make(map[string]interface{})


    if event.Name != "" {
        fieldsEvent["name"] = event.Name
    }

    if event.Description != "" {
        fieldsEvent["description"] = event.Description
    }

    if event.Organizer != "" {
        fieldsEvent["organizer"] = event.Organizer
    }

    if event.StartTime != (time.Time{}) {
        fieldsEvent["start_date"] = event.StartTime
    }

    if event.EndTime != (time.Time{}) {
        fieldsEvent["end_date"] = event.EndTime
    }

    
    // if event.MustUploadSubmission !=  {
    //     fieldsEvent["must_upload_submission"] = event.MustUploadSubmission
    // }

    if event.Province != "" {
        fieldsEvent["province"] = event.Province
    }

    if event.City != "" {
        fieldsEvent["city"] = event.City
    }

    if event.District != "" {
        fieldsEvent["district"] = event.District
    }

    if event.Longitude != 0.0 {
        fieldsEvent["longitude"] = event.Longitude
    }

    if event.Latitude != 0.0 {
        fieldsEvent["latitude"] = event.Latitude
    }

    if event.FullAddress != "" {
        fieldsEvent["full_address"] = event.FullAddress
    }

    return fieldsEvent
}

func setUpdateFieldPricing(pricing *entity.Pricing) map[string]interface{} {
    fields := make(map[string]interface{})

    if pricing.Fee != 0.0 {
        fields["fee"] = pricing.Fee
    }

    if pricing.Quota != 0 {
        fields["quota"] = pricing.Quota
    }

    if pricing.Remaining != 0 {
        fields["remaining"] = pricing.Remaining
    }

    if pricing.Name != "" {
        fields["name"] = pricing.Name
    }

    return fields
}