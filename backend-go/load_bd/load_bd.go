package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	kzip "github.com/klauspost/compress/zip"
	"github.com/snabb/httpreaderat"
)

const (
	fiasURL   = "https://fias-file.nalog.ru/downloads/2026.04.07/gar_xml.zip"
	targetReg = "59"
	batchSize = 5000
	workers   = 12
)

func getConnStr() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	return "postgres://postgres:postgres@" + host + ":5432/postgres?sslmode=disable"
}

func main() {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(getConnStr())
	if err != nil {
		log.Fatalf("Ошибка конфигурации БД: %v", err)
	}
	config.MaxConns = int32(workers * 4)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	fmt.Println("Подключение к архиву...")
	req, _ := http.NewRequest("GET", fiasURL, nil)
	htr, err := httpreaderat.New(nil, req, nil)
	if err != nil {
		log.Fatalf("Ошибка HTTP ридера: %v", err)
	}

	reader, err := kzip.NewReader(htr, htr.Size())
	if err != nil {
		log.Fatalf("Ошибка ZIP ридера: %v", err)
	}

	targets := map[string]string{
		"AS_ADDR_OBJ_":      "addr_obj",
		"AS_ADM_HIERARCHY_": "hierarchy",
		"AS_HOUSES_":        "houses",
	}

	var wgFiles sync.WaitGroup

	for prefix, table := range targets {
		for _, zf := range reader.File {
			if strings.HasPrefix(zf.Name, targetReg+"/") &&
				strings.Contains(zf.Name, prefix) &&
				!strings.Contains(zf.Name, "_PARAMS") {

				wgFiles.Add(1)
				go func(f *kzip.File, tName string) {
					defer wgFiles.Done()
					processFileParallel(ctx, pool, f, tName)
				}(zf, table)
			}
		}
	}

	wgFiles.Wait()
	fmt.Println("\n--- Все данные успешно загружены! ---")
}

func processFileParallel(ctx context.Context, pool *pgxpool.Pool, zf *kzip.File, tableName string) {
	jobs := make(chan []any, 10000)
	var wgWorkers sync.WaitGroup

	for w := 0; w < workers; w++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			dbWorker(ctx, pool, tableName, jobs)
		}()
	}

	rc, err := zf.Open()
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v", zf.Name, err)
		close(jobs)
		return
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	count := 0
	startTime := time.Now()

	fmt.Printf("[START] %s -> %s\n", zf.Name, tableName)

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if se, ok := t.(xml.StartElement); ok {
			var row []any
			switch tableName {
			case "addr_obj":
				if se.Name.Local == "OBJECT" {
					var o struct {
						ID string `xml:"OBJECTID,attr"`
						Nm string `xml:"NAME,attr"`
						Tp string `xml:"TYPENAME,attr"`
						Lv string `xml:"LEVEL,attr"`
						Ac string `xml:"ISACTIVE,attr"`
					}
					decoder.DecodeElement(&o, &se)
					if o.Ac == "1" {
						row = []any{o.ID, o.Tp + " " + o.Nm, o.Lv}
					}
				}
			case "hierarchy":
				if se.Name.Local == "ITEM" {
					var h struct {
						ID  string `xml:"OBJECTID,attr"`
						Pid string `xml:"PARENTOBJID,attr"`
						Ac  string `xml:"ISACTIVE,attr"`
					}
					decoder.DecodeElement(&h, &se)
					if h.Ac == "1" && h.Pid != "" {
						row = []any{h.ID, h.Pid}
					}
				}
			case "houses":
				if se.Name.Local == "HOUSE" {
					var hs struct {
						ID string `xml:"OBJECTID,attr"`
						Nr string `xml:"HOUSENUM,attr"`
						Ac string `xml:"ISACTIVE,attr"`
					}
					decoder.DecodeElement(&hs, &se)
					if hs.Ac == "1" {
						row = []any{hs.ID, hs.Nr}
					}
				}
			}

			if row != nil {
				jobs <- row
				count++
				if count%10000 == 0 {
					elapsed := time.Since(startTime).Seconds()
					rps := float64(count) / elapsed
					fmt.Printf("[%s] Загружено: %d | Скорость: %.0f строк/сек\n", tableName, count, rps)
				}
			}
		}
	}

	close(jobs)
	wgWorkers.Wait()
	fmt.Printf("[FINISH] %s | Всего: %d | Время: %v\n", zf.Name, count, time.Since(startTime).Round(time.Second))
}

func dbWorker(ctx context.Context, pool *pgxpool.Pool, table string, jobs <-chan []any) {
	var cols []string
	switch table {
	case "addr_obj":
		cols = []string{"id", "name", "level"}
	case "hierarchy":
		cols = []string{"id", "parent_id"}
	case "houses":
		cols = []string{"id", "house_num"}
	}

	batch := make([][]any, 0, batchSize)

	for row := range jobs {
		batch = append(batch, row)
		if len(batch) >= batchSize {
			_, err := pool.CopyFrom(ctx, pgx.Identifier{table}, cols, pgx.CopyFromRows(batch))
			if err != nil {
				log.Printf("Ошибка COPY %s: %v", table, err)
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		pool.CopyFrom(ctx, pgx.Identifier{table}, cols, pgx.CopyFromRows(batch))
	}
}
