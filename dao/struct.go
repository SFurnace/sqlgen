package dao

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	"github.com/huandu/go-sqlbuilder"
)

var globalStructMap = new(sync.Map)

// Struct 对 sqlbuilder.Struct 进行了封装，使其更易使用
type Struct struct {
	*sqlbuilder.Struct
	typ reflect.Type
}

// S ...
func S(val interface{}) *Struct {
	typ := dereferencedType(reflect.TypeOf(val))
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("invalid value: %v", val))
	}
	if v, ok := globalStructMap.Load(typ); ok {
		return v.(*Struct)
	}

	v := &Struct{Struct: sqlbuilder.NewStruct(reflect.New(typ).Interface()), typ: typ}
	globalStructMap.Store(typ, v)
	return v
}

// scanRow ...
func (s *Struct) scanRow(row *sql.Row, tag string, destPtr interface{}) error {
	dTyp := reflect.TypeOf(destPtr)
	if dTyp.Kind() != reflect.Ptr || dTyp.Elem() != s.typ {
		return fmt.Errorf("invalid dest type: %v", dTyp)
	}

	if err := row.Scan(s.AddrForTag(tag, destPtr)...); err != nil {
		return err
	}
	return nil
}

// scanRows ...
func (s *Struct) scanRows(rows *sql.Rows, tag string, destPtr interface{}) error {
	dTyp := reflect.TypeOf(destPtr)
	if dTyp.Kind() != reflect.Ptr || dTyp.Elem().Kind() != reflect.Slice || dTyp.Elem().Elem() != s.typ {
		return fmt.Errorf("invalid dest type: %v", dTyp)
	}

	var (
		dVal = reflect.ValueOf(destPtr).Elem()
		err  error
	)
	for rows.Next() {
		tmp := reflect.New(s.typ)
		if err = rows.Scan(s.AddrForTag(tag, tmp.Interface())...); err != nil {
			return err
		}
		dVal.Set(reflect.Append(dVal, tmp.Elem()))
	}
	if rows.Err() != nil {
		return err
	}
	return nil
}

// Query ...
func (s *Struct) Query(ctx context.Context, db Executor, result interface{}, expr string, args ...interface{}) error {
	return s.TagQuery(ctx, db, result, "", expr, args...)
}

// QueryB ...
func (s *Struct) QueryB(ctx context.Context, db Executor, result interface{}, b sqlbuilder.Builder) error {
	expr, args := b.Build()
	return s.TagQuery(ctx, db, result, "", expr, args...)
}

// TagQuery ...
func (s *Struct) TagQuery(
	ctx context.Context, db Executor, result interface{}, tag, expr string, args ...interface{},
) error {
	rows, err := Query(ctx, db, expr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err = s.scanRows(rows, tag, result); err != nil {
		// ecmlog.ErrorEx(ctx, "scanRows failed", "expr", expr, "args", args)
		return err
	}
	return nil
}

// TagQueryB ...
func (s *Struct) TagQueryB(
	ctx context.Context, db Executor, result interface{}, tag string, b sqlbuilder.Builder,
) error {
	expr, args := b.Build()
	return s.TagQuery(ctx, db, result, tag, expr, args...)
}

// QueryRow ...
func (s *Struct) QueryRow(
	ctx context.Context, db Executor, result interface{}, expr string, args ...interface{},
) error {
	return s.TagQueryRow(ctx, db, result, "", expr, args...)
}

// QueryRowB ...
func (s *Struct) QueryRowB(ctx context.Context, db Executor, result interface{}, b sqlbuilder.Builder) error {
	expr, args := b.Build()
	return s.TagQueryRow(ctx, db, result, "", expr, args...)
}

// TagQueryRow ...
func (s *Struct) TagQueryRow(
	ctx context.Context, db Executor, result interface{}, tag, expr string, args ...interface{},
) error {
	err := s.scanRow(QueryRow(ctx, db, expr, args...), tag, result)
	if err != nil {
		// ecmlog.ErrorEx(ctx, "scanRow failed", "expr", expr, "args", args)
	}
	return err
}

// TagQueryRowB ...
func (s *Struct) TagQueryRowB(
	ctx context.Context, db Executor, result interface{}, tag string, b sqlbuilder.Builder,
) error {
	expr, args := b.Build()
	return s.TagQueryRow(ctx, db, result, tag, expr, args...)
}

// Exec ...
func (s *Struct) Exec(ctx context.Context, db Executor, expr string, args ...interface{}) (sql.Result, error) {
	return Exec(ctx, db, expr, args...)
}

// ExecB ...
func (s *Struct) ExecB(ctx context.Context, db Executor, b sqlbuilder.Builder) (sql.Result, error) {
	expr, args := b.Build()
	return s.Exec(ctx, db, expr, args...)
}
