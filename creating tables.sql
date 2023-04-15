DROP TABLE supply;
DROP TABLE sale;
CREATE TABLE supply (
    id INTEGER PRIMARY KEY,
    barcode TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    supply_time TEXT NOT NULL,
    price REAL NOT NULL,
    sold_amount REAL NOT NULL DEFAULT 0
);

CREATE TABLE sale (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    barcode TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    sale_time DATETIME NOT NULL,
    price REAL NOT NULL,
    margin REAL NOT NULL DEFAULT 0
);