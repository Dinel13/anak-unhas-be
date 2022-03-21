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
CREATE INDEX IF NOT EXISTS messageread_key ON messages (read) ;


CREATE TABLE messages (     
    id uuid,     
    from_user int,          
    to_user int,          
    body text,          
    time timeuuid,
    PRIMARY KEY ((to_user), id) )  ;

CREATE TABLE readed (
    id uuid,     
    from_user int,          
    to_user int,          
    body text,          
    time timeuuid,
    PRIMARY KEY ((to_user), id) )  ; 
);

CREATE TABLE anakunhas.message (     
    from_user int,
    to_user int,
    read boolean,         
    time timestamp,
    body text,
    PRIMARY KEY ((to_user, read), from_user, time) 
) WITH CLUSTERING ORDER BY (from_user DESC, time ASC);

CREATE TABLE anakunhas.friend (     
    user int,
    friend int,
    time timestamp,
    last_message text,
    PRIMARY KEY ((user), friend) 
) WITH CLUSTERING ORDER BY (friend DESC);