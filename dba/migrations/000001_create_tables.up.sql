/**
* Copyright 2022-present Coinbase Global, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*  http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

CREATE EXTENSION "uuid-ossp";
CREATE TYPE DIRECTION AS ENUM ('DEBIT', 'CREDIT');
CREATE TYPE TRANSACTION_TYPE AS ENUM ('ORDER', 'TRANSFER', 'CONVERT');
CREATE TYPE TRANSACTION_STATUS as ENUM ('PENDING', 'FILLED', 'FAILED', 'CANCELED');

CREATE TABLE IF NOT EXISTS account (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    qldb_id VARCHAR(64) NOT NULL,
    user_id UUID NOT NULL,
    currency VARCHAR(64) NOT NULL
);

CREATE INDEX user_accounts ON account USING HASH(user_id);
CREATE INDEX idem_user_account ON account (currency, user_id);

CREATE TABLE IF NOT EXISTS account_balance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES account (id),
    balance NUMERIC DEFAULT 0 NOT NULL CHECK (balance >= 0),
    hold NUMERIC DEFAULT 0 NOT NULL CHECK (hold >= 0),
    available NUMERIC DEFAULT 0 NOT NULL CHECK (available >= 0),
    created_at TIMESTAMPTZ(3) DEFAULT NOW() NOT NULL,
    request_id UUID,
    idem VARCHAR(32),
    CONSTRAINT balance_change UNIQUE (account_id, idem)
);

CREATE INDEX account_balance_index ON account_balance USING HASH(account_id);

CREATE TABLE IF NOT EXISTS transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    qldb_id VARCHAR(64) NOT NULL,
    sender_id UUID NOT NULL REFERENCES account (id),
    receiver_id UUID NOT NULL REFERENCES account (id),
    transaction_type TRANSACTION_TYPE NOT NULL,
    created_at TIMESTAMPTZ(3) DEFAULT NOW(),
    finalized_at TIMESTAMPTZ(3),
    transaction_status TRANSACTION_STATUS
);

CREATE INDEX sender_transactions ON transaction USING HASH(sender_id);
CREATE INDEX receiver_transactions ON transaction USING HASH(receiver_id);

CREATE TABLE IF NOT EXISTS entry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES account (id),
    venue_order_id UUID NOT NULL REFERENCES transaction (id),
    fill_id UUID NOT NULL,
    amount NUMERIC NOT NULL CHECK (amount > 0),
    direction DIRECTION NOT NULL,
    created_at TIMESTAMPTZ(3) DEFAULT NOW() NOT NULL
);

CREATE INDEX transaction_entries ON entry USING HASH(venue_order_id);
CREATE INDEX fill_entries ON entry USING HASH(fill_id);
CREATE INDEX account_entries ON entry USING HASH(account_id);

CREATE TABLE IF NOT EXISTS hold (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES account (id),
    venue_order_id UUID NOT NULL REFERENCES transaction (id),
    amount NUMERIC NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMPTZ(3) DEFAULT NOW() NOT NULL, 
    released_at TIMESTAMPTZ(3),
    released BOOLEAN
);

