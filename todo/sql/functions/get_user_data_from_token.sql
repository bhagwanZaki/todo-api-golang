CREATE OR REPLACE FUNCTION get_user_data_from_token(
    user_token VARCHAR
)
  RETURNS TABLE(id INTEGER, username VARCHAR, email VARCHAR, fullname VARCHAR)
  LANGUAGE plpgsql AS
$func$
BEGIN
  RETURN QUERY
  SELECT users.id, users.username, users.email, users.fullname 
  FROM users 
  INNER JOIN token 
  ON token.user_id = users.id 
  WHERE token.token = user_token;
END
$func$;