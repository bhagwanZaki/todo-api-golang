CREATE OR REPLACE PROCEDURE logout_from_all_device_except_one(in_user_id INTEGER, current_token VARCHAR)
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM token WHERE token.user_id = in_user_id AND token.token != current_token;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;