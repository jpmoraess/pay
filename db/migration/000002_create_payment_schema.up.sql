CREATE TABLE "payments" (
    "id" uuid PRIMARY KEY,
    "external_id" varchar NOT NULL,
    "value" decimal(10,2) NOT NULL,
    "due_date" date NOT NULL,
    "status" varchar NOT NULL DEFAULT 'PENDING',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "payments" ("external_id");