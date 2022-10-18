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

CREATE TYPE get_account_result AS
(
    account_id UUID,
    currency   TEXT,
    balance    NUMERIC,
    hold       NUMERIC,
    available  NUMERIC,
    created_at TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION get_balances_for_users(
    arg_user_id UUID
) RETURNS SETOF get_account_result
    LANGUAGE plpgsql
AS
$$
BEGIN
    return QUERY
        SELECT acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        FROM (select id, currency FROM account WHERE user_id = arg_user_id) acc
                 INNER JOIN
             (SELECT account_id, MAX(count) as max
              FROM account_balance
              WHERE account_id IN (select id
                                   FROM account
                                   WHERE user_id = arg_user_id)
              GROUP BY account_id) recent_balance
             ON acc.id = recent_balance.account_id
                 INNER JOIN
             account_balance ab
             ON recent_balance.account_id = ab.account_id and recent_balance.max = ab.count
        GROUP BY ab.count, acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        HAVING recent_balance.count = acc.count;
END
$$;