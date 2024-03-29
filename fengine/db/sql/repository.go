package sql

import (
	ctx "context"
	"fmt"

	"github.com/duclmse/fengine/pkg/logger"
	. "github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ Repository = (*fengineRepository)(nil)

type Repository interface {
	GetEntity(ctx ctx.Context, id UUID) (*EntityDefinition, error)
	UpsertEntity(ctx ctx.Context, def EntityDefinition) (int64, error)
	DeleteEntity(ctx ctx.Context, thingId UUID) (int64, error)

	GetThingAllServices(ctx ctx.Context, thingId UUID) ([]EntityService, error)
	GetThingService(ctx ctx.Context, id ThingServiceId) (*EntityService, error)
	UpsertThingService(ctx ctx.Context, service ...ThingService) (int, error)
	DeleteThingService(ctx ctx.Context, id ThingServiceId) (int, error)

	GetThingAllSubscriptions(ctx ctx.Context, thingId UUID) ([]EntitySubscription, error)
	GetThingSubscriptions(ctx ctx.Context, id ThingSubscriptionId) (*EntitySubscription, error)
	UpsertThingSubscription(ctx ctx.Context, sub ...ThingSubscription) (int64, error)
	DeleteThingSubscription(ctx ctx.Context, id ThingSubscriptionId) (int64, error)

	GetThingAttributes(ctx ctx.Context, attrs ...string) ([]Variable, error)
	SetThingAttributes(ctx ctx.Context, attrs []Variable) (int64, error)
	GetAttributeHistory(cts ctx.Context, attrs AttributeHistoryRequest) ([]Variable, error)

	Select(ctx ctx.Context, sql string, params ...any) (r []map[string]Variable, e error)
	Insert(ctx ctx.Context, sql string, params ...any) (r int64, e error)
	Update(ctx ctx.Context, sql string, params ...any) (r int64, e error)
	Delete(ctx ctx.Context, sql string, params ...any) (r int64, e error)
}

type RowMapper func(pgx.Rows) error

// NewFEngineRepository instantiates a PostgresSQL implementation of PricingRepository
func NewFEngineRepository(db *pgxpool.Pool, log logger.Logger) Repository {
	return &fengineRepository{
		db:  db,
		log: log,
	}
}

type fengineRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func (fer fengineRepository) GetEntity(ctx ctx.Context, thingId UUID) (*EntityDefinition, error) {
	// language=sql
	query := `SELECT "id", "name", "type", "description", "project_id", "base_template", "base_shapes", "create_ts",
       "update_ts" FROM entity WHERE id = $1`
	rows, err := fer.db.Query(ctx, query, thingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		def := new(EntityDefinition)
		//FIXME err := rows.StructScan(def)
		def.BaseShapes, err = def.BaseShapesStr.ToUuidArray()
		return def, err
	}
	return nil, nil
}

func (fer fengineRepository) UpsertEntity(ctx ctx.Context, def EntityDefinition) (int64, error) {
	// language=sql
	query := `INSERT INTO entity("id", "name", "type", "description", "project_id", "base_template", "base_shapes"
 		) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO UPDATE SET base_template = $6, base_shapes = $7`
	res, err := fer.db.Exec(ctx, query, def.Id, def.Name, def.Type, def.Description, def.ProjectId,
		def.BaseTemplate, def.BaseShapes)
	if err != nil {
		return 0, err
	}
	affected := res.RowsAffected()
	ts, err := def.ToThingServices()
	if err != nil {
		return 0, err
	}
	subs, err := def.ToThingSubscriptions()
	if err != nil {
		return 0, err
	}
	_, err = fer.UpsertThingService(ctx, ts...)
	if err != nil {
		return 0, err
	}
	_, err = fer.UpsertThingSubscription(ctx, subs...)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (fer fengineRepository) DeleteEntity(ctx ctx.Context, thingId UUID) (int64, error) {
	// language=postgresql
	result, err := fer.db.Exec(ctx, `DELETE FROM entity WHERE id = $1::UUID`, thingId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (fer fengineRepository) GetThingAllServices(ctx ctx.Context, thingId UUID) ([]EntityService, error) {
	// language=postgresql
	query := `SELECT m1.entity_id AS id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "service" m1 LEFT JOIN "service" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = $1::UUID`
	entities, err := fer.db.Query(ctx, query, thingId)
	if err != nil {
		return nil, err
	}
	defer entities.Close()

	result := []EntityService{}
	for entities.Next() {
		entity := EntityService{}
		//if err := entities.StructScan(&entity); err != nil {
		//	return nil, err
		//}
		result = append(result, entity)
	}
	return result, nil
}

func (fer fengineRepository) GetThingService(ctx ctx.Context, id ThingServiceId) (*EntityService, error) {
	// language=postgresql
	query := `SELECT m1.entity_id AS id, m1.name, m1."input", m1."output", m1."from",
    	CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
		FROM "service" m1 LEFT JOIN "service" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
		WHERE m1.entity_id = $1::UUID AND m1.name = $2`
	entities, err := fer.db.Query(ctx, query, id.EntityId, id.Name)
	if err != nil {
		fmt.Printf("err selecting %s", err.Error())
		return nil, err
	}
	defer entities.Close()

	for entities.Next() {
		result := new(EntityService)
		//if err := entities.StructScan(result); err != nil {
		//	return nil, err
		//}
		return result, nil
	}

	return nil, nil
}

func (fer fengineRepository) UpsertThingService(ctx ctx.Context, service ...ThingService) (int, error) {
	// language=postgresql
	query := `INSERT INTO service("entity_id", "name", "input", "output", "code") VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO UPDATE SET "input" = $3, "output" = $4, "code" = $5, update_ts = NOW()`
	result, err := fer.db.Exec(ctx, query, service)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (fer fengineRepository) DeleteThingService(ctx ctx.Context, id ThingServiceId) (int, error) {
	// language=postgresql
	query := `DELETE FROM service s WHERE s.entity_id = $1::UUID AND s.name = $2;`
	result, err := fer.db.Exec(ctx, query, id.EntityId, id.Name)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (fer fengineRepository) GetThingAllSubscriptions(ctx ctx.Context, thingId UUID) (subs []EntitySubscription, err error) {
	// language=postgresql
	query := `SELECT "entity_id", "name", "subs_on", "event", "from", "code", "create_ts", "update_ts" FROM subscription
		WHERE entity_id = $1::UUID`
	rows, err := fer.db.Query(ctx, query, thingId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {

	}
	return
}

func (fer fengineRepository) GetThingSubscriptions(ctx ctx.Context, id ThingSubscriptionId) (*EntitySubscription, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) UpsertThingSubscription(ctx ctx.Context, sub ...ThingSubscription) (int64, error) {
	// language=postgresql
	query := `INSERT INTO subscription("entity_id", "name", "event", "subs_on", "code") VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO UPDATE SET "event" = $3, "subs_on" = $4, "code" = $5, update_ts = NOW()`
	result, err := fer.db.Exec(ctx, query, sub)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (fer fengineRepository) DeleteThingSubscription(ctx ctx.Context, id ThingSubscriptionId) (int64, error) {
	// language=postgresql
	query := `DELETE FROM "subscription" s WHERE s.entity_id = $1::UUID AND s.name = $2;`
	result, err := fer.db.Exec(ctx, query, id.EntityId, id.Name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (fer fengineRepository) GetThingAttributes(ctx ctx.Context, attrs ...string) ([]Variable, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) SetThingAttributes(ctx ctx.Context, attrs []Variable) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) GetAttributeHistory(cts ctx.Context, attrs AttributeHistoryRequest) ([]Variable, error) {
	//TODO implement me
	panic("implement me")
}

func (fer fengineRepository) Select(ctx ctx.Context, sql string, params ...any) (r []map[string]Variable, err error) {
	fmt.Printf("param.l = %d: %v\n", len(params), params)
	rows, err := fer.db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		return
	}
	result := []map[string]Variable{}
	dess := rows.FieldDescriptions()
	for rows.Next() {
		row := map[string]Variable{}
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		for i, des := range dess {
			row[string(des.Name)] = Variable{Value: values[i]}
		}
		result = append(result, row)
	}
	return result, nil
}

func (fer fengineRepository) Insert(ctx ctx.Context, sql string, params ...any) (r int64, err error) {
	inserted, err := fer.db.Exec(ctx, sql, params)
	if err != nil {
		return
	}
	return inserted.RowsAffected(), nil
}

func (fer fengineRepository) Update(ctx ctx.Context, sql string, params ...any) (r int64, err error) {
	updated, err := fer.db.Exec(ctx, sql, params)
	if err != nil {
		return
	}
	return updated.RowsAffected(), nil
}

func (fer fengineRepository) Delete(ctx ctx.Context, sql string, params ...any) (r int64, err error) {
	deleted, err := fer.db.Exec(ctx, sql, params)
	if err != nil {
		return 0, err
	}
	return deleted.RowsAffected(), nil
}
