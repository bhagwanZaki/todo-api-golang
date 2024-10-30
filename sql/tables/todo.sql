CREATE SEQUENCE IF NOT EXISTS public.todo_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

ALTER SEQUENCE public.todo_id_seq
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.todo
(
    id integer NOT NULL DEFAULT nextval('todo_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default",
    completed BOOLEAN DEFAULT FALSE,
    user_id INTEGER,
    created_at DATE,
    updated_at DATE,

    CONSTRAINT todo_pkey PRIMARY KEY (id),
    CONSTRAINT todo_user_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER SEQUENCE public.todo_id_seq
    OWNED BY public.todo.id;
    
ALTER SEQUENCE public.todo_id_seq
    OWNER TO postgres;


ALTER TABLE IF EXISTS public.todo
    OWNER to postgres;


CREATE INDEX IF NOT EXISTS ix_todo_id
    ON public.todo USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;