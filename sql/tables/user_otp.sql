CREATE SEQUENCE IF NOT EXISTS public.user_otp_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

ALTER SEQUENCE public.user_otp_id_seq
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.user_otp
(
    id integer NOT NULL DEFAULT nextval('user_otp_id_seq'::regclass),
    email character varying COLLATE pg_catalog."default",
    otp INTEGER,
    request_type INTEGER,
    created_at DATE,

    CONSTRAINT user_otp_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.user_otp_id_seq
    OWNED BY public.user_otp.id;
    
ALTER SEQUENCE public.user_otp_id_seq
    OWNER TO postgres;


ALTER TABLE IF EXISTS public.user_otp
    OWNER to postgres;


CREATE INDEX IF NOT EXISTS ix_user_otp_id
    ON public.user_otp USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;