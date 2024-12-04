CREATE OR REPLACE PROCEDURE create_feedback(
    IN inUserID INTEGER,
    IN inUserName VARCHAR,
    IN inUserEmail VARCHAR,
    IN inFeedback VARCHAR,
    IN inImageAddr VARCHAR,
    IN inAppName VARCHAR,
    IN inCreatedAt DATE
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.feedback (
        feedback_user_id,
        feedback_user_username,
        feedback_user_email,
        feedback,
        imageAddr,
        appName,
        created_at
    ) VALUES (
        inUserID, 
        inUserName, 
        inUserEmail, 
        inFeedback,
        inImageAddr, 
        inAppName,
        inCreatedAt
    );
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;