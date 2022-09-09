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

/*
 * Neoworks Fee Accounts
*/
INSERT INTO account(id, portfolio_id, user_id, currency, created_at)
VALUES ('B72D0E55-F53A-4DB0-897E-2CE4A73CB94B', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '116BDE43-7733-43A1-A85A-FC8627E6DA8E', 'USD', now());

/*
 * Coinbase Fee Accounts
*/
INSERT INTO account(id, portfolio_id, user_id, currency, created_at)
VALUES ('C4D0E14E-1B2B-4023-AFA6-8891AD1960C9', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '433C0C15-0A44-49C4-A207-4501BB11F48C', 'USD', now());

/*
 * Demo Accounts
 */

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('0adbb104-fc18-46ca-a4eb-beee7775eb69', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('0adbb104-fc18-46ca-a4eb-beee7775eb69', 10000000, 0, 10000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('0e4f9af4-4972-4e09-aee4-69b96db0c129', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'ETH', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('0e4f9af4-4972-4e09-aee4-69b96db0c129', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('084b5346-3214-419b-b7ec-e9c26a32e7c9', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'SOL', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('084b5346-3214-419b-b7ec-e9c26a32e7c9', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('55886255-6035-44ca-b531-a738670aa133', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'BTC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('55886255-6035-44ca-b531-a738670aa133', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('3a8450a7-6f5c-4972-82f6-3b4af620a1d6', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'ADA', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('3a8450a7-6f5c-4972-82f6-3b4af620a1d6', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('a32ea75d-38d7-42b5-84c5-0a00db2d8b1d', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'MATIC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('a32ea75d-38d7-42b5-84c5-0a00db2d8b1d', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('8b1a6d8a-44b8-4dee-8e53-6a69fa182f39', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c5af3271-7185-4a52-9d0c-1c4b418317d8', 'ATOM', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('8b1a6d8a-44b8-4dee-8e53-6a69fa182f39', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

/*
 * Demo1 Accounts
 */

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('8e2ee8eb-2057-4996-8d32-9eef6a6ab824', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('8e2ee8eb-2057-4996-8d32-9eef6a6ab824', 10000000, 0, 10000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('55bd611c-cc12-45b5-8596-53c3f18e176d', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'ETH', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('55bd611c-cc12-45b5-8596-53c3f18e176d', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('f950c105-8555-466c-9383-cc35bedc9b0f', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'SOL', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('f950c105-8555-466c-9383-cc35bedc9b0f', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('653e3e55-5299-4539-aee4-98f6507e26c3', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'BTC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('653e3e55-5299-4539-aee4-98f6507e26c3', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('9d4df5d5-8b83-4142-95bf-4c9487e93eec', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'ADA', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('9d4df5d5-8b83-4142-95bf-4c9487e93eec', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('f655baa0-c103-41a6-a7af-d02446dd1f2c', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'MATIC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('f655baa0-c103-41a6-a7af-d02446dd1f2c', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('acc5b856-8571-49c7-abbb-99c794f469a3', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '36ae23e7-d79e-4901-89d5-3fa87cca1abf', 'ATOM', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('acc5b856-8571-49c7-abbb-99c794f469a3', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);
