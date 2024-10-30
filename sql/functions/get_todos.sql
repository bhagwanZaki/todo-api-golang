CREATE OR REPLACE FUNCTION get_todos(userid INT)
  RETURNS TABLE(id INTEGER, name VARCHAR, completed boolean)
  LANGUAGE plpgsql AS
$func$
BEGIN
  RETURN QUERY
  SELECT 
	todo.id, 
	todo.name, 
	todo.completed 
FROM todo 
WHERE todo.user_id = userid
ORDER BY todo.id DESC;
END
$func$;