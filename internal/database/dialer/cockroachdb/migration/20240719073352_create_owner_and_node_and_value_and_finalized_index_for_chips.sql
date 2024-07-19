-- +goose Up
-- +goose StatementBegin
CREATE INDEX "idx_owner_node_value_finalized" ON "stake"."chips" ("owner", "node") STORING ("value", "finalized");

CREATE VIEW "stake"."stakings" AS
SELECT "owner" AS "staker", "node", count(*) AS "count", sum("value") AS "value"
FROM "stake"."chips"
WHERE "finalized" IS TRUE
GROUP BY "owner", "node";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "stake"."idx_owner_node_value_finalized";

DROP VIEW "stake"."stakings";
-- +goose StatementEnd
