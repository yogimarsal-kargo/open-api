-- migrate:up transaction:false
CREATE INDEX CONCURRENTLY IF NOT EXISTS order_details__order_ksuid_idx
    ON order_details USING btree
    (order_ksuid ASC NULLS LAST)
;

-- migrate:down
DROP INDEX IF EXISTS order_details__order_ksuid_idx;
