CREATE TABLE messages (
    id uuid,
    from_user text,     
    to_user text,     
    body text,     
    read boolean,     
    time timeuuid,     
    PRIMARY KEY (id) 
)
CREATE INDEX IF NOT EXISTS messagetouser_key ON messages (to_user)

CREATE TABLE chat_read (
    id uuid,
    user_id text,     
    PRIMARY KEY (id) 
);