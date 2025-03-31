ALTER TABLE "accounts" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT (now());
ALTER TABLE "entries" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT (now());
ALTER TABLE "transfers" ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT (now());
