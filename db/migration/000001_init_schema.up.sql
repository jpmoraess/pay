CREATE TABLE "users" (
    "id" uuid PRIMARY KEY,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "role" varchar NOT NULL DEFAULT 'simple',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

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