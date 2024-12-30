CREATE OR REPLACE FUNCTION check_user_exists(
    p_username VARCHAR,
    p_email VARCHAR
) RETURNS INT AS $$
DECLARE
    username_exists BOOLEAN;
    email_exists BOOLEAN;
BEGIN
    -- Check if the username exists
    SELECT EXISTS(SELECT 1 FROM users WHERE username = p_username) INTO username_exists;

    -- Check if the email exists
    SELECT EXISTS(SELECT 1 FROM users WHERE email = p_email) INTO email_exists;

    -- Determine the return value based on the results
    IF username_exists AND email_exists THEN
        RETURN 3; -- Both found
    ELSIF username_exists THEN
        RETURN 1; -- Username found
    ELSIF email_exists THEN
        RETURN 2; -- Email found
    ELSE
        RETURN 0; -- Neither found
    END IF;
END;
$$ LANGUAGE plpgsql;
