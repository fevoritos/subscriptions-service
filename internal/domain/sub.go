package subscription

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionID uuid.UUID
type UserID uuid.UUID

type Subscription struct {
	ID          SubscriptionID
	UserID      UserID
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
}

func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(uuid.New())
}

func (id SubscriptionID) String() string {
	return uuid.UUID(id).String()
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

func ParseSubscriptionID(s string) (SubscriptionID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return SubscriptionID{}, err
	}
	return SubscriptionID(u), nil
}

func ParseID(s string) (UserID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, err
	}

	return UserID(u), err
}

func (id SubscriptionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}
