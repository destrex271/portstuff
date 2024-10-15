CREATE TABLE Users (
    Id  BIGSERIAL PRIMARY KEY,
    Username VARCHAR(255) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    Dob DATE NOT NULL,
    EmailID VARCHAR(255) NOT NULL UNIQUE,
    Mobile VARCHAR(20) NOT NULL
);

CREATE TABLE Admins (
    Role VARCHAR(100) NOT NULL  -- Additional admin-specific column
    -- Additional admin-specific columns can be added here
) INHERITS (Users);

CREATE TABLE Drivers (
    LicenseNumber VARCHAR(50) NOT NULL UNIQUE,
    AdminId INT,  -- Foreign key to track which Admin manages this driver
    -- Additional driver-specific columns can be added here
    FOREIGN KEY (AdminId) REFERENCES Users(Id) ON DELETE SET NULL
) INHERITS (Users);

CREATE TABLE Vehicles (
    Id SERIAL PRIMARY KEY,
    LicensePlate VARCHAR(20) NOT NULL UNIQUE,
    Model VARCHAR(100) NOT NULL,
    Year INT NOT NULL CHECK (Year > 1885),  -- Cars were invented in the late 19th century
    DriverId INT,  -- Foreign key to track which driver is associated with this vehicle
    -- Additional vehicle-specific columns can be added here
    VehicleeType INT CHECK
    FOREIGN KEY (DriverId) REFERENCES Users(Id) ON DELETE SET NULL
);

