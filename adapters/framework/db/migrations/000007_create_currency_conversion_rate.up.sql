UPDATE "currencies" SET "name" = 'Franc', "symbol" = 'CHF' WHERE "name" = 'Bitcoin';

ALTER TABLE "currencies" ADD COLUMN "conversion_rate_to_usd" NUMERIC(10,2) NOT NULL DEFAULT 1.0;

UPDATE "currencies" SET "conversion_rate_to_usd" = 0.9716 WHERE "name" = 'Franc';
UPDATE "currencies" SET "conversion_rate_to_usd" = 414.8065 WHERE "name" = 'Naira';
UPDATE "currencies" SET "conversion_rate_to_usd" = 1.0 WHERE "name" = 'Dollar';
UPDATE "currencies" SET "conversion_rate_to_usd" = 0.9480 WHERE "name" = 'Euro';
UPDATE "currencies" SET "conversion_rate_to_usd" = 0.7972 WHERE "name" = 'Pound';
UPDATE "currencies" SET "conversion_rate_to_usd" = 130.0931 WHERE "name" = 'Yen';
UPDATE "currencies" SET "conversion_rate_to_usd" = 6.6318 WHERE "name" = 'Yuan';
UPDATE "currencies" SET "conversion_rate_to_usd" = 76.4901 WHERE "name" = 'Rupee';