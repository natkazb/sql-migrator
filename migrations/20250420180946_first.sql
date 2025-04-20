-- SQL
-- Up begin
CREATE TABLE IF NOT EXISTS temp (
    id SERIAL PRIMARY KEY,
    start_date TIMESTAMP NOT NULL
);

INSERT INTO temp (start_date) VALUES (NOW());
-- Up end

-- Down begin
DROP TABLE temp;
-- Down end
