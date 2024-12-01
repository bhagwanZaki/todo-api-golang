CREATE OR REPLACE FUNCTION public.verify_user_otp(
	user_email character varying,
	register_otp integer,
	user_request_type integer
  )
    RETURNS TABLE(id integer, email character varying, otp integer) 
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
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
  SELECT * FROM deleted_row;

  -- Raise an exception if no rows were deleted
  IF NOT FOUND THEN
      RAISE EXCEPTION 'Invalid id';
  END IF;
END
$BODY$;