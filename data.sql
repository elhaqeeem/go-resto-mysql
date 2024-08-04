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
