CREATE OR REPLACE PROCEDURE update_profile(
    data_id INTEGER, 
    user_username VARCHAR,
    user_email VARCHAR,
    user_fullname VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_updated INT;
BEGIN
    -- Attempt to update the user profile
    UPDATE users 
    SET email = user_email, 
        username = user_username, 
        fullname = user_fullname
    WHERE id = data_id;
    
    -- Check if any row was actually updated
    GET DIAGNOSTICS rows_updated = ROW_COUNT;
    IF rows_updated = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;

EXCEPTION
    -- Handle duplicate key violation for unique constraints
    WHEN unique_violation THEN
        IF SQLERRM LIKE '%users_email_key%' THEN
            RAISE EXCEPTION 'Duplicate email';
        ELSIF SQLERRM LIKE '%users_username_key%' THEN
            RAISE EXCEPTION 'Duplicate username';
        ELSIF SQLERRM LIKE '%users_email_key%' AND SQLERRM LIKE '%users_username_key%' THEN
            RAISE EXCEPTION 'Duplicate email and username';
        ELSE
            RAISE EXCEPTION 'Duplicate data error';
        END IF;
    
    -- Catch-all for other exceptions
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;
