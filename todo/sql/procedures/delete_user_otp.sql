CREATE OR REPLACE PROCEDURE delete_user_top(user_otp_id INTEGER, user_email VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_deleted INT;
BEGIN
    DELETE FROM user_otp WHERE user_otp.id = user_otp_id AND user_otp.email = user_email;

    GET DIAGNOSTICS rows_deleted = ROW_COUNT;

    IF rows_deleted = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;
