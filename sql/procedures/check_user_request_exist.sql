CREATE OR REPLACE PROCEDURE check_user_request_exist(
    user_email VARCHAR,
    user_request_type INTEGER,
    user_created_at DATE,
    OUT deletion_successful BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_deleted INT;
BEGIN
    DELETE FROM user_request
    WHERE email = user_email
      AND request_type = user_request_type
      AND created_at = user_created_at;

    GET DIAGNOSTICS rows_deleted = ROW_COUNT;
    deletion_successful := (rows_deleted > 0);
END;
$$;
