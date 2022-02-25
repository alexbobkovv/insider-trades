CREATE TABLE insider
(
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    cik  int          NOT NULL UNIQUE,
    name varchar(255) NOT NULL
);

CREATE TABLE company
(
    id     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    cik    int          NOT NULL UNIQUE,
    name   varchar(255) NOT NULL,
    ticker varchar(255) NOT NULL UNIQUE
);

CREATE TABLE sec_filings
(
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    filing_type      smallint,
    url              varchar(1024) UNIQUE NOT NULL,
    insider_id       uuid,
    company_id       uuid,
    officer_position varchar(255),
    reported_on      date,
    CONSTRAINT fk_insider_id
        FOREIGN KEY (insider_id)
            REFERENCES insider (id)
            ON DELETE SET NULL,
    CONSTRAINT fk_company_id
        FOREIGN KEY (company_id)
            REFERENCES company (id)
            ON DELETE SET NULL,
    UNIQUE(url, insider_id, company_id)
);

CREATE TABLE transaction
(
    id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    sec_filings_id        uuid    NOT NULL,
    transaction_type_name varchar(125),
    average_price         decimal NOT NULL,
    total_shares          bigint  NOT NULL,
    total_value           decimal,
    created_at            timestamptz      DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sec_filings_id
        FOREIGN KEY (sec_filings_id)
            REFERENCES sec_filings (id)
            ON DELETE CASCADE
);

CREATE TABLE security_transaction_holdings
(
    id                                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id                       uuid,
    sec_filings_id                       uuid,
    quantity_owned_following_transaction decimal,
    security_title                       varchar(1024),
    security_type                        smallint,
    quantity                             bigint,
    price_per_security                   decimal,
    transaction_date                     date,
    transaction_code                     smallint,
    created_at                           timestamptz      DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction_id
        FOREIGN KEY (transaction_id)
            REFERENCES transaction (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_sec_filings_id
        FOREIGN KEY (sec_filings_id)
            REFERENCES sec_filings (id)
            ON DELETE CASCADE
);
