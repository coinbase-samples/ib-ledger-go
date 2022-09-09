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

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('B183F5E2-B72A-4AA5-B7AE-95E0D548D84D', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '620E62FD-DAF1-4738-84CE-1DBC4393ED29', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('B183F5E2-B72A-4AA5-B7AE-95E0D548D84D', 100000, 0, 100000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('9AA945E8-05FB-4D8E-88C9-1986F0813292', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '620E62FD-DAF1-4738-84CE-1DBC4393ED29', 'ETH', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('9AA945E8-05FB-4D8E-88C9-1986F0813292', 100000, 0, 100000, now(), null, 1);

/*
Carl's User
*/
INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('266f4811-8906-4292-8e8a-ce18ebe33d1b', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('266f4811-8906-4292-8e8a-ce18ebe33d1b', 10000000, 0, 10000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('100a4079-3237-4302-8bd5-44b5425d0823', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'ETH', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('100a4079-3237-4302-8bd5-44b5425d0823', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('cc367c26-c192-437c-8fa3-43ce39ca32ef', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'SOL', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('cc367c26-c192-437c-8fa3-43ce39ca32ef', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('04ca5482-9f3a-42df-b70d-78ae13cec589', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'BTC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('04ca5482-9f3a-42df-b70d-78ae13cec589', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('997e24a9-4377-4df4-9315-30f99ecfa269', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'ADA', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('997e24a9-4377-4df4-9315-30f99ecfa269', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('6d8f8af0-5e06-48f1-9ff6-ad3d95dc6837', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'MATIC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('6d8f8af0-5e06-48f1-9ff6-ad3d95dc6837', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('502ff8be-99ce-494a-9e27-54bbf2a37fd6', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        'c7e34d37-f678-4096-94f7-3cad7d3258b9', 'ATOM', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('502ff8be-99ce-494a-9e27-54bbf2a37fd6', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);


/*
Jay's Test Account - 4f5a6336-8101-4634-a458-73b7f6fcf49f
*/

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('61e3221e-9434-448f-9cd3-73d32f42a647', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'USD', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('61e3221e-9434-448f-9cd3-73d32f42a647', 10000000, 0, 10000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('c71f9a2f-2522-4ec1-bc68-bc697fcd7b17', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'ETH', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('c71f9a2f-2522-4ec1-bc68-bc697fcd7b17', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('1934cfb7-90a4-4e84-a615-f91db1c757ba', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'SOL', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('1934cfb7-90a4-4e84-a615-f91db1c757ba', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('8b75d9cb-756f-4e12-a342-7a2b512126af', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'BTC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('8b75d9cb-756f-4e12-a342-7a2b512126af', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('6eb0d347-2dd9-45aa-8da2-fca26467eef9', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'ADA', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('6eb0d347-2dd9-45aa-8da2-fca26467eef9', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('744b8753-ed2b-4a9a-8d7e-6e6c4e89de8c', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'MATIC', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('744b8753-ed2b-4a9a-8d7e-6e6c4e89de8c', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES ('2660a259-8c5a-406d-9ca2-e1466f3ab2e4', 'D263E7E3-24D7-4C04-8D67-EA3A0BE7907E',
        '4f5a6336-8101-4634-a458-73b7f6fcf49f', 'ATOM', now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES ('2660a259-8c5a-406d-9ca2-e1466f3ab2e4', 15000000000000000000, 0, 15000000000000000000, now(), null, 1);