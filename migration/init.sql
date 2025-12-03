CREATE TABLE IF EXISTS "devices" (
    "id" serial PRIMARY KEY,
    "username" varchar(255) NOT NULL,
    "password" bytea NOT NULL,
    "color" text NOT NULL DEFAULT 'white',
    "timezone" text NOT NULL DEFAULT 'UTC',
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF EXISTS "alarms" (
    "id" serial PRIMARY KEY,
    "device_id" INTEGER NOT NULL,
    "time" Timestamp WITH TIME ZONE NOT NULL,
    "is_repeat" boolean NOT NULL DEFAULT false,
    "is_active" boolean NOT NULL DEFAULT true,
    "days" jsonb NOT NULL DEFAULT '[]',
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    FOREIGN KEY ("device_id") REFERENCES "devices" ("id")
);

INSERT INTO "devices" ("username", "password", "color", "timezone") VALUES ('admin', '\xc00001e240', 'white', 'UTC');
