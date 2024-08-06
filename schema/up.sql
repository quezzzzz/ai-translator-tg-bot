CREATE TABLE IF NOT EXISTS translations (
    id  serial primary key,
    user_id integer,
    original_text text,
    translated_text text,
    timestamp date
);