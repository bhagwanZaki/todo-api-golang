CREATE OR REPLACE PROCEDURE update_todo(todo_id INTEGER, todo_name VARCHAR, todo_status Boolean, todo_user_id INTEGER)
LANGUAGE plpgsql
AS $$
DECLARE
    rows_deleted INT;
BEGIN
	UPDATE todo 
    SET name = todo_name, completed = todo_status 
    WHERE todo.id = todo_id AND todo.user_id = todo_user_id;
	
	GET DIAGNOSTICS rows_deleted = ROW_COUNT;
	IF rows_deleted = 0 THEN
        RAISE EXCEPTION 'Invalid id';
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Transaction failed and rolled back: %', SQLERRM;
END;
$$;
