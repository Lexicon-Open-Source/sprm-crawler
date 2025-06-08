-- POSTGRESQL

-- Create UrlFrontier Table
CREATE TABLE "url_frontiers" (
  "id" varchar(64) NOT NULL PRIMARY KEY,
  "domain" varchar(255) NOT NULL,
  "url" varchar(255) NOT NULL,
  "crawler" varchar(255) NOT NULL,
  "status" smallint NOT NULL DEFAULT 0,
  "metadata" jsonb NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE url_frontiers ADD CONSTRAINT url_frontiers_unique UNIQUE (url);

-- Set Comment for Status Column
COMMENT ON COLUMN "url_frontiers"."status" IS '0: Pending, 1: Crawled, 2: Changed';

-- Create index for domain
CREATE INDEX idx_domain ON url_frontiers (domain);

-- Create index for url
CREATE INDEX idx_url ON url_frontiers (url);

-- Create index for crawler
CREATE INDEX idx_crawler ON url_frontiers (crawler);

-- Create extractions table
CREATE TABLE "extractions" (
  "id" varchar(64) NOT NULL PRIMARY KEY,
  "url_frontier_id" varchar(64) NOT NULL,
  "site_content" text NULL,
  "artifact_link" varchar(255) NULL,
  "raw_page_link" varchar(255) NULL,
  "metadata" jsonb NULL,
  "language" varchar(10) NOT NULL DEFAULT 'en',
  "page_hash" varchar(255) NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_extractions_url_frontiers_id FOREIGN KEY (url_frontier_id) REFERENCES url_frontiers(id)
)