CREATE OR REPLACE FUNCTION get_feedback_excluding_appname(app_name_input character varying)
RETURNS TABLE (
    id INTEGER,
    feedback_user_id INTEGER,
    feedback_user_username character varying,
    feedback_user_email character varying,
    feedback character varying,
    imageAddr character varying,
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
    FROM public.feedback
    WHERE appName = app_name_input;
END;
$$ LANGUAGE plpgsql;
