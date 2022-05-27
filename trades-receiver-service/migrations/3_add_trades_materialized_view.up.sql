CREATE MATERIALIZED VIEW trades_matview AS
SELECT transaction.id,
       transaction.sec_filings_id,
       transaction.transaction_type_name,
       transaction.average_price,
       transaction.total_shares,
       transaction.total_value,
       transaction.created_at,
       sf.url,
       sf.insider_id,
       sf.company_id,
       sf.officer_position,
       sf.reported_on,
       i.cik  AS insider_cik,
       i.name AS insider_name,
       c.cik  AS company_cik,
       c.name AS company_name,
       c.ticker
FROM transaction
         INNER JOIN sec_filings sf ON sf.id = transaction.sec_filings_id
         INNER JOIN insider i ON i.id = sf.insider_id
         INNER JOIN company c ON c.id = sf.company_id
ORDER BY transaction.created_at DESC;