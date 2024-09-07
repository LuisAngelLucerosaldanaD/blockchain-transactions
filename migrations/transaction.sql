CREATE TABLE bc.transaction
(
    id         uuid      NOT NULL,
    from_id    uuid      NOT NULL,
    to_id      uuid      NOT NULL,
    amount     float8    NOT NULL,
    type_id    int4      NOT NULL,
    "data"     varchar   NOT NULL,
    block      int8      NOT NULL,
    "files"    varchar   NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT transactions_pkey PRIMARY KEY (id),
    constraint FK_from_wallet foreign key (from_id) references auth.wallet (id),
    constraint FK_to_wallet foreign key (to_id) references auth.wallet (id)
);
