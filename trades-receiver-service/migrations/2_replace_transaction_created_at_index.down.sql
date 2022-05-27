DROP INDEX IF EXISTS idx_transaction_pagination;
CREATE INDEX idx_transaction_pagination ON transaction (created_at);
