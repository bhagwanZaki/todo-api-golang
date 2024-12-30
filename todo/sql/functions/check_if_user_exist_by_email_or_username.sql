CREATE OR REPLACE FUNCTION check_if_user_exist_by_email_or_username(
    user_email_username VARCHAR
)
RETURNS INT
LANGUAGE plpgsql AS
$func$
DECLARE user_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO user_count
    FROM users
    WHERE users.email = user_email_username OR users.username = user_email_username;
    RETURN user_count;
END
$func$;