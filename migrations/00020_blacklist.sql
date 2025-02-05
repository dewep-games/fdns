CREATE TABLE IF NOT EXISTS black_list
(
    id         BIGSERIAL PRIMARY KEY,
    url        TEXT         NOT NULL UNIQUE,
    type       VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL,
    deleted_at TIMESTAMPTZ
);
