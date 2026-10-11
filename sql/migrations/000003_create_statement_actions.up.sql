CREATE TABLE public.statement_actions
(
    statement_id BIGINT NOT NULL,
    action       VARCHAR(255) NOT NULL,

    CONSTRAINT statement_actions_pk
        PRIMARY KEY (
                     statement_id,
                     action
            ),

    CONSTRAINT statement_actions_statement_fk
        FOREIGN KEY (statement_id)
            REFERENCES public.statement (id)
            ON DELETE CASCADE
);

CREATE INDEX statement_actions_action_idx
    ON public.statement_actions (action);