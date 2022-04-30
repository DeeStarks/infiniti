CREATE TABLE "beneficiaries" (
    "id" bigserial PRIMARY KEY,
    "user_id" int NOT NULL,
    "beneficiary_account_id" int NOT NULL
);

ALTER TABLE "beneficiaries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "beneficiaries" ADD FOREIGN KEY ("beneficiary_account_id") REFERENCES "user_accounts" ("id");