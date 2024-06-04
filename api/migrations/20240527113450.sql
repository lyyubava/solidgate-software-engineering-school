-- Create "rates" table
CREATE TABLE "public"."rates" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "rate" numeric NULL,
  "exchange_date" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_rates_deleted_at" to table: "rates"
CREATE INDEX "idx_rates_deleted_at" ON "public"."rates" ("deleted_at");
