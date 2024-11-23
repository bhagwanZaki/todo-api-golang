CREATE OR REPLACE PROCEDURE update_fullname(data_id INTEGER, user_fullname VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_updated INT;
BEGIN
	UPDATE users 
    SET fullname = user_fullname
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