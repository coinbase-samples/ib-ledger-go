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

CREATE OR REPLACE FUNCTION get_account_and_latest_balance(
    arg_user_id UUID,
    arg_currency TEXT
) RETURNS SETOF get_account_result
    LANGUAGE plpgsql
AS
$$
BEGIN
    return QUERY
        SELECT acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        FROM (select id, currency FROM account WHERE user_id = arg_user_id AND currency = arg_currency) acc
                 INNER JOIN
             (SELECT account_id, balance, hold, available, created_at FROM account_balance HAVING count = MAX(count)) ab
             ON acc.id = ab.account_id;
END
$$;
