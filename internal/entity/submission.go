
package entity

import (
    "github.com/google/uuid"
)

type Submission struct {
    ID        uuid.UUID `json:"id"`
    EventID   uuid.UUID `json:"event_id"`
    UserID    uuid.UUID `json:"user_id"`
    Name      string    `json:"name"`
    Filename  string    `json:"filename"`
    Status    string    `json:"status,omitempty"`
    Type      string    `json:"type"`
    Auditable
}

func NewSubmission(eventID, userID uuid.UUID, name, filename, submissionType, status string) *Submission {
    return &Submission{
        ID:        uuid.New(),
        EventID:   eventID,
        UserID:    userID,
        Name:      name,
        Filename:  filename,
        Type:      submissionType,
        Status:    status,
        Auditable: NewAuditable(),
    }
}
