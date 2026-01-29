CREATE TABLE subscriptions (
    subscription_id UUID PRIMARY KEY,
    service_name VARCHAR NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);