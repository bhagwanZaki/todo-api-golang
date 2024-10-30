-- CREATE SEQUENCE
CREATE SEQUENCE IF NOT EXISTS public.user_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

-- CREATE TABLE
CREATE TABLE IF NOT EXISTS public.users
(
    id INTEGER NOT NULL DEFAULT nextval('user_id_seq'::regclass),
    username character varying COLLATE pg_catalog."default",
    fullname character varying COLLATE pg_catalog."default",
    email character varying COLLATE pg_catalog."default",
    password character varying COLLATE pg_catalog."default",
    is_active BOOLEAN DEFAULT TRUE,
    is_premium BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at DATE,
    updated_at DATE,

    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email),
    CONSTRAINT user_username_key UNIQUE (username)
);

ALTER SEQUENCE public.user_id_seq
    OWNED BY public.users.id;
    
ALTER SEQUENCE public.user_id_seq
    OWNER TO postgres;

-- ALTER TABLE
ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

-- CREATE INDEX
CREATE INDEX IF NOT EXISTS ix_user_id
    ON public.users USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;