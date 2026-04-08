<?php
require_once 'config.php';

// Если пришел поисковый запрос — отдаем JSON
if (isset($_GET['q'])) {
    header('Content-Type: application/json');
    $query = $_GET['q'];

    if (strlen($query) < 3) {
        echo json_encode([]);
        exit;
    }

    try {
        $dsn = "pgsql:host=" . DB_HOST . ";port=" . DB_PORT . ";dbname=" . DB_NAME;
        $pdo = new PDO($dsn, DB_USER, DB_PASS, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

        $words = explode(' ', trim($query));
        $conditions = [];
        $params = [];

        foreach ($words as $word) {
            if (empty($word)) continue;
            $conditions[] = "full_address ILIKE ?";
            $params[] = "%$word%";
        }

        $sql = "SELECT full_address FROM addresses 
                WHERE " . implode(' AND ', $conditions) . " 
                LIMIT 10";
        
        $stmt = $pdo->prepare($sql);
        $stmt->execute($params);
        
        echo json_encode($stmt->fetchAll(PDO::FETCH_COLUMN));
    } catch (PDOException $e) {
        echo json_encode(["error" => $e->getMessage()]);
    }
    exit;
}

// Если запроса нет — просто отдаем HTML-файл
include 'index.html';