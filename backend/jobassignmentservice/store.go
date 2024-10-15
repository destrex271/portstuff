package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(ctx context.Context, connString string) (*Store, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) CreateTrip(ctx context.Context, trip Trip) (int, error) {
	var id int
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO trips (user_id, pickup_location, dropoff_location, request_time, status)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7)
		RETURNING trip_id
	`, trip.UserID, trip.PickupLocation.Lon, trip.PickupLocation.Lat, trip.DropoffLocation.Lon, trip.DropoffLocation.Lat, trip.RequestTime, trip.Status).Scan(&id)
	return id, err
}

func (s *Store) UpdateTrip(ctx context.Context, trip Trip) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE trips
		SET driver_id = $1, status = $2, start_time = $3, end_time = $4, fare = $5, distance = $6
		WHERE trip_id = $7
	`, trip.DriverID, trip.Status, trip.StartTime, trip.EndTime, trip.Fare, trip.Distance, trip.ID)
	return err
}

func (s *Store) GetTripByID(ctx context.Context, tripID int) (Trip, error) {
	var trip Trip
	err := s.db.QueryRowContext(ctx, `
		SELECT trip_id, user_id, driver_id, ST_X(pickup_location::geometry) as pickup_lon, ST_Y(pickup_location::geometry) as pickup_lat,
			   ST_X(dropoff_location::geometry) as dropoff_lon, ST_Y(dropoff_location::geometry) as dropoff_lat,
			   request_time, start_time, end_time, status, fare, distance
		FROM trips
		WHERE trip_id = $1
	`, tripID).Scan(&trip.ID, &trip.UserID, &trip.DriverID, &trip.PickupLocation.Lon, &trip.PickupLocation.Lat,
		&trip.DropoffLocation.Lon, &trip.DropoffLocation.Lat, &trip.RequestTime, &trip.StartTime, &trip.EndTime,
		&trip.Status, &trip.Fare, &trip.Distance)
	return trip, err
}

func (s *Store) FindNearestAvailableDriver(ctx context.Context, pickupLocation Point) (Driver, error) {
	var driver Driver
	// Calculating proximity using Postgis
	// But we need to switch to a more accurate implementation
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			d.id, 
			ST_X(d.last_known_location::geometry) AS lon, 
			ST_Y(d.last_known_location::geometry) AS lat,
			ST_Distance(
				d.last_known_location::geography, 
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
			) AS distance
		FROM 
			drivers d
		WHERE 
			d.is_available = true
			AND d.last_known_location IS NOT NULL
		ORDER BY 
			d.last_known_location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT 1
	`, pickupLocation.Lon, pickupLocation.Lat).Scan(
		&driver.ID,
		&driver.LastKnownLocation.Lon,
		&driver.LastKnownLocation.Lat,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Driver{}, fmt.Errorf("no available drivers found near the specified location")
		}
		return Driver{}, fmt.Errorf("error finding nearest available driver: %w", err)
	}

	return driver, nil
}

func (s *Store) UpdateDriverStatus(ctx context.Context, driverID int, isAvailable bool) error {
	_, err := s.db.ExecContext(ctx, "UPDATE drivers SET is_available = $1 WHERE id = $2", isAvailable, driverID)
	return err
}

func (s *Store) UpdateDriverLocation(ctx context.Context, driverID int, location Point) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE drivers
		SET last_known_location = ST_SetSRID(ST_MakePoint($1, $2), 4326)
		WHERE id = $3
	`, location.Lon, location.Lat, driverID)
	return err
}
