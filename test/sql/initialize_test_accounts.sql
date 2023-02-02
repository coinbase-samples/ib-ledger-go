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

