INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES('B183F5E2-B72A-4AA5-B7AE-95E0D548D84D','D263E7E3-24D7-4C04-8D67-EA3A0BE7907E','620E62FD-DAF1-4738-84CE-1DBC4393ED29','wei',now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES('B183F5E2-B72A-4AA5-B7AE-95E0D548D84D',1000,0,1000,now(),null,0);

INSERT INTO account (id, portfolio_id, user_id, currency, created_at)
VALUES('9AA945E8-05FB-4D8E-88C9-1986F0813292','D263E7E3-24D7-4C04-8D67-EA3A0BE7907E','620E62FD-DAF1-4738-84CE-1DBC4393ED29','avax',now());

INSERT INTO account_balance(account_id, balance, hold, available, created_at, request_id, count)
VALUES('9AA945E8-05FB-4D8E-88C9-1986F0813292',1000,0,1000,now(),null,0);