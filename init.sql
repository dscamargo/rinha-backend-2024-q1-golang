CREATE TABLE IF NOT EXISTS "clientes" (
                                          "id" serial PRIMARY KEY NOT NULL,
                                          "nome" text NOT NULL,
                                          "saldo" integer DEFAULT 0 NOT NULL,
                                          "limite" integer DEFAULT 0 NOT NULL
);

CREATE INDEX clientes_id_idx ON "clientes" USING HASH(id);

CREATE TABLE IF NOT EXISTS "transacoes" (
                                            "id" serial PRIMARY KEY NOT NULL,
                                            "cliente_id" integer NOT NULL ,
                                            "valor" integer NOT NULL,
                                            "tipo" char(1) NOT NULL,
                                            "descricao" varchar(10) NOT NULL,
                                            "created_at" timestamp DEFAULT now(),
                                            CONSTRAINT fk_cliente FOREIGN KEY(cliente_id) references clientes(id) on delete set null on update no action

);

CREATE INDEX transacoes_id_idx ON "transacoes" USING HASH(id);
CREATE INDEX transacoes_cliente_id_idx ON "transacoes" USING HASH(cliente_id);

create or replace procedure criar_transacao(
    in_cliente_id INTEGER,
    in_valor integer,
    in_tipo text,
    in_descricao text,
    inout in_saldo_atualizado integer default null,
    inout in_limite_atualizado integer default null
)

    language plpgsql
as $$

begin
    UPDATE clientes
    set saldo = saldo + in_valor
    where id = in_cliente_id and saldo + in_valor >= - limite
    returning saldo, limite into in_saldo_atualizado, in_limite_atualizado;

    if in_saldo_atualizado is null or in_limite_atualizado is null then return; end if;

    commit;

    INSERT INTO transacoes (valor, tipo, descricao, cliente_id)
    VALUES (ABS(in_valor), in_tipo, in_descricao, in_cliente_id);
end;
$$;

DO $$
    BEGIN
        INSERT INTO clientes (nome, limite)
        VALUES
            ('o barato sai caro', 1000 * 100),
            ('zan corp ltda', 800 * 100),
            ('les cruders', 10000 * 100),
            ('padaria joia de cocaia', 100000 * 100),
            ('kid mais', 5000 * 100);
    END; $$