<?php
require_once 'config.php';

try {
    $dsn = "pgsql:host=" . DB_HOST . ";port=" . DB_PORT . ";dbname=" . DB_NAME;
    $pdo = new PDO($dsn, DB_USER, DB_PASS, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

    $sql = "
    DROP TABLE IF EXISTS addresses;
    CREATE TABLE addresses (
        id SERIAL PRIMARY KEY,
        full_address TEXT,
        region VARCHAR(100),
        city VARCHAR(200),
        street VARCHAR(200),
        house VARCHAR(200),
        external_id VARCHAR(50),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    DROP TABLE IF EXISTS addr_obj;
    CREATE TABLE addr_obj (
        id VARCHAR(50) PRIMARY KEY,
        name TEXT,
        level VARCHAR(10)
    );

    DROP TABLE IF EXISTS hierarchy;
    CREATE TABLE hierarchy (
        id VARCHAR(50) PRIMARY KEY,
        parent_id VARCHAR(50)
    );

    DROP TABLE IF EXISTS houses;
    CREATE TABLE houses (
        id VARCHAR(50) PRIMARY KEY,
        house_num VARCHAR(50)
    );

    CREATE EXTENSION IF NOT EXISTS pg_trgm;
    CREATE INDEX IF NOT EXISTS idx_addresses_full_search ON addresses USING gin (full_address gin_trgm_ops);

    CREATE INDEX IF NOT EXISTS idx_addr_obj_id ON addr_obj(id);
    CREATE INDEX IF NOT EXISTS idx_hierarchy_id ON hierarchy(id);
    CREATE INDEX IF NOT EXISTS idx_hierarchy_parent ON hierarchy(parent_id);
    CREATE INDEX IF NOT EXISTS idx_houses_id ON houses(id);
    ";

    $pdo->exec($sql);
    echo "[PHP] Таблицы обновлены.\n";
} catch (PDOException $e) {
    die("[PHP ERROR] " . $e->getMessage() . "\n");
}