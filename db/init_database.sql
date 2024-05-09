CREATE TABLE commands (
    id          serial PRIMARY KEY,
    name        text NOT null,
    raw         text NOT null,
    status      text,
    error_msg   text,
    created_at  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE command_logs (
    command_id integer,
    logs       text,
    FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE SET NULL
)