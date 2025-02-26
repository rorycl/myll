/*
My Lens Locked schema file

Started 23 February 2025
*/

-- CREATE SCHEMA IF NOT EXISTS myll;

SET search_path = myll;

-- mode type
CREATE TYPE e_mode AS ENUM (
    'create'
    ,'update'
    ,'delete'
    ,'view'
);

/*
users are the users of the service
*/
CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY
    ,created TIMESTAMPTZ NOT NULL DEFAULT current_timestamp --  timezone with timezone
    ,modified TIMESTAMPTZ NOT NULL DEFAULT current_timestamp --  timezone with timezone
    ,name TEXT NOT NULL
    ,email TEXT NOT NULL
    ,hpassword TEXT NOT NULL
);

