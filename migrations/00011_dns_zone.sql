INSERT INTO "dns_zone" ("name", "data", "created_at", "updated_at")
VALUES ('.', '{"1.1.1.1", "1.0.0.1", "8.8.8.8", "8.8.4.4"}', now(), now()),
       ('ru.', '{"77.88.8.8", "77.88.8.1"}', now(), now()),
       ('su.', '{"77.88.8.8", "77.88.8.1"}', now(), now()),
       ('xn--p1ai.', '{"77.88.8.8", "77.88.8.1"}', now(), now())
ON CONFLICT ("name") DO NOTHING;