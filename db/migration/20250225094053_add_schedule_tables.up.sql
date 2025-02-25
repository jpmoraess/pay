CREATE TABLE services (
    "id" uuid PRIMARY KEY,
    "name" varchar NOT NULL,
    "duration_minutes" int NOT NULL,
    "price" decimal(10,2) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE schedules (
    "id" uuid PRIMARY KEY,
    "user_id" uuid REFERENCES users(id),
    "weekday" int NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "interval_minutes" int NOT NULL
);

CREATE TABLE schedule_overrides (
    "id" uuid PRIMARY KEY,
    "user_id" uuid REFERENCES users(id),
    "date" date NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "available" boolean NOT NULL
);

CREATE TABLE appointments (
    "id" uuid PRIMARY KEY,
    "service_id" uuid REFERENCES services(id),
    "customer_name" varchar NOT NULL,
    "appointment_time" timestamptz NOT NULL UNIQUE,
    "status" varchar NOT NULL CHECK (status IN ('confirmed', 'cancelled')),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);