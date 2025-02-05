CREATE TABLE IF NOT EXISTS regexp_list
(
    id         BIGSERIAL PRIMARY KEY,
    qtype      INT           NOT NULL,
    name       VARCHAR(1000) NOT NULL,
    data       TEXT[]        NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL,
    updated_at TIMESTAMPTZ   NOT NULL,
    deleted_at TIMESTAMPTZ,
    UNIQUE (qtype, name)
);
