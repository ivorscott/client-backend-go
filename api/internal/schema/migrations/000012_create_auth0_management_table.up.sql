CREATE TABLE auth0_management (
    auth0_management_id UUID PRIMARY KEY,
    token text NOT NULL,
    expiration text NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc')
);