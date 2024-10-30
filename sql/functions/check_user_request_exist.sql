CREATE OR REPLACE FUNCTION check_user_request_exist(
    user_email VARCHAR,
    user_request INTEGER,
    user_created_at DATE
)
  RETURNS INT
  LANGUAGE plpgsql AS
$func$
DECLARE request_exist INTEGER;
BEGIN
  SELECT COUNT(*)
  INTO request_exist
  FROM user_request
  WHERE 
    user_request.email = user_email AND
    user_request.request_type = user_request AND
    user_request.created_at = user_created_at;
  RETURN request_exist;
END
$func$;