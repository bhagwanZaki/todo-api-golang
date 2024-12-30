CREATE OR REPLACE FUNCTION get_user_detail(
    user_email_username VARCHAR
)
RETURNS TABLE(id INTEGER, username VARCHAR, email VARCHAR, fullname VARCHAR)
LANGUAGE plpgsql AS
$func$
BEGIN
    RETURN QUERY 
    SELECT 
    users.id,
    users.username,
    users.email 
    FROM users
    WHERE users.email = user_email_username OR users.username = user_email_username;
END
$func$;