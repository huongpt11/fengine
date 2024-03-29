package sql_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	l "log"
	"os"
	"testing"

	. "github.com/duclmse/fengine/fengine/db/sql"
	"github.com/goccy/go-json"

	"github.com/google/uuid"

	"github.com/duclmse/fengine/pkg/logger"
)

func connect(t *testing.T) (logger.Logger, *pgxpool.Pool) {
	config := Config{Host: "localhost", Port: "5432", User: "postgres", Pass: "1", Name: "postgres", SSLMode: "disable"}
	log, err := logger.New(os.Stdout, "debug")
	if err != nil {
		t.Fatalf("cannot create logger %s", err)
		return nil, nil
	}
	db, err := Connect(config, log)
	if err != nil {
		log.Fatalf("cannot connect %s", err)
	}
	return log, db
}

func TestGenUUID(t *testing.T) {
	if id, err := uuid.NewRandom(); err != nil {
		fmt.Printf("err")
	} else {
		fmt.Printf("uuid:= %s\n", id)
	}
}

func TestUuid(t *testing.T) {
	log, db := connect(t)
	rows, err := db.Query(context.Background(), `SELECT * FROM entity WHERE id = $1::UUID`, `21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab`)
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	defer rows.Close()

	things := make([]Entity, 0)
	for rows.Next() {
		e := &Entity{}
		//FIXME
		//if err := rows.StructScan(e); err != nil {
		//	t.Errorf("err %s", err.Error())
		//}
		things = append(things, *e)
	}
	t.Logf("result len %d", len(things))
	log.Struct(things)
}

//func TestFengineRepository_GetThingService(t *testing.T) {
//	log, db := connect(t)
//	repository := NewFEngineRepository(db, log)
//	service, err := repository.GetThingService(context.Background(), ThingServiceId{
//		EntityId: uuid.MustParse("1d6d5123-3fb8-4ab1-956f-c6f96847471d"),
//		Name:     "templ_method",
//	})
//	if err != nil {
//		t.Errorf("err %s\n", err.Error())
//		return
//	}
//	log.Struct(service)
//}

//func TestFengineRepository_GetAllThingService(t *testing.T) {
//	log, db := connect(t)
//	repository := NewFEngineRepository(db, log)
//	service, err := repository.GetThingAllServices(context.Background(), uuid.MustParse("21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab"))
//	if err != nil {
//		t.Errorf("err %s\n", err.Error())
//		return
//	}
//	fmt.Printf("%+v\n", service)
//}

func TestPanic(t *testing.T) {
	a := func(msg string) {
		defer func() {
			if err := recover(); err != nil {
				l.Printf("err %s\n", err)
			}
		}()
		panic(msg)
	}
	a("1")
	a("2")
	a("3")
}

func TestGeneratedQuery(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"table":   "tbl_test",
		"fields":  ["id", "name as n", "description d", "a", "b", "c"],
		"filter":  {
			"$and": [
				{"a": {"$gt": 10, "$lt": 20}},
				{"$or": [{"b": {"$gt": 50}}, {"b": {"$lt": 20}}, {"c": {"$in": [123, 234]}}]}
			]
		},
		"limit":   1000,
		"offset":  1,
		"order_by": ["name"]
	}`)
	req := new(SelectRequest)
	if err := json.Unmarshal(jsonb, &req); err != nil {
		t.Logf("error unmarshalling req: %s", err.Error())
		t.FailNow()
	}

	sql, err := req.ToSQL()
	if err != nil {
		t.Fatalf("cannot generate query command %s", err)
	}
	t.Logf("%s\n", sql)

	log, db := connect(t)
	rows, err := db.Query(context.Background(), sql)
	if err != nil {
		log.Info("err=%v", err)
	}
	cols := Columns(rows)

	fmt.Printf("%+v\n", cols)
	for rows.Next() {
		rowMap := make(map[string]interface{})
		MapScan(rows, rowMap)
		t.Logf("%v\n", rowMap)
	}
}
