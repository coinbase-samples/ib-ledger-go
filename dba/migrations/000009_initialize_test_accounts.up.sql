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
        '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('0adbb104-fc18-46ca-a4eb-beee7775eb69', 10000, 0, 10000, now(), null, 1);

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'ETH');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'SOL');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'BTC');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'ADA');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'MATIC');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '37d10e18-34a2-4bd2-b7bc-b8e6dd6358f1', 'ATOM');

/*
 * Demo1 Accounts
 */

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('8e2ee8eb-2057-4996-8d32-9eef6a6ab824', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('8e2ee8eb-2057-4996-8d32-9eef6a6ab824', 10000, 0, 10000, now(), null, 0);

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'ETH');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'SOL');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'BTC');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'ADA');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'MATIC');

SELECT *
FROM initialize_account('D263E7E3-24D7-4C04-8D67-EA3A0BE7907E', '6a56c92e-10f5-4f05-9f3f-414cc98c7292', 'ATOM');
