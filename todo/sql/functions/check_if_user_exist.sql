CREATE OR REPLACE FUNCTION check_if_user_exist(
    user_email VARCHAR , user_username VARCHAR
)
RETURNS INT
LANGUAGE plpgsql AS
$func$
DECLARE user_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO user_count
    FROM users
    WHERE users.email = user_email OR users.username = user_username;
    RETURN user_count;
END
$func$;