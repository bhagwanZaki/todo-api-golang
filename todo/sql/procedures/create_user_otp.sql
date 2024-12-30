CREATE OR REPLACE PROCEDURE create_user_otp(
    IN user_email VARCHAR,
    IN user_pin INTEGER,
    IN user_request INTEGER,
    IN created_at Date
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.user_otp (
        email,
        otp,
        request_type,
        created_at
    ) VALUES (
        user_email,
        user_pin,
        user_request,
        created_at
    );
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;