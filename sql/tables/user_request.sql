CREATE SEQUENCE IF NOT EXISTS public.user_request_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

ALTER SEQUENCE public.user_request_id_seq
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.user_request
(
    id integer NOT NULL DEFAULT nextval('user_request_id_seq'::regclass),
    email character varying COLLATE pg_catalog."default",
    request_type INTEGER,
    created_at DATE,

    CONSTRAINT user_request_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.user_request_id_seq
    OWNED BY public.user_request.id;
    
ALTER SEQUENCE public.user_request_id_seq
    OWNER TO postgres;


ALTER TABLE IF EXISTS public.user_request
    OWNER to postgres;


CREATE INDEX IF NOT EXISTS ix_user_request_id
    ON public.user_request USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;