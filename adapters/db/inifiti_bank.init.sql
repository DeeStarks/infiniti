CREATE TABLE "db_tables" (
  "id" bigserial PRIMARY KEY,
  "table_name" varchar NOT NULL
);

CREATE TABLE "permissions" (
  "id" bigserial PRIMARY KEY,
  "table_id" int NOT NULL,
  "method" varchar NOT NULL
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_permissions" (
  "id" bigserial PRIMARY KEY,
  "user_id" int,
  "permission_id" int
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "group_permissions" (
  "id" bigserial PRIMARY KEY,
  "group_id" int,
  "permission_id" int
);

CREATE TABLE "user_groups" (
  "group_id" int,
  "user_id" int
);

CREATE TABLE "currencies" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "symbol" varchar NOT NULL
);

CREATE TABLE "account_types" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "user_accounts" (
  "id" bigserial PRIMARY KEY,
  "user_id" int NOT NULL,
  "account_type_id" int,
  "account_number" bigint UNIQUE NOT NULL,
  "balance" double PRECISION DEFAULT 0.0,
  "currency_id" int
);

CREATE TABLE "transaction_types" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "user_transactions" (
  "id" bigserial PRIMARY KEY,
  "user_id" int,
  "transaction_type_id" int,
  "amount" double PRECISION NOT NULL,
  "sender_id" int,
  "reciever_id" int,
  "remark" varchar,
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "permissions" ADD FOREIGN KEY ("table_id") REFERENCES "db_tables" ("id");

ALTER TABLE "user_permissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "group_permissions" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "group_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "user_groups" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "user_groups" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_accounts" ADD FOREIGN KEY ("account_type_id") REFERENCES "account_types" ("id");

ALTER TABLE "user_accounts" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");

ALTER TABLE "user_transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_transactions" ADD FOREIGN KEY ("transaction_type_id") REFERENCES "transaction_types" ("id");

ALTER TABLE "user_transactions" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id");

ALTER TABLE "user_transactions" ADD FOREIGN KEY ("reciever_id") REFERENCES "users" ("id");

CREATE INDEX ON "permissions" ("table_id");

CREATE INDEX ON "permissions" ("method");

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("last_name");

CREATE INDEX ON "users" ("first_name", "last_name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "user_permissions" ("user_id");

CREATE INDEX ON "group_permissions" ("group_id");

CREATE INDEX ON "user_accounts" ("user_id");

CREATE INDEX ON "user_accounts" ("account_type_id");

CREATE INDEX ON "user_accounts" ("account_number");

CREATE INDEX ON "user_accounts" ("currency_id");

CREATE INDEX ON "user_transactions" ("user_id");

CREATE INDEX ON "user_transactions" ("transaction_type_id");

CREATE INDEX ON "user_transactions" ("sender_id");

CREATE INDEX ON "user_transactions" ("reciever_id");
