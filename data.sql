-- Membuat database baru
CREATE DATABASE restoran;
USE restoran;

-- Tabel untuk kategori
CREATE TABLE kategori (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(50) NOT NULL
);

-- Menambahkan kategori
INSERT INTO kategori (nama) VALUES ('Minuman');
INSERT INTO kategori (nama) VALUES ('Makanan');
INSERT INTO kategori (nama) VALUES ('Promo');

-- Tabel untuk minuman
CREATE TABLE minuman (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(50) NOT NULL,
    varian VARCHAR(50),
    harga DECIMAL(10,2) NOT NULL,
    kategori_id INT,
    FOREIGN KEY (kategori_id) REFERENCES kategori(id)
);

-- Menambahkan minuman
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Jeruk', 'Dingin', 12000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Jeruk', 'Panas', 10000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Teh', 'Manis', 8000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Teh', 'Tawar', 5000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Kopi', 'Dingin', 8000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Kopi', 'Panas', 6000, (SELECT id FROM kategori WHERE nama='Minuman'));
INSERT INTO minuman (nama, varian, harga, kategori_id) VALUES ('Es Batu', NULL, 2000, (SELECT id FROM kategori WHERE nama='Minuman'));

-- Tabel untuk makanan
CREATE TABLE makanan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(50) NOT NULL,
    varian VARCHAR(50),
    harga DECIMAL(10,2) NOT NULL,
    kategori_id INT,
    FOREIGN KEY (kategori_id) REFERENCES kategori(id)
);

-- Menambahkan makanan
INSERT INTO makanan (nama, varian, harga, kategori_id) VALUES ('Mie', 'Goreng', 15000, (SELECT id FROM kategori WHERE nama='Makanan'));
INSERT INTO makanan (nama, varian, harga, kategori_id) VALUES ('Mie', 'Kuah', 15000, (SELECT id FROM kategori WHERE nama='Makanan'));
INSERT INTO makanan (nama, varian, harga, kategori_id) VALUES ('Nasi Goreng', NULL, 15000, (SELECT id FROM kategori WHERE nama='Makanan'));

-- Tabel untuk promo
CREATE TABLE promo (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    harga DECIMAL(10,2) NOT NULL
);

-- Menambahkan promo
INSERT INTO promo (nama, harga) VALUES ('Nasi Goreng + Jeruk Dingin', 23000);

-- Tabel untuk printer
CREATE TABLE printer (
    id CHAR(1) PRIMARY KEY,
    nama VARCHAR(50) NOT NULL
);

-- Menambahkan printer
INSERT INTO printer (id, nama) VALUES ('A', 'Printer Kasir');
INSERT INTO printer (id, nama) VALUES ('B', 'Printer Dapur (Makanan)');
INSERT INTO printer (id, nama) VALUES ('C', 'Printer Bar (Minuman)');

-- Tabel untuk meja
CREATE TABLE meja (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nomor INT NOT NULL
);

-- Menambahkan meja
INSERT INTO meja (nomor) VALUES (1);
INSERT INTO meja (nomor) VALUES (2);
INSERT INTO meja (nomor) VALUES (3);

-- Tabel untuk orders (pesanan)
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    meja_id INT,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (meja_id) REFERENCES meja(id)
);

-- Tabel untuk order_items (item dalam pesanan)
CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    item_type ENUM('Minuman', 'Makanan', 'Promo') NOT NULL,
    item_id INT,
    jumlah INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (item_id) REFERENCES (
        SELECT id FROM minuman WHERE item_type = 'Minuman'
        UNION ALL
        SELECT id FROM makanan WHERE item_type = 'Makanan'
        UNION ALL
        SELECT id FROM promo WHERE item_type = 'Promo'
    )
);

-- Menambahkan pesanan
INSERT INTO orders (meja_id) VALUES ((SELECT id FROM meja WHERE nomor = 1));

-- Menambahkan item pesanan
INSERT INTO order_items (order_id, item_type, item_id, jumlah) VALUES (
    (SELECT id FROM orders WHERE meja_id = (SELECT id FROM meja WHERE nomor = 1)),
    'Minuman',
    (SELECT id FROM minuman WHERE nama = 'Es Batu'),
    1
);

INSERT INTO order_items (order_id, item_type, item_id, jumlah) VALUES (
    (SELECT id FROM orders WHERE meja_id = (SELECT id FROM meja WHERE nomor = 1)),
    'Minuman',
    (SELECT id FROM minuman WHERE nama = 'Kopi' AND varian = 'Panas'),
    1
);

INSERT INTO order_items (order_id, item_type, item_id, jumlah) VALUES (
    (SELECT id FROM orders WHERE meja_id = (SELECT id FROM meja WHERE nomor = 1)),
    'Promo',
    (SELECT id FROM promo WHERE nama = 'Nasi Goreng + Jeruk Dingin'),
    2
);

INSERT INTO order_items (order_id, item_type, item_id, jumlah) VALUES (
    (SELECT id FROM orders WHERE meja_id = (SELECT id FROM meja WHERE nomor = 1)),
    'Minuman',
    (SELECT id FROM minuman WHERE nama = 'Teh' AND varian = 'Manis'),
    1
);

INSERT INTO order_items (order_id, item_type, item_id, jumlah) VALUES (
    (SELECT id FROM orders WHERE meja_id = (SELECT id FROM meja WHERE nomor = 1)),
    'Makanan',
    (SELECT id FROM makanan WHERE nama = 'Mie' AND varian = 'Goreng'),
    1
);

ALTER TABLE order_items
ADD INDEX idx_order_id (order_id),
ADD INDEX idx_item_type (item_type),
ADD INDEX idx_item_id (item_id);
