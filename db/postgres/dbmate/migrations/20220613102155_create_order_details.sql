-- migrate:up
SET LOCAL lock_timeout = '60s';
CREATE TABLE IF NOT EXISTS order_details (
    ksuid varchar,
    order_ksuid varchar,
    order_type varchar,
    origin varchar,
    destination varchar,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT order_details_pk PRIMARY KEY (ksuid),
    CONSTRAINT order_details__order_ksuid_fk FOREIGN KEY (order_ksuid) REFERENCES orders(ksuid)
);

-- migrate:down
SET LOCAL lock_timeout = '60s';
DROP TABLE IF EXISTS order_details;
