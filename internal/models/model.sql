-- ========================================================
-- 🗄️ Table: upload_files
-- Stores metadata for each uploaded file.
-- ========================================================
CREATE TABLE upload_files (
  id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  file_name         TEXT        NOT NULL,
  original_file_name TEXT       NOT NULL,
  file_size_bytes   BIGINT      NOT NULL CHECK (file_size_bytes <= 10 * 1024 * 1024),
  mime_type         TEXT        NOT NULL,
  uploaded_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  version           TEXT        NOT NULL DEFAULT '1.0'
);

-- ========================================================
-- 🗄️ Table: threads
-- Each user “conversation” or session.
-- ========================================================
CREATE TABLE threads (
  id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  slug        TEXT        NOT NULL UNIQUE,
  user_id     UUID        NULL,                     -- future auth
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  version     TEXT        NOT NULL DEFAULT '1.0'
);

-- ========================================================
-- 🗄️ Table: messages
-- All prompts, LLM responses, widget events, etc.
-- ========================================================
CREATE TABLE messages (
  id             UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
  thread_id      UUID            NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
  query_text     TEXT            NULL,
  response_text  TEXT            NULL,
  event_type     TEXT            NULL CHECK (event_type IN (
                   'START','END','PLAN','WEB_RESULTS',
                   'MARKDOWN_ANSWER','RELATED_QUERIES','WIDGET'
                 )),
  model          TEXT            NOT NULL,
  input_token    INTEGER         NOT NULL,
  output_token   INTEGER         NOT NULL,
  response_time  DOUBLE PRECISION NOT NULL,         -- seconds
  stream_status  TEXT            NULL CHECK (stream_status IN ('IN_PROGRESS','DONE','FAILED')),
  message_index  INTEGER         NOT NULL,
  meta_data      JSONB           NOT NULL DEFAULT '{}'::jsonb,
  created_at     TIMESTAMPTZ     NOT NULL DEFAULT now(),
  version        TEXT            NOT NULL DEFAULT '1.0',
  UNIQUE(thread_id, message_index)
);

-- ========================================================
-- 🗄️ Table: message_files
-- Pivot for many-to-many messages ↔ upload_files
-- ========================================================
CREATE TABLE message_files (
  message_id UUID NOT NULL REFERENCES messages(id)  ON DELETE CASCADE,
  file_id    UUID NOT NULL REFERENCES upload_files(id) ON DELETE CASCADE,
  PRIMARY KEY (message_id, file_id)
);

-- ========================================================
-- 🔎 Indexes for performance
-- ========================================================
CREATE INDEX idx_messages_thread        ON messages(thread_id);
CREATE INDEX idx_messages_created       ON messages(created_at);
CREATE INDEX idx_threads_created        ON threads(created_at);
CREATE INDEX idx_upload_files_uploaded  ON upload_files(uploaded_at);

-- ========================================================
-- ⚙️ Triggers to bump `updated_at` on threads
-- ========================================================
CREATE OR REPLACE FUNCTION bump_thread_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_threads_update
BEFORE UPDATE ON threads
FOR EACH ROW EXECUTE PROCEDURE bump_thread_updated_at();
