CREATE SEQUENCE IF NOT EXISTS public.token_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.token
(
    id INTEGER NOT NULL DEFAULT nextval('token_id_seq'::regclass),
    user_id INTEGER,
    token character varying COLLATE pg_catalog."default",
    digest character varying COLLATE pg_catalog."default",
    created_at date,

    CONSTRAINT token_pkey PRIMARY KEY (id),
    CONSTRAINT token_user_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER SEQUENCE public.token_id_seq
    OWNED BY public.token.id;
    
ALTER SEQUENCE public.token_id_seq
    OWNER TO postgres;


ALTER TABLE IF EXISTS public.token
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS ix_token_id
    ON public.token USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;