DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'principal_type') THEN
CREATE TYPE principal_type AS ENUM (
            'USER',
            'ROLE',
            'GROUP',
            'ACCOUNT',
            'SERVICE',
            'ANONYMOUS'
        );
END IF;
END $$;

CREATE TABLE public.statement_principals
(
    statement_id   BIGINT NOT NULL,
    principal_type principal_type NOT NULL,
    principal_id   VARCHAR(255) NOT NULL,

    CONSTRAINT statement_principals_pk
        PRIMARY KEY (
                     statement_id,
                     principal_type,
                     principal_id
            ),

    CONSTRAINT statement_principals_statement_fk
        FOREIGN KEY (statement_id)
            REFERENCES public.statement (id)
            ON DELETE CASCADE
);

CREATE INDEX statement_principals_lookup_idx
    ON public.statement_principals (
                                    principal_type,
                                    principal_id
        );