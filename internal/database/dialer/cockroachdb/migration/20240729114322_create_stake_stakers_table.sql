-- +goose Up
-- +goose StatementBegin
CREATE VIEW stake.stakers AS
SELECT staker AS address, count(node) AS nodes, sum(count) AS chip_number, sum(value) AS chip_value
FROM stake.stakings
GROUP BY staker;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW stake.stakers;
-- +goose StatementEnd
