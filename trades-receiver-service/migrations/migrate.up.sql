create table insider (
                         id uuid primary key default gen_random_uuid(),
                         cik int unique not null,
                         name varchar(255) not null
);

create table company (
                         id uuid primary key default gen_random_uuid(),
                         cik int not null,
                         name varchar(255) not null,
                         ticker varchar(255) not null,
                         unique (cik, ticker)
);

create table sec_filings (
                             id uuid primary key default gen_random_uuid(),
                             filing_type smallint,
                             url varchar(1024) not null,
                             insider_id uuid not null,
                             company_id uuid not null,
                             officer_position varchar(255),
                             reported_on date not null,
                             constraint fk_insider_id
                                 foreign key (insider_id)
                                     references insider (id)
                                     on delete set null,
                             constraint fk_company_id
                                 foreign key (company_id)
                                     references company (id)
                                     on delete set null
);

create table transaction (
                             id uuid primary key default gen_random_uuid(),
                             sec_filings_id uuid not null,
                             transaction_type_name varchar(125) not null,
                             average_price decimal not null,
                             total_shares bigint not null,
                             total_value decimal not null,
                             created_at timestamptz default current_timestamp,
                             constraint fk_sec_filings_id
                                 foreign key (sec_filings_id)
                                     references sec_filings (id)
                                     on delete cascade
);

create table security_transaction_holdings (
                                               id uuid primary key default gen_random_uuid(),
                                               transaction_id uuid,
                                               sec_filings_id uuid not null,
                                               quantity_owned_following_transaction decimal,
                                               security_title varchar(1024) not null,
                                               security_type smallint,
                                               quantity bigint not null,
                                               price_per_security decimal not null,
                                               transaction_date date not null,
                                               transaction_code smallint not null,
                                               created_at timestamptz default current_timestamp,
                                               constraint fk_transaction_id
                                                   foreign key (transaction_id)
                                                       references transaction (id)
                                                       on delete cascade,
                                               constraint fk_sec_filings_id
                                                   foreign key (sec_filings_id)
                                                       references sec_filings (id)
                                                       on delete cascade
);

