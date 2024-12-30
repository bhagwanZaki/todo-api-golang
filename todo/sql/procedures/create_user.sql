CREATE OR REPLACE PROCEDURE create_user(
    IN username VARCHAR,
    IN fullname VARCHAR,
    IN email VARCHAR,
    IN password VARCHAR,
    IN created_at Date,
    OUT inserted_id INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.users (
        username,
        fullname,
        email,
        password,
        is_active,
        is_premium,
        is_admin,
        created_at,
        updated_at
    ) VALUES (
        username, 
        fullname, 
        email, 
        password,
        TRUE, 
        FALSE,
        FALSE,
        created_at, 
        created_at
    )
    RETURNING id INTO inserted_id;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;