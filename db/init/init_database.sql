CREATE TABLE commands (
    id          serial PRIMARY KEY,
    name        text NOT null,
    raw         text NOT null,
    status      text,
    error_msg   text,
    created_at  timestamp NOT null DEFAULT current_timestamp,
    updated_at  timestamp NOT null DEFAULT current_timestamp
);

CREATE TABLE command_logs (
    command_id integer UNIQUE,
    logs       text,
    FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE SET null
)