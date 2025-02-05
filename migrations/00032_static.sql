INSERT INTO regexp_list (qtype, name, data, created_at, updated_at, deleted_at)
VALUES (1, '^(\\d+)-(\\d+)-(\\d+)-(\\d+)\\.local\\.$','{"$1.$2.$3.$4"}', now(), now(), null)
ON CONFLICT (qtype, name) DO NOTHING;
