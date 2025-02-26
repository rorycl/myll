/*
MyLensLocked

User management function
*/

CREATE OR REPLACE FUNCTION fn_user_manage (
    mode          e_mode
    ,in_name      TEXT
    ,in_email     TEXT
    ,in_hpassword TEXT -- hpassword is a bcrypt hashed password
    ,in_id        BIGINT
)
RETURNS users -- only returns one row
AS $$

DECLARE

    hereUser users%rowtype;

BEGIN

    IF mode IN ('create', 'update') THEN
    
        in_name = btrim(in_name);
        in_email = btrim(lower(in_email));

        IF (in_name IS NULL OR in_email IS NULL OR in_hpassword IS NULL) THEN 
            RAISE EXCEPTION 'all fields required in create mode';
        ELSIF length(in_name) < 5 THEN
            RAISE EXCEPTION 'names must be at least 5 characters in length';
        ELSIF length(in_email) < 5 THEN
            RAISE EXCEPTION 'emails must be at least 5 characters in length';
        ELSIF position('@' in in_email) < 2 THEN
            RAISE EXCEPTION 'emails require an @ symbol';
        ELSIF length(in_hpassword) != 60 THEN
            RAISE EXCEPTION 'hashed passwords are expected to be 60 characters in length';
        ELSIF position('$2a$' in in_hpassword) != 1 THEN
            RAISE EXCEPTION 'hashed password does not appear to be a bcrypt password';
        END IF;

    ELSIF mode NOT IN ('create') THEN
        IF in_id IS NULL THEN 
            RAISE EXCEPTION 'user id required';
        END IF;
    END IF;

    IF mode = 'create' THEN

        SELECT FROM
            users
        WHERE
            in_email = email
            OR
            lower(in_name) = lower(name)
        ;

        IF FOUND THEN
            RAISE EXCEPTION 'user with same name or email already exists';
        END IF;

        INSERT INTO users
            (name, email, hpassword)
        VALUES
            (in_name, in_email, in_hpassword)
        RETURNING
            *
        INTO
            hereUser;

    ELSIF mode = 'update' THEN

        SELECT
            *
        FROM
            users
        WHERE
            id = in_id
        INTO
            hereUser
        ;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'user % not found', in_id;
        END IF;

        IF hereUser.name <> in_name THEN
            RAISE EXCEPTION 'user name cannot be changed';
        END IF;

        UPDATE
            users
        SET
            name = in_name
            ,email = in_email
            ,hpassword = in_hpassword
            ,modified = current_timestamp
        WHERE
            id = in_id
        RETURNING
            *
        INTO
            hereUser;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'user % not found', in_id;
        END IF;

    ELSIF mode = 'delete' THEN

        DELETE FROM
            users
        WHERE
            id = in_id
        RETURNING
            *
        INTO
            hereUser;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'user % not found', in_id;
        END IF;

    ELSE -- view

        SELECT
            *
        FROM
            users
        WHERE
            id = in_id
        INTO
            hereUser;

    END IF;

    RETURN hereUser;

END;
$$
LANGUAGE plpgsql
;
