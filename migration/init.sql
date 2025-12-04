CREATE TABLE IF NOT EXISTS "devices" (
    "id" serial PRIMARY KEY,
    "username" varchar(255) NOT NULL,
    "password" bytea NOT NULL,
    "color" text NOT NULL DEFAULT 'white',
    "timezone" text NOT NULL DEFAULT 'UTC',
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "alarms" (
    "id" serial PRIMARY KEY,
    "device_id" INTEGER NOT NULL,
    "time" Timestamp WITH TIME ZONE NOT NULL,
    "is_repeat" boolean NOT NULL DEFAULT false,
    "is_active" boolean NOT NULL DEFAULT true,
    "days" jsonb NOT NULL DEFAULT '[]',
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    FOREIGN KEY ("device_id") REFERENCES "devices" ("id")
);

INSERT INTO "devices" ("username", "password", "color", "timezone") VALUES ('admin', '\x2432612431302441624866564a414239385776397258524674544c504f54696c6e6a75326932552e323868496f4f776931327261393738336e556c57', 'white', 'UTC');
