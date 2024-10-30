CREATE OR REPLACE PROCEDURE update_password(user_email VARCHAR, user_password VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_updated INT;
BEGIN
	UPDATE users 
    SET password = user_password 
    WHERE users.email = user_email;
	
	GET DIAGNOSTICS rows_updated = ROW_COUNT;
	IF rows_updated = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;
