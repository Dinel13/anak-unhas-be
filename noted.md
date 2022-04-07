cqlsh
CREATE KEYSPACE anakunhas WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use anakunhas;

gunakan bactch insert untuk memastikan konsistensi data

BEGIN BATCH
  INSERT INTO learn_cassandra.todo_by_user_email (user_email,creation_date,name) VALUES('alice@email.com', toTimestamp(now()), 'My first todo entry')

  INSERT INTO learn_cassandra.todos_shared_by_target_user_email (target_user_email, source_user_email,creation_date,name) VALUES('bob@email.com', 'alice@email.com',toTimestamp(now()), 'My first todo entry')

  INSERT INTO learn_cassandra.todos_shared_by_source_user_email (target_user_email, source_user_email,creation_date,name) VALUES('alice@email.com', 'bob@email.com', toTimestamp(now()), 'My first todo entry')

APPLY BATCH;


error di midleware token cek masih selalu server error