package main

import (
	"context"
	"errors"
	"time"
)

// TripRequest represents a request for a new trip.
type TripRequest struct {
	UserID     int     `json:"user_id"`
	PickupLat  float64 `json:"pickup_lat"`
	PickupLon  float64 `json:"pickup_lon"`
	DropoffLat float64 `json:"dropoff_lat"`
	DropoffLon float64 `json:"dropoff_lon"`
}

// Validate checks the validity of the TripRequest coordinates.
func (r *TripRequest) Validate() error {
	if r.PickupLat < -90 || r.PickupLat > 90 || r.DropoffLat < -90 || r.DropoffLat > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if r.PickupLon < -180 || r.PickupLon > 180 || r.DropoffLon < -180 || r.DropoffLon > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

// Trip represents a trip in the system.
type Trip struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	DriverID        *int       `json:"driver_id,omitempty"`
	PickupLocation  Point      `json:"pickup_location"`
	DropoffLocation Point      `json:"dropoff_location"`
	RequestTime     time.Time  `json:"request_time"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	Status          TripStatus `json:"status"`
	Fare            *float64   `json:"fare,omitempty"`
	Distance        *float64   `json:"distance,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"` // Track last update time
}

// Point represents a geographical point with latitude and longitude.
type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// TripStatus represents the status of a trip.
type TripStatus string

// Trip statuses
const (
	TripStatusRequested  TripStatus = "requested"
	TripStatusAssigned   TripStatus = "assigned"
	TripStatusInProgress TripStatus = "in_progress"
	TripStatusCompleted  TripStatus = "completed"
	TripStatusCancelled  TripStatus = "cancelled"
)

// Driver represents a driver in the system.
type Driver struct {
	ID                int   `json:"id"`
	LastKnownLocation Point `json:"last_known_location"`
	IsAvailable       bool  `json:"is_available"`
}

// JobAssignmentService defines methods for managing trip requests.
type JobAssignmentService interface {
	RequestTrip(context.Context, TripRequest) error
	AssignTrip(context.Context, int) error
	UpdateTripStatus(context.Context, int, TripStatus) error
	GetTripByID(context.Context, int) (Trip, error)
	ProcessTripRequest(context.Context, TripRequest) error
}

// JobAssignmentStore defines methods for data persistence.
type JobAssignmentStore interface {
	CreateTrip(context.Context, Trip) (int, error)
	UpdateTrip(context.Context, Trip) error
	GetTripByID(context.Context, int) (Trip, error)
	FindNearestAvailableDriver(context.Context, Point) (Driver, error)
	UpdateDriverStatus(context.Context, int, bool) error
	UpdateDriverLocation(context.Context, int, Point) error
}

// KafkaProducer defines methods for producing Kafka messages.
type KafkaProducer interface {
	ProduceTripRequest(context.Context, TripRequest) error
}

// KafkaConsumer defines methods for consuming Kafka messages.
type KafkaConsumer interface {
	ConsumeTripRequests(context.Context) (<-chan TripRequest, error)
}
