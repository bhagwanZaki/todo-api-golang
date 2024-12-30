CREATE OR REPLACE FUNCTION get_feedback()
RETURNS TABLE (
    id INTEGER,
    feedback_user_id INTEGER,
    feedback_user_username character varying,
    feedback_user_email character varying,
    feedback character varying,
    imageAddr character varying,
    appName character varying,
    created_at DATE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        id,
        feedback_user_id,
        feedback_user_username,
        feedback_user_email,
        feedback,
        imageAddr,
        created_at
    FROM public.feedback;
END;
$$ LANGUAGE plpgsql;
