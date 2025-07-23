CREATE TABLE IF NOT EXISTS pairs
(
    pair character varying(8) NOT NULL,
    ask_unit character varying(4) NOT NULL,
    bid_unit character varying(4) NOT NULL,
    min_ask bigint NOT NULL,
    min_bid bigint NOT NULL,
    maker_fee bigint NOT NULL,
    taker_fee bigint NOT NULL,
    price_precision integer NOT NULL,
    volume_precision integer NOT NULL,
    created_at bigint NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at bigint NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    PRIMARY KEY (pair)
);