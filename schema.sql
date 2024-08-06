CREATE TABLE cloudflare_image (
  id bigint NOT NULL PRIMARY KEY DEFAULT '{}'::jsonb,
  path jsonb DEFAULT '{}'::jsonb,
  author_name text
);