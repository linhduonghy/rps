CREATE TABLE public.statement_resources
(
    statement_id BIGINT NOT NULL,
    resource     VARCHAR(500) NOT NULL,

    CONSTRAINT statement_resources_pk
        PRIMARY KEY (
                     statement_id,
                     resource
            ),

    CONSTRAINT statement_resources_statement_fk
        FOREIGN KEY (statement_id)
            REFERENCES public.statement (id)
            ON DELETE CASCADE
);

CREATE INDEX statement_resources_resource_idx
    ON public.statement_resources (resource);