CREATE OR REPLACE PROCEDURE logout_from_all_device(in_user_id INTEGER)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_deleted INT;
BEGIN
    DELETE FROM token WHERE token.user_id = in_user_id;

    GET DIAGNOSTICS rows_deleted = ROW_COUNT;

    IF rows_deleted = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;
