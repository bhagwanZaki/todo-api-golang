CREATE OR REPLACE FUNCTION check_auth(
    user_token VARCHAR
)
  RETURNS TABLE(id INTEGER, token VARCHAR, digest VARCHAR)
  LANGUAGE plpgsql AS
$func$
BEGIN
  RETURN QUERY
  SELECT token.id, token.token, token.digest
  FROM token
  WHERE token.token = user_token;
END
$func$;