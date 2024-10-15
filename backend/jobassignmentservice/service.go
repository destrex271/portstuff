package main

import (
	"context"
	"time"
)

type Service struct {
	store     JobAssignmentStore
	publisher KafkaProducer
}

func NewService(store JobAssignmentStore, publisher KafkaProducer) *Service {
	return &Service{store: store, publisher: publisher}
}

func (s *Service) RequestTrip(ctx context.Context, req TripRequest) error {
	return s.publisher.ProduceTripRequest(ctx, req)
}

func (s *Service) ProcessTripRequest(ctx context.Context, req TripRequest) error {
	trip := Trip{
		UserID:          req.UserID,
		PickupLocation:  Point{Lat: req.PickupLat, Lon: req.PickupLon},
		DropoffLocation: Point{Lat: req.DropoffLat, Lon: req.DropoffLon},
		RequestTime:     time.Now(),
		Status:          TripStatusRequested,
	}

	tripID, err := s.store.CreateTrip(ctx, trip)
	if err != nil {
		return err
	}

	return s.AssignTrip(ctx, tripID)
}

func (s *Service) AssignTrip(ctx context.Context, tripID int) error {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	driver, err := s.store.FindNearestAvailableDriver(ctx, trip.PickupLocation)
	if err != nil {
		return err
	}

	trip.DriverID = &driver.ID
	trip.Status = TripStatusAssigned

	if err := s.store.UpdateTrip(ctx, trip); err != nil {
		return err
	}

	return s.store.UpdateDriverStatus(ctx, driver.ID, false)
}

func (s *Service) UpdateTripStatus(ctx context.Context, tripID int, status TripStatus) error {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	trip.Status = status

	if status == TripStatusCompleted && trip.DriverID != nil {
		if err := s.store.UpdateDriverStatus(ctx, *trip.DriverID, true); err != nil {
			return err
		}
	}

	return s.store.UpdateTrip(ctx, trip)
}

func (s *Service) GetTripByID(ctx context.Context, tripID int) (Trip, error) {
	return s.store.GetTripByID(ctx, tripID)
}
