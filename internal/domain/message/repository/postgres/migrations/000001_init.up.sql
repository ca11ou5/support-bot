BEGIN;

CREATE TABLE IF NOT EXISTS support_employee(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(63) NOT NULL
);

INSERT INTO support_employee (username, password) VALUES ('trofimovaa2', 'qwerty');
INSERT INTO support_employee (username, password) VALUES ('callous', '123456');

COMMIT;