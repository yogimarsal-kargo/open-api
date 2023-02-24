-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."audit_resources" (
    "ksuid" varchar PRIMARY KEY,
    "actor_ksuid" varchar,
    "actor_name" varchar,
    "resource" varchar,
    "resource_ksuid" varchar,
    "action" varchar,
    "change" varchar,
    "from" jsonb,
    "to" jsonb,
    "event_ksuid" varchar,
    "created_at" timestamptz
    -- CONSTRAINT audit_resources__event_ksuid_fk FOREIGN KEY ("event_ksuid") REFERENCES "audit_events"("ksuid")
);


-- migrate:down
DROP TABLE IF EXISTS "public"."audit_resources";