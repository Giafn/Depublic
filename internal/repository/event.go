package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type EventRepository interface {
    CreateEvent(event *entity.Event,  pricings []entity.Pricing) (*entity.Event, error)
    CreatePricing(pricing *entity.Pricing) error
    FindEventByID(id uuid.UUID) (*entity.Event, error)
    FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error)
    GetEvents(filter map[string] interface{}, sort string)([]entity.Event, error)
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

func (r *eventRepository) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Preload("Pricings").First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindPricingByEventID(eventID uuid.UUID) ([]entity.Pricing, error) {
	var priceList []entity.Pricing
	if err := r.db.Where("event_id = ?",eventID).Find(&priceList).Error; err != nil {
		return nil, err
	}
	return priceList, nil
}

func (r *eventRepository) GetEvents(filter map[string] interface{}, sort string)([]entity.Event, error){
    var events []entity.Event

    query := r.db.Preload("Pricings")

    // Apply filters
    if price, ok := filter["price"]; ok {
        query = query.Joins("JOIN pricings ON pricings.event_id = events.id")
        switch price {
        case "0":
            query = query.Where("pricings.fee = ?", 0)
        case "<100000":
            query = query.Where("pricings.fee < ?", 100000)
        case "<500000":
            query = query.Where("pricings.fee < ?", 500000)
        case "<1000000":
            query = query.Where("pricings.fee < ?", 1000000)
        case "<2500000":
            query = query.Where("pricings.fee < ?", 2500000)
        case "<5000000":
            query = query.Where("pricings.fee < ?", 5000000)
        case ">5000000":
            query = query.Where("pricings.fee > ?", 5000000)
        }
    }

    if province, ok := filter["province"]; ok {
        query = query.Where("province = ?", province)
    }

    if city, ok := filter["city"]; ok {
        query = query.Where("city = ?", city)
    }

    if timeStart, ok := filter["timeStart"]; ok {
        switch timeStart {
        case "day":
            query = query.Where("(EXTRACT(HOUR FROM start_time) >= ? AND EXTRACT(HOUR FROM start_time) <= ?)", 6, 17)
        case "night":
            query = query.Where("(EXTRACT(HOUR FROM start_time) >= ? OR EXTRACT(HOUR FROM start_time) <= ?)", 18, 5)
        }
    }

    if category, ok := filter["category"]; ok {
        switch category {
        case "withSubmission":
            query = query.Where("must_upload_submission = ?", true)
        case "withoutSubmission":
            query = query.Where("must_upload_submission = ?", false)
        }
    }

    // Apply sorting
    switch sort {
    case "terbaru":
        query = query.Order("created_at DESC")
    case "termahal":
        query = query.Joins("JOIN pricings ON pricings.event_id = events.id").
            Select("events.*, MAX(pricings.fee) as max_fee").
            Group("events.id").
            Order("max_fee DESC")
    case "termurah":
        query = query.Joins("JOIN pricings ON pricings.event_id = events.id").
            Select("events.*, MIN(pricings.fee) as min_fee").
            Group("events.id").
            Order("min_fee ASC")
    case "terdekat":
        radius := 6371
        lat := 0.5069023316903872
        lon := 101.3985472713651
        query = query.Select("events.*, "+
        "(? * ACOS(SIN(RADIANS(?)) * SIN(RADIANS(latitude)) + COS(RADIANS(?)) * COS(RADIANS(latitude)) * COS(RADIANS(?) - RADIANS(longitude)))) as distance", radius, lat, lat, lon).
        Order("distance ASC")
    }

    if err := query.Find(&events).Error; err != nil {
        return nil, err
    }

    return events, nil
}   
