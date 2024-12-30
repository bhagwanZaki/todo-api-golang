CREATE OR REPLACE PROCEDURE create_token(
    IN user_id INTEGER,
    IN token VARCHAR,
    IN digest VARCHAR,
    IN created_at Date
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.token (
        user_id,
        token,
        digest,
        created_at
    ) VALUES (
        user_id, 
        token, 
        digest, 
        created_at
    );
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;