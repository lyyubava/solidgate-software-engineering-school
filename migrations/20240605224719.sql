-- Create "emails" table
CREATE TABLE "public"."emails" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_emails_deleted_at" to table: "emails"
CREATE INDEX "idx_emails_deleted_at" ON "public"."emails" ("deleted_at");
