CREATE TABLE IF NOT EXISTS black_list_rules
(
    id         BIGSERIAL PRIMARY KEY,
    list_id    BIGINT       NOT NULL,
    data       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (list_id) REFERENCES black_list (id)
);