CREATE TABLE city
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR UNIQUE NOT NULL,
    latitude   NUMERIC        NOT NULL,
    longitude  NUMERIC        NOT NULL,
    created_at TIMESTAMP      NOT NULL DEFAULT NOW()
);
CREATE INDEX cities_name_idx ON city (name);

CREATE TABLE temperature
(
    id         SERIAL PRIMARY KEY,
    city_id    INTEGER   NOT NULL REFERENCES city (id),
    min        INTEGER   NOT NULL,
    max        INTEGER   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX temperatures_city_id_idx ON temperature (city_id);
CREATE INDEX temperatures_min_idx ON temperature (min);
CREATE INDEX temperatures_max_idx ON temperature (max);
CREATE INDEX temperatures_created_at_idx ON temperature (created_at);

CREATE TABLE webhook
(
    id           SERIAL PRIMARY KEY,
    city_id      INTEGER   NOT NULL REFERENCES city (id),
    callback_url VARCHAR   NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX webhooks_city_id_idx ON temperature (city_id);
