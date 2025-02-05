INSERT INTO "black_list" ("url", "type", created_at, updated_at, deleted_at)
VALUES ('LOCAL', 'static', now(), now(), NULL),
       ('https://cdn.osspkg.com/adblock/ublock3.txt', 'dynamic', now(), now(), NULL),
       ('https://cdn.osspkg.com/adblock/ublock20.txt', 'dynamic', now(), now(), NULL),
       ('https://cdn.osspkg.com/adblock/adtidy2.txt', 'dynamic', now(), now(), NULL),
       ('https://cdn.osspkg.com/adblock/adtidy11.txt', 'dynamic', now(), now(), NULL),
       ('https://cdn.osspkg.com/adblock/easylist.txt', 'dynamic', now(), now(), NULL)
ON CONFLICT ("url") DO NOTHING;

-- # https://filters.adtidy.org/extension/ublock/filters/3.txt
-- # https://filters.adtidy.org/extension/ublock/filters/20.txt
-- # https://filters.adtidy.org/extension/ublock/filters/2_without_easylist.txt
-- # https://filters.adtidy.org/extension/ublock/filters/11.txt
-- # https://cdn.statically.io/gh/uBlockOrigin/uAssetsCDN/main/thirdparties/easylist.txt
