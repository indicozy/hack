CREATE TABLE supply (
    id INTEGER PRIMARY KEY,
    barcode TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    datetime TEXT NOT NULL,
    price REAL NOT NULL,
    sold_amount REAL NOT NULL
);

CREATE TABLE sale (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    barcode TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    datetime DATETIME NOT NULL,
    price REAL NOT NULL,
    margin REAL NOT NULL
);

CREATE TABLE supply_sale (
    supply_id INTEGER PRIMARY KEY,
    sale_id INTEGER PRIMARY KEY,
    supply_quantity INTEGER NOT NULL,
);

INSERT INTO supply (id, barcode, quantity, datetime, price, sold_amount)
VALUES
    (1, 'ABC123', 10, '2022-01-01 10:00:00', 2, 0),
    (2, 'DEF456', 20, '2022-01-02 10:00:00', 3, 0),
    (3, 'GHI789', 30, '2022-01-03 10:00:00', 4, 0),
    (4, 'JKL012', 40, '2022-01-04 10:00:00', 5, 0),
    (5, 'MNO345', 50, '2022-01-05 10:00:00', 6, 0);
