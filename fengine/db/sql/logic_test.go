package sql_test

import (
	"encoding/json"
	. "github.com/duclmse/fengine/fengine/db/sql"
	"log"
	"testing"

	pb "github.com/duclmse/fengine/pb"
)

func Test_DynamicUnmarshall(t *testing.T) {
	checkUnmarshall(t, `{
		"$and": [
			{"a": {"$gt": 10, "$lt": 20}},
			{"$or": [{"b": {"$gt": 50}},{"b": {"$lt": 20}},{"c": {"$in": ["abc", "def", 123]}}]}
		]
	}`)

	checkUnmarshall(t, `{"a": {"$gt": 10}}`)
}

func checkUnmarshall(t *testing.T, s string) {
	filter := Filter{}
	if err := json.Unmarshal([]byte(s), &filter); err != nil {
		t.Errorf("err> %v\n", err)
		return
	}

	if sb, err := filter.BuildLogic(); err == nil {
		t.Logf("%s\n", sb.String())
	}
}

func TestSelectRequest_ToSQL(t *testing.T) {
	// language=json
	jsonb := []byte(`{
		"table":   "tbl_test",
		"fields":  ["id", "name as n", "description"],
		"filter":  {
			"$and": [
				{"a": {"$gt": 10, "$lt": 20}},
				{"$or": [
					{"b": {"$gt": 50}},
					{"b": {"$lt": 20}},
					{"c": {"$in": ["abc", "def", 123]}}
				]}
			]
		},
		"group_by": ["name"],
		"limit":   1000,
		"offset":  10,
		"order_by": ["name"]
	}`)
	req := SelectRequest{}
	if err := json.Unmarshal(jsonb, &req); err != nil {
		log.Printf("error unmarshalling req: %s", err.Error())
		t.FailNow()
	}

	sql, err := req.ToSQL()
	if err != nil {
		t.FailNow()
		return
	}
	t.Logf("%s\n", sql)
}

func TestTableDefinition_ToSQL(t *testing.T) {
	def := TableDefinition{
		Name: "test",
		Fields: []Field{
			{Name: "id", Type: pb.Type_i32, IsPrimaryKey: true, IsLogged: false},
			{Name: "name", Type: pb.Type_string, IsPrimaryKey: false, IsLogged: false},
		},
	}
	sql, err := def.ToSQL()
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	t.Logf("%s\n", sql)
}
