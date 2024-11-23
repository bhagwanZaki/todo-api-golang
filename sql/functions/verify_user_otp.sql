CREATE OR REPLACE FUNCTION verify_user_otp(
  user_email VARCHAR, 
  register_otp INTEGER, 
  user_request_type INTEGER
  )
  RETURNS TABLE(id INTEGER, email VARCHAR, otp INTEGER)
  LANGUAGE plpgsql AS
$func$
DECLARE
    rows_deleted INT;
BEGIN
  RETURN QUERY
  WITH deleted_row AS (
    DELETE FROM user_otp
    WHERE 
      user_otp.email = user_email 
      AND user_otp.otp = register_otp 
      AND user_otp.request_type = user_request_type
    RETURNING user_otp.id, user_otp.email, user_otp.otp
  )
  
  GET DIAGNOSTICS rows_deleted = ROW_COUNT;

  IF rows_deleted = 0 THEN
      RAISE EXCEPTION 'Invalid id';
  END IF;

  SELECT * FROM deleted_row
  LIMIT 1;
END
$func$;