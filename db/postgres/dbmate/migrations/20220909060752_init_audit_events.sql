-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."audit_events" (
    "ksuid" varchar PRIMARY KEY,
    "actor_ksuid" varchar,
    "actor_name" varchar,
    "event_name" varchar,
    "resource_ksuid" varchar,
    "metadata" jsonb,
    "created_at" timestamptz
);

-- migrate:down

DROP TABLE IF EXISTS public.audit_events;

