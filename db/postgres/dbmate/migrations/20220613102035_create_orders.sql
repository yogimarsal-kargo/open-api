-- migrate:up
SET LOCAL lock_timeout = '60s';
CREATE TABLE IF NOT EXISTS orders (
    ksuid varchar,
    client_id varchar,
    product_id varchar,
    num_sales bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT orders_pk PRIMARY KEY (ksuid)
);

-- migrate:down
SET LOCAL lock_timeout = '60s';
DROP TABLE IF EXISTS orders; 
