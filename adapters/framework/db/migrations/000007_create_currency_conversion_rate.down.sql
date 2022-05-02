ALTER TABLE "currencies" DROP COLUMN "conversion_rate_to_usd";

UPDATE "currencies" SET "name" = 'Bitcoin', "symbol" = 'BTC' WHERE "name" = 'Franc';