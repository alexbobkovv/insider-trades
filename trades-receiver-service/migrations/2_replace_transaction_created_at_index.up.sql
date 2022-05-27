DROP INDEX IF EXISTS idx_transaction_pagination;
CREATE INDEX idx_transaction_pagination ON transaction USING brin (created_at);
