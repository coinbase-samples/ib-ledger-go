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

DROP INDEX user_accounts;
DROP INDEX idem_user_account;
DROP INDEX account_balance_index;
DROP INDEX sender_transactions;
DROP INDEX receiver_transactions;
DROP INDEX transaction_entries;
DROP INDEX request_entries;
DROP INDEX account_entries;
DROP INDEX transaction_holds;
DROP INDEX request_holds;
DROP INDEX account_holds;
DROP TABLE entry;
DROP TABLE finalized_transaction;
DROP TABLE released_hold;
DROP TABLE hold;
DROP TABLE transaction;
DROP TABLE account_balance;
DROP TABLE account;
DROP TYPE direction;
DROP TYPE ttype;

DROP EXTENSION "uuid-ossp";
