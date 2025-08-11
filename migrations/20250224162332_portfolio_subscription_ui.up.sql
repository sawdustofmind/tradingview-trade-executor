CREATE UNIQUE INDEX portfolio_subscription_ui
    ON portfolio_subscription (customer_id, portfolio_id, exchange, is_test)
    WHERE portfolio_subscription.status = 'active';