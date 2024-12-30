CREATE OR REPLACE PROCEDURE update_email(data_id INTEGER, user_email VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_updated INT;
BEGIN
	UPDATE users 
    SET email = user_email 
    WHERE users.id = data_id;
	
	GET DIAGNOSTICS rows_updated = ROW_COUNT;
	IF rows_updated = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Duplicate data error';
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;