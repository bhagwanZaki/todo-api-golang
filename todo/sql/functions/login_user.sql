CREATE OR REPLACE FUNCTION login_user(
    dataUsername VARCHAR
)
  RETURNS TABLE(id INTEGER,username VARCHAR, email VARCHAR, password VARCHAR, fullname VARCHAR)
  LANGUAGE plpgsql AS
$func$
BEGIN
   RETURN QUERY
   SELECT users.id, users.username, users.email, users.password, users.fullname
   FROM users
   WHERE users.username=dataUsername;
END
$func$;