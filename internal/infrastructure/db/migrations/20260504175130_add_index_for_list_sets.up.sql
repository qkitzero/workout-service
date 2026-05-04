CREATE INDEX IF NOT EXISTS idx_sets_user_id_trained_at_id ON sets (user_id, trained_at DESC, id DESC);
