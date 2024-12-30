CREATE OR REPLACE FUNCTION check_if_otp_already_exist(user_email VARCHAR, user_request_type INTEGER)
  RETURNS TABLE(id INTEGER, email VARCHAR, otp INTEGER)
  LANGUAGE plpgsql AS
$func$
BEGIN
  RETURN QUERY
  SELECT
    user_otp.id,
    user_otp.email,
    user_otp.otp
  FROM user_otp 
  WHERE 
    user_otp.email = user_email AND 
    user_otp.request_type = user_request_type
  LIMIT 1;
END
$func$;