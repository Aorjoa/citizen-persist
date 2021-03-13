CREATE SEQUENCE IF NOT EXISTS citizens_id_seq;

CREATE TABLE "public"."citizens" (
    "id" int8 NOT NULL DEFAULT nextval('citizens_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "cid" text,
    PRIMARY KEY ("id")
);