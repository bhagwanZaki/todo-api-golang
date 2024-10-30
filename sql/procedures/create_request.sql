CREATE OR REPLACE PROCEDURE create_request(
    IN email VARCHAR,
    IN request_type VARCHAR,
    IN created_at DATE
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.user_request (
        email,
        request_type,
        created_at
    ) VALUES (
        email,
        request_type,
        created_at
    );
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;