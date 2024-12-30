CREATE OR REPLACE PROCEDURE insert_todo(
	IN todo_user_id integer,
	IN todo_name text,
	IN status boolean,
	IN todo_created_at DATE,
	OUT inserted_id integer
)
LANGUAGE plpgsql
AS $$
BEGIN
	INSERT INTO public.todo (name, completed, user_id, created_at, updated_at)
    VALUES (todo_name, status, todo_user_id, todo_created_at, todo_created_at)
    RETURNING id INTO inserted_id;
EXCEPTION
	WHEN OTHERS THEN
		RAISE EXCEPTION '%', SQLERRM;
END;
$$;