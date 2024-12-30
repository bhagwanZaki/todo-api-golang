-- CREATE SEQUENCE
CREATE SEQUENCE IF NOT EXISTS public.feedback_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.feedback
(
    id INTEGER NOT NULL DEFAULT nextval('feedback_id_seq'::regclass),
    feedback_user_id INTEGER,
    feedback_user_username character varying COLLATE pg_catalog."default",
    feedback_user_email character varying COLLATE pg_catalog."default",
    feedback character varying COLLATE pg_catalog."default",
    imageAddr character varying COLLATE pg_catalog."default",
    appName character varying COLLATE pg_catalog."default",
    created_at DATE,

    CONSTRAINT feedback_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.feedback_id_seq
    OWNED BY public.feedback.id;
    
ALTER SEQUENCE public.feedback_id_seq
    OWNER TO postgres;

-- ALTER TABLE
ALTER TABLE IF EXISTS public.feedback
    OWNER to postgres;

-- CREATE INDEX
CREATE INDEX IF NOT EXISTS ix_feedback_id
    ON public.feedback USING btree
    (id ASC NULLS LAST)
    TABLESPACE pg_default;