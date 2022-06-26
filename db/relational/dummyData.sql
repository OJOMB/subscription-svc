USE subscriptions;

INSERT INTO products VALUES
    ('7Yn_IvvYsfkeo7-ysixd7','premium','Access all content for 30 days',30,30.00,'2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W','basic','Access basic content for 30 days',30,20.00,'2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('UMp3k41eV5mY_iOkiElGm','premium','Access all content for 1 year',365,300.00,'2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('mApy9b9Fqpt_WjghgUkSY','basic','Access basic content for 1 year',365,200.00,'2021-11-28 00:00:00','2021-11-28 00:00:00');

INSERT INTO users VALUES
    ('b3lFU5zF9zB37DxKk-zCC','Christopher','Ebert','morar.doug@example.org','2020-08-05','2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('FUQQzY_-4Tv_p7SFeHJEI','Emmie','Mayert','fhilpert@example.org','1994-07-26','2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('Yh4WkMFE0MTZnUjD9aYFy','Sven','Smith','rempel.ulices@example.org','2006-04-18','2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('yne6EtUQ9azjDXGLGZFs8','Claudie','Pouros','derrick.schroeder@example.org','2018-01-24','2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('521Qlk96BbJoGKseV1nPZ','Vergie','Carroll','ntillman@example.com','1999-04-25','2021-11-28 00:00:00','2021-11-28 00:00:00'),
    /* the following users have no subscription-plan */
    ('qo6_0keqGKqDA9EB5obql','Jeffrey','Jordan','jjordan123@example.com','1991-04-23','2021-11-28 00:00:00','2021-11-28 00:00:00'),
    ('MFE0MTZqGKqDA9EB5obql','Ben','Fisher','benf@example.com','1967-04-23','2021-11-28 00:00:00','2021-11-28 00:00:00');

INSERT INTO vouchers VALUES
    ('10-percent-off','PERCENTAGE', 10.00, '2021-11-28 00:00:00', '2025-05-20 00:00:00', 5, '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('10-off','FIXED_AMOUNT', 10.00, '2021-11-28 00:00:00', '2025-05-20 00:00:00', 5, '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('im-an-expired-voucher','FIXED_AMOUNT', 20.00, '2021-11-28 00:00:00', '2021-11-28 00:01:00', 5, '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('im-a-used-up-voucher','FIXED_AMOUNT', 20.00, '2021-11-28 00:00:00', '2025-05-20 00:00:00', 0, '2021-11-28 00:00:00', '2021-11-28 00:00:00');

/* here vouchers only apply to the 30 day products */
INSERT INTO products_vouchers VALUES
    ('7Yn_IvvYsfkeo7-ysixd7', '10-percent-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('7Yn_IvvYsfkeo7-ysixd7', '10-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('7Yn_IvvYsfkeo7-ysixd7', '20-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('7Yn_IvvYsfkeo7-ysixd7', 'im-an-expired-voucher', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('7Yn_IvvYsfkeo7-ysixd7', 'im-a-used-up-voucher', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W', '10-percent-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W', '10-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W', '20-off', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W', 'im-an-expired-voucher', '2021-11-28 00:00:00', '2021-11-28 00:00:00'),
    ('Ft2GgLRgN3FbMveklTy-W', 'im-a-used-up-voucher', '2021-11-28 00:00:00', '2021-11-28 00:00:00');

INSERT INTO subscription_plans VALUES
    (id,                       user_id,                product_id,              status,    start_date,            end_date,             net_price, gross_price, tax,   discount, voucher_code,     created_at,            updated_at)
    ('rhmeplLbg8bxWsqLzZQ6i', 'b3lFU5zF9zB37DxKk-zCC', '7Yn_IvvYsfkeo7-ysixd7', 'ACTIVE',  '2022-06-01 00:00:00', '2022-07-01 00:00:00', 31.50,    30.00,       1.50,  0,        NULL,             '2022-06-01 00:00:00', '2022-06-01 00:00:00'),
    ('sFF_eBjQgBcTZAKPCSzo5', 'FUQQzY_-4Tv_p7SFeHJEI', 'Ft2GgLRgN3FbMveklTy-W', 'ACTIVE',  '2022-06-01 00:00:00', '2022-07-01 00:00:00', 21.00,    20.00,       1.00,  0,        NULL,             '2022-06-01 00:00:00', '2022-06-01 00:00:00'),
    /* the following user (Yh4WkMFE0MTZnUjD9aYFy) has 2 plans on record - one is active and the other expired */
    ('sFtJdT5uTuYsWpLVIAt-m', 'Yh4WkMFE0MTZnUjD9aYFy', 'Ft2GgLRgN3FbMveklTy-W', 'EXPIRED', '2021-11-01 00:00:00', '2021-12-01 00:00:00', 18.90,    20.00,       1.00,  2.00,     '10-percent-off', '2021-11-28 00:00:00', '2021-12-01 00:00:00'),
    ('KlRuHEGEseQBogzXpc8ns', 'Yh4WkMFE0MTZnUjD9aYFy', 'Ft2GgLRgN3FbMveklTy-W', 'ACTIVE',  '2022-06-23 00:00:00', '2022-07-23 00:00:00', 18.90,    20.00,       1.00,  2.00,     '10-percent-off', '2022-06-23 00:00:00', '2022-06-23 00:00:00'),
    ('DDkg0NNvkwzV6miDLZX5E', 'yne6EtUQ9azjDXGLGZFs8', 'UMp3k41eV5mY_iOkiElGm', 'PAUSED',  '2022-01-01 00:00:00', '2023-01-31 00:00:00', 210.00,   200.00,      10.00, 0,        NULL,             '2022-01-01 00:00:00', '2022-04-02 00:00:00'),
    ('LwBhyg4FcWwUH1ORGVXDU', '521Qlk96BbJoGKseV1nPZ', 'mApy9b9Fqpt_WjghgUkSY', 'ACTIVE',  '2022-01-01 00:00:00', '2023-01-01 00:00:00', 315.00,   300.00,      15.00, 0,        NULL,             '2022-01-01 00:00:00', '2022-01-01 00:00:00');

INSERT INTO subscription_plan_pauses VALUES
    ('ouDToOusald2HiXATnZbf', 'DDkg0NNvkwzV6miDLZX5E', '2022-01-31 00:00:00', '2023-01-01 00:00:00', '2022-03-02 00:00:00','2022-01-31 00:00:00', '2022-03-02 00:00:00'),
    ('ouDToOusald2HiXATnZbf', 'DDkg0NNvkwzV6miDLZX5E', '2022-04-02 00:00:00', '2023-01-31 00:00:00', NULL,'2022-04-02 00:00:00', '2022-04-02 00:00:00');
