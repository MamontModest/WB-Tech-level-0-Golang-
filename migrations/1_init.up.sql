CREATE TABLE ORDERS(
    orderUid  VARCHAR(63) PRIMARY KEY,
    trackNumber  VARCHAR(31) NOT NULL UNIQUE,
    entry  VARCHAR(31) NOT NULL,
    locale  VARCHAR(3) NOT NULL,
    internalSignature  VARCHAR(63) NOT NULL,
    customerId  VARCHAR(31) NOT NULL,
    deliveryService  VARCHAR(31) NOT NULL,
    shardkey  VARCHAR(31) NOT NULL,
    smId BIGINT NOT NULL,
    dateCreated TIMESTAMP NOT NULL,
    oofShard  VARCHAR(31) NOT NULL
);

CREATE TABLE PAYMENTS(
    orderId VARCHAR(30) PRIMARY KEY REFERENCES ORDERS (orderUid),
    transaction VARCHAR(50) NOT NULL,
    requestId VARCHAR(30) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(30) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    paymentDt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    deliveryCost DECIMAL(10, 2) NOT NULL,
    goodsTotal INTEGER NOT NULL,
    customFee DECIMAL(10, 2) NOT NULL
);

CREATE TABLE ITEM(
     id SERIAL PRIMARY KEY,
     orderId  VARCHAR(31) REFERENCES ORDERS (orderUid),
     chrtId BIGINT NOT NULL,
     trackNumber VARCHAR(63) NOT NULL,
     price DECIMAL(10, 2) NOT NULL,
     rid VARCHAR(63) NOT NULL,
     name VARCHAR(63) NOT NULL,
     sale INTEGER CHECK(sale >=0 AND sale <= 100) NOT NULL,
     size VARCHAR(15) NOT NULL,
     totalPrice DECIMAL(10, 2) NOT NULL,
     nmId BIGINT NOT NULL,
     brand VARCHAR(31) NOT NULL,
     status INTEGER NOT NULL
);

CREATE TABLE DELIVERY(
     orderId VARCHAR(30) PRIMARY KEY REFERENCES ORDERS(orderUid),
     name VARCHAR(50) NOT NULL,
     phone VARCHAR(30) NOT NULL,
     zip VARCHAR(10) NOT NULL,
     city VARCHAR(30) NOT NULL,
     address VARCHAR(50) NOT NULL,
     region VARCHAR(50) NOT NULL,
     email VARCHAR(50) NOT NULL
);