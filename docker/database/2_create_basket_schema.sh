#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
    CREATE SCHEMA baskets;

    CREATE TABLE baskets.baskets
    (
        id      text NOT NULL,
        customer_id text NOT NULL,
        payment_id text NOT NULL,
        items bytea NOT NULL,
        status bytea NOT NULL,
        created_at timestamptz NOT NULL DEFAULT NOW(),
        updated_at timestamptz NOT NULL DEFAULT NOW(),
        PRIMARY KEY (id)
    );

    CREATE TRIGGER created_at_baskets_trgr BEFORE UPDATE ON baskets.baskets FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
    CREATE TRIGGER updated_at_baskets_trgr BEFORE UPDATE ON baskets.baskets FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

    CREATE TABLE baskets.events
    (
        stream_id       text    NOT NULL,
        stream_name     text    NOT NULL,
        stream_version  int     NOT NULL,
        event_id        text    NOT NULL,
        event_name      text    NOT NULL,
        event_data      bytea   NOT NULL,
        occurred_at     timestamptz    NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (stream_id, stream_name, stream_version)
    );

    CREATE TABLE baskets.snapshots
    (
        stream_id       text        NOT NULL,
        stream_name     text        NOT NULL,
        stream_version  int         NOT NULL,
        snapshot_name   text        NOT NULL,
        snapshot_data   bytea       NOT NULL,
        updated_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (stream_id, stream_name)
    );

    CREATE TRIGGER updated_at_snapshot_trgr BEFORE UPDATE ON baskets.snapshots FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

    GRANT USAGE ON SCHEMA baskets TO mallbots_user;
    GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA baskets TO mallbots_user;
EOSQL