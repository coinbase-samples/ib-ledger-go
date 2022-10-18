/**
 * Copyright 2022 Coinbase Global, Inc.
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

CREATE TYPE direction as ENUM ('DEBIT', 'CREDIT');
CREATE TYPE ttype as ENUM ('ORDER', 'TRANSFER', 'CONVERT');

CREATE TABLE IF NOT EXISTS account
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    portfolio_id UUID                           NOT NULL,
    user_id      UUID                           NOT NULL,
    currency     TEXT                           NOT NULL,
    created_at   TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS account_balance
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID                           NOT NULL REFERENCES account (id),
    balance    NUMERIC          DEFAULT 0     NOT NULL CHECK (balance >= 0),
    hold       NUMERIC          DEFAULT 0     NOT NULL CHECK (hold >= 0),
    available  NUMERIC          DEFAULT 0     NOT NULL CHECK (available >= 0),
    created_at TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL,
    request_id UUID,
    count      NUMERIC          DEFAULT 1     NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction
(
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id        UUID  NOT NULL REFERENCES account (id),
    receiver_id      UUID  NOT NULL REFERENCES account (id),
    request_id       UUID  NOT NULL,
    transaction_type TTYPE NOT NULL,
    created_at       TIMESTAMPTZ(3)   DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS finalized_transaction
(
    transaction_id UUID PRIMARY KEY REFERENCES transaction (id),
    completed_at   TIMESTAMPTZ(3),
    canceled_at    TIMESTAMPTZ(3),
    failed_at      TIMESTAMPTZ(3),
    request_id     UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS entry
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id     UUID                           NOT NULL REFERENCES account (id),
    transaction_id UUID                           NOT NULL REFERENCES transaction (id),
    request_id     UUID                           NOT NULL,
    amount         NUMERIC                        NOT NULL CHECK (amount > 0),
    direction      DIRECTION                      NOT NULL,
    created_at     TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS hold
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id     UUID                           NOT NULL REFERENCES account (id),
    transaction_id UUID                           NOT NULL REFERENCES transaction (id),
    request_id     UUID                           NOT NULL,
    amount         NUMERIC                        NOT NULL CHECK (amount > 0),
    created_at     TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS released_hold
(
    hold_id     UUID PRIMARY KEY REFERENCES hold (id),
    released_at TIMESTAMPTZ(3) DEFAULT NOW() NOT NULL,
    request_id  UUID                         NOT NULL
);
