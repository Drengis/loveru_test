package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const conn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer pool.Close()

	fmt.Println("[STEP 1] Очистка таблицы addresses...")
	_, err = pool.Exec(ctx, "TRUNCATE TABLE addresses")
	if err != nil {
		log.Fatalf("Ошибка очистки: %v", err)
	}

	fmt.Println("[STEP 2] Сборка полных адресов (это может занять пару минут)...")
	start := time.Now()

	// Cобираем цепочку Дом -> Улица -> Город с правильной пунктуацией
	query := `
    INSERT INTO addresses (full_address, region, city, street, house, external_id)
    SELECT 
        -- Логика для Города: если точки нет, добавляем её к сокращению
        (CASE 
            WHEN city.name ~ '\.' THEN city.name 
            ELSE REGEXP_REPLACE(city.name, '^([гпсд])\s', '\1. ', 'i') 
        END) || ', ' ||
        -- Логика для Улицы: если точки нет, добавляем её к сокращению
        (CASE 
            WHEN street.name ~ '\.' THEN street.name 
            ELSE REGEXP_REPLACE(street.name, '^(ул|пер|пр|бульв|ш)\s', '\1. ', 'i') 
        END) || ', д. ' || h.house_num,
        '59',
        city.name,
        street.name,
        h.house_num,
        h.id
    FROM houses h
    JOIN hierarchy h_to_s ON h.id = h_to_s.id
    JOIN addr_obj street ON h_to_s.parent_id = street.id
    JOIN hierarchy s_to_c ON street.id = s_to_c.id
    JOIN addr_obj city ON s_to_c.parent_id = city.id
    WHERE street.level IN ('7', '8') 
      AND city.level <= '6';
    `

	res, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Ошибка сборки: %v", err)
	}

	fmt.Printf("[FINISH] Сборка завершена за %v\n", time.Since(start).Round(time.Second))
	fmt.Printf("Добавлено записей: %d\n", res.RowsAffected())
}
