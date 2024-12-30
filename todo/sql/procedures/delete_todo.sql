CREATE OR REPLACE PROCEDURE delete_todo(todo_id INTEGER, todo_user_id INTEGER)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_deleted INT;
BEGIN
    DELETE FROM todo WHERE todo.id = todo_id AND todo.user_id = todo_user_id;

    GET DIAGNOSTICS rows_deleted = ROW_COUNT;

    IF rows_deleted = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION '%', SQLERRM;
END;
$$;
