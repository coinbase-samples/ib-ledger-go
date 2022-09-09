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