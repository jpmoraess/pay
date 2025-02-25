-- tenants
CREATE TABLE "tenants" (
    "id" uuid PRIMARY KEY,
    "name" varchar NOT NULL,
    "email" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL
);

CREATE INDEX ON "tenants" ("email");

-- users
CREATE TABLE "users" (
    "id" uuid PRIMARY KEY,
    "tenant_id" uuid NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "role" varchar(50) CHECK (role IN ('admin', 'service_provider')) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE INDEX ON "users" ("email");
CREATE INDEX ON "users" ("tenant_id");

-- sessions
CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "email" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY ("email") REFERENCES "users" ("email")
);

CREATE INDEX ON "sessions" ("email");

-- services
CREATE TABLE "services" (
    "id" uuid PRIMARY KEY,
    "tenant_id" uuid NOT NULL,
    "name" varchar NOT NULL,
    "price" decimal(10,2) NOT NULL,
    "description" varchar,
    FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id")
);

CREATE INDEX ON "services" ("tenant_id");

-- user_services
CREATE TABLE "user_services" (
    "id" uuid PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "service_id" uuid NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

-- schedules
CREATE TABLE "schedules" (
    "id" uuid PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "weekday" int NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "interval_minutes" int NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- schedule_overrides
CREATE TABLE "schedule_overrides" (
    "id" uuid PRIMARY KEY,
    "schedule_id" uuid NOT NULL,
    "date" date NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "reason" varchar,
    FOREIGN KEY ("schedule_id") REFERENCES "schedules" ("id")
);

-- appointments
CREATE TABLE "appointments" (
    "id" uuid PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "client_name" varchar NOT NULL,
    "service_id" uuid NOT NULL,
    "appointment_time" timestamptz NOT NULL,
    "status" varchar(50) CHECK (status IN ('pending', 'confirmed', 'canceled')) NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);