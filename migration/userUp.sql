CREATE TABLE "users" (
  "id" bigserial primary key,
  "name" varchar NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL,
  "verified" boolean NOT NULL DEFAULT false,
  "image" varchar DEFAULT '',
  "wa" varchar DEFAULT NULL,
  "gender" varchar DEFAULT NULL,
  "address" varchar DEFAULT NULL,
  "jurusan" varchar DEFAULT NULL,
  "fakultas" varchar DEFAULT NULL,
  "ig" varchar DEFAULT NULL,
  "bio" varchar DEFAULT NULL,
  "tertarik" varchar DEFAULT NULL,
  "angkatan" VARCHAR DEFAULT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");
CREATE INDEX ON "users" ("email");
CREATE INDEX ON "users" ("name");