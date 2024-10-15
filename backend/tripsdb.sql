-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create enum types for status fields
CREATE TYPE trip_status AS ENUM ('requested', 'assigned', 'in_progress', 'completed', 'cancelled');
CREATE TYPE payment_status AS ENUM ('pending', 'processing', 'completed', 'failed');

-- Trips table
CREATE TABLE trips (
    trip_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    driver_id INT,
    pickup_location GEOGRAPHY(POINT, 4326) NOT NULL,  -- Using Geography type for pickup location
    dropoff_location GEOGRAPHY(POINT, 4326) NOT NULL, -- Using Geography type for dropoff location
    request_time TIMESTAMP NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    status trip_status NOT NULL DEFAULT 'requested',
    fare DECIMAL(10, 2),
    distance DECIMAL(10, 2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trip assignments table
CREATE TABLE trip_assignments (
    assignment_id SERIAL PRIMARY KEY,
    trip_id INT NOT NULL REFERENCES trips(trip_id),
    driver_id INT NOT NULL,
    assignment_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trip locations table (for tracking real-time locations)
CREATE TABLE trip_locations (
    location_id SERIAL PRIMARY KEY,
    trip_id INT NOT NULL REFERENCES trips(trip_id),
    location GEOGRAPHY(POINT, 4326) NOT NULL,  -- Storing location as Geography type
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Payments table
CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    trip_id INT NOT NULL REFERENCES trips(trip_id),
    amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50),
    status payment_status NOT NULL DEFAULT 'pending',
    transaction_id VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trip ratings table
CREATE TABLE trip_ratings (
    rating_id SERIAL PRIMARY KEY,
    trip_id INT NOT NULL REFERENCES trips(trip_id),
    user_rating INT CHECK (user_rating BETWEEN 1 AND 5),
    driver_rating INT CHECK (driver_rating BETWEEN 1 AND 5),
    user_comment TEXT,
    driver_comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_trips_user_id ON trips(user_id);
CREATE INDEX idx_trips_driver_id ON trips(driver_id);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trip_assignments_driver_id ON trip_assignments(driver_id);
CREATE INDEX idx_trip_locations_trip_id ON trip_locations(trip_id);
CREATE INDEX idx_trip_locations_location ON trip_locations(location);  -- Index on location for spatial queries
CREATE INDEX idx_payments_trip_id ON payments(trip_id);
CREATE INDEX idx_trip_ratings_trip_id ON trip_ratings(trip_id);

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers to automatically update the updated_at column
CREATE TRIGGER update_trips_modtime
    BEFORE UPDATE ON trips
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_trip_assignments_modtime
    BEFORE UPDATE ON trip_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_payments_modtime
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
