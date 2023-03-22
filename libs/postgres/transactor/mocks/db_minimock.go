package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i route256/libs/postgres/transactor.DB -o ./mocks/db_minimock.go -n DBMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// DBMock implements transactor.DB
type DBMock struct {
	t minimock.Tester

	funcBeginTx          func(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error)
	inspectFuncBeginTx   func(ctx context.Context, txOptions pgx.TxOptions)
	afterBeginTxCounter  uint64
	beforeBeginTxCounter uint64
	BeginTxMock          mDBMockBeginTx

	funcExec          func(ctx context.Context, query string, args ...interface{}) (c2 pgconn.CommandTag, err error)
	inspectFuncExec   func(ctx context.Context, query string, args ...interface{})
	afterExecCounter  uint64
	beforeExecCounter uint64
	ExecMock          mDBMockExec

	funcQuery          func(ctx context.Context, query string, args ...interface{}) (r1 pgx.Rows, err error)
	inspectFuncQuery   func(ctx context.Context, query string, args ...interface{})
	afterQueryCounter  uint64
	beforeQueryCounter uint64
	QueryMock          mDBMockQuery

	funcQueryRow          func(ctx context.Context, query string, args ...interface{}) (r1 pgx.Row)
	inspectFuncQueryRow   func(ctx context.Context, query string, args ...interface{})
	afterQueryRowCounter  uint64
	beforeQueryRowCounter uint64
	QueryRowMock          mDBMockQueryRow
}

// NewDBMock returns a mock for transactor.DB
func NewDBMock(t minimock.Tester) *DBMock {
	m := &DBMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginTxMock = mDBMockBeginTx{mock: m}
	m.BeginTxMock.callArgs = []*DBMockBeginTxParams{}

	m.ExecMock = mDBMockExec{mock: m}
	m.ExecMock.callArgs = []*DBMockExecParams{}

	m.QueryMock = mDBMockQuery{mock: m}
	m.QueryMock.callArgs = []*DBMockQueryParams{}

	m.QueryRowMock = mDBMockQueryRow{mock: m}
	m.QueryRowMock.callArgs = []*DBMockQueryRowParams{}

	return m
}

type mDBMockBeginTx struct {
	mock               *DBMock
	defaultExpectation *DBMockBeginTxExpectation
	expectations       []*DBMockBeginTxExpectation

	callArgs []*DBMockBeginTxParams
	mutex    sync.RWMutex
}

// DBMockBeginTxExpectation specifies expectation struct of the DB.BeginTx
type DBMockBeginTxExpectation struct {
	mock    *DBMock
	params  *DBMockBeginTxParams
	results *DBMockBeginTxResults
	Counter uint64
}

// DBMockBeginTxParams contains parameters of the DB.BeginTx
type DBMockBeginTxParams struct {
	ctx       context.Context
	txOptions pgx.TxOptions
}

// DBMockBeginTxResults contains results of the DB.BeginTx
type DBMockBeginTxResults struct {
	t1  pgx.Tx
	err error
}

// Expect sets up expected params for DB.BeginTx
func (mmBeginTx *mDBMockBeginTx) Expect(ctx context.Context, txOptions pgx.TxOptions) *mDBMockBeginTx {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DBMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &DBMockBeginTxExpectation{}
	}

	mmBeginTx.defaultExpectation.params = &DBMockBeginTxParams{ctx, txOptions}
	for _, e := range mmBeginTx.expectations {
		if minimock.Equal(e.params, mmBeginTx.defaultExpectation.params) {
			mmBeginTx.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmBeginTx.defaultExpectation.params)
		}
	}

	return mmBeginTx
}

// Inspect accepts an inspector function that has same arguments as the DB.BeginTx
func (mmBeginTx *mDBMockBeginTx) Inspect(f func(ctx context.Context, txOptions pgx.TxOptions)) *mDBMockBeginTx {
	if mmBeginTx.mock.inspectFuncBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("Inspect function is already set for DBMock.BeginTx")
	}

	mmBeginTx.mock.inspectFuncBeginTx = f

	return mmBeginTx
}

// Return sets up results that will be returned by DB.BeginTx
func (mmBeginTx *mDBMockBeginTx) Return(t1 pgx.Tx, err error) *DBMock {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DBMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &DBMockBeginTxExpectation{mock: mmBeginTx.mock}
	}
	mmBeginTx.defaultExpectation.results = &DBMockBeginTxResults{t1, err}
	return mmBeginTx.mock
}

// Set uses given function f to mock the DB.BeginTx method
func (mmBeginTx *mDBMockBeginTx) Set(f func(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error)) *DBMock {
	if mmBeginTx.defaultExpectation != nil {
		mmBeginTx.mock.t.Fatalf("Default expectation is already set for the DB.BeginTx method")
	}

	if len(mmBeginTx.expectations) > 0 {
		mmBeginTx.mock.t.Fatalf("Some expectations are already set for the DB.BeginTx method")
	}

	mmBeginTx.mock.funcBeginTx = f
	return mmBeginTx.mock
}

// When sets expectation for the DB.BeginTx which will trigger the result defined by the following
// Then helper
func (mmBeginTx *mDBMockBeginTx) When(ctx context.Context, txOptions pgx.TxOptions) *DBMockBeginTxExpectation {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DBMock.BeginTx mock is already set by Set")
	}

	expectation := &DBMockBeginTxExpectation{
		mock:   mmBeginTx.mock,
		params: &DBMockBeginTxParams{ctx, txOptions},
	}
	mmBeginTx.expectations = append(mmBeginTx.expectations, expectation)
	return expectation
}

// Then sets up DB.BeginTx return parameters for the expectation previously defined by the When method
func (e *DBMockBeginTxExpectation) Then(t1 pgx.Tx, err error) *DBMock {
	e.results = &DBMockBeginTxResults{t1, err}
	return e.mock
}

// BeginTx implements transactor.DB
func (mmBeginTx *DBMock) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error) {
	mm_atomic.AddUint64(&mmBeginTx.beforeBeginTxCounter, 1)
	defer mm_atomic.AddUint64(&mmBeginTx.afterBeginTxCounter, 1)

	if mmBeginTx.inspectFuncBeginTx != nil {
		mmBeginTx.inspectFuncBeginTx(ctx, txOptions)
	}

	mm_params := &DBMockBeginTxParams{ctx, txOptions}

	// Record call args
	mmBeginTx.BeginTxMock.mutex.Lock()
	mmBeginTx.BeginTxMock.callArgs = append(mmBeginTx.BeginTxMock.callArgs, mm_params)
	mmBeginTx.BeginTxMock.mutex.Unlock()

	for _, e := range mmBeginTx.BeginTxMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.t1, e.results.err
		}
	}

	if mmBeginTx.BeginTxMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmBeginTx.BeginTxMock.defaultExpectation.Counter, 1)
		mm_want := mmBeginTx.BeginTxMock.defaultExpectation.params
		mm_got := DBMockBeginTxParams{ctx, txOptions}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmBeginTx.t.Errorf("DBMock.BeginTx got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmBeginTx.BeginTxMock.defaultExpectation.results
		if mm_results == nil {
			mmBeginTx.t.Fatal("No results are set for the DBMock.BeginTx")
		}
		return (*mm_results).t1, (*mm_results).err
	}
	if mmBeginTx.funcBeginTx != nil {
		return mmBeginTx.funcBeginTx(ctx, txOptions)
	}
	mmBeginTx.t.Fatalf("Unexpected call to DBMock.BeginTx. %v %v", ctx, txOptions)
	return
}

// BeginTxAfterCounter returns a count of finished DBMock.BeginTx invocations
func (mmBeginTx *DBMock) BeginTxAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.afterBeginTxCounter)
}

// BeginTxBeforeCounter returns a count of DBMock.BeginTx invocations
func (mmBeginTx *DBMock) BeginTxBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.beforeBeginTxCounter)
}

// Calls returns a list of arguments used in each call to DBMock.BeginTx.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmBeginTx *mDBMockBeginTx) Calls() []*DBMockBeginTxParams {
	mmBeginTx.mutex.RLock()

	argCopy := make([]*DBMockBeginTxParams, len(mmBeginTx.callArgs))
	copy(argCopy, mmBeginTx.callArgs)

	mmBeginTx.mutex.RUnlock()

	return argCopy
}

// MinimockBeginTxDone returns true if the count of the BeginTx invocations corresponds
// the number of defined expectations
func (m *DBMock) MinimockBeginTxDone() bool {
	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginTxMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBeginTx != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		return false
	}
	return true
}

// MinimockBeginTxInspect logs each unmet expectation
func (m *DBMock) MinimockBeginTxInspect() {
	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DBMock.BeginTx with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginTxMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		if m.BeginTxMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to DBMock.BeginTx")
		} else {
			m.t.Errorf("Expected call to DBMock.BeginTx with params: %#v", *m.BeginTxMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBeginTx != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		m.t.Error("Expected call to DBMock.BeginTx")
	}
}

type mDBMockExec struct {
	mock               *DBMock
	defaultExpectation *DBMockExecExpectation
	expectations       []*DBMockExecExpectation

	callArgs []*DBMockExecParams
	mutex    sync.RWMutex
}

// DBMockExecExpectation specifies expectation struct of the DB.Exec
type DBMockExecExpectation struct {
	mock    *DBMock
	params  *DBMockExecParams
	results *DBMockExecResults
	Counter uint64
}

// DBMockExecParams contains parameters of the DB.Exec
type DBMockExecParams struct {
	ctx   context.Context
	query string
	args  []interface{}
}

// DBMockExecResults contains results of the DB.Exec
type DBMockExecResults struct {
	c2  pgconn.CommandTag
	err error
}

// Expect sets up expected params for DB.Exec
func (mmExec *mDBMockExec) Expect(ctx context.Context, query string, args ...interface{}) *mDBMockExec {
	if mmExec.mock.funcExec != nil {
		mmExec.mock.t.Fatalf("DBMock.Exec mock is already set by Set")
	}

	if mmExec.defaultExpectation == nil {
		mmExec.defaultExpectation = &DBMockExecExpectation{}
	}

	mmExec.defaultExpectation.params = &DBMockExecParams{ctx, query, args}
	for _, e := range mmExec.expectations {
		if minimock.Equal(e.params, mmExec.defaultExpectation.params) {
			mmExec.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmExec.defaultExpectation.params)
		}
	}

	return mmExec
}

// Inspect accepts an inspector function that has same arguments as the DB.Exec
func (mmExec *mDBMockExec) Inspect(f func(ctx context.Context, query string, args ...interface{})) *mDBMockExec {
	if mmExec.mock.inspectFuncExec != nil {
		mmExec.mock.t.Fatalf("Inspect function is already set for DBMock.Exec")
	}

	mmExec.mock.inspectFuncExec = f

	return mmExec
}

// Return sets up results that will be returned by DB.Exec
func (mmExec *mDBMockExec) Return(c2 pgconn.CommandTag, err error) *DBMock {
	if mmExec.mock.funcExec != nil {
		mmExec.mock.t.Fatalf("DBMock.Exec mock is already set by Set")
	}

	if mmExec.defaultExpectation == nil {
		mmExec.defaultExpectation = &DBMockExecExpectation{mock: mmExec.mock}
	}
	mmExec.defaultExpectation.results = &DBMockExecResults{c2, err}
	return mmExec.mock
}

// Set uses given function f to mock the DB.Exec method
func (mmExec *mDBMockExec) Set(f func(ctx context.Context, query string, args ...interface{}) (c2 pgconn.CommandTag, err error)) *DBMock {
	if mmExec.defaultExpectation != nil {
		mmExec.mock.t.Fatalf("Default expectation is already set for the DB.Exec method")
	}

	if len(mmExec.expectations) > 0 {
		mmExec.mock.t.Fatalf("Some expectations are already set for the DB.Exec method")
	}

	mmExec.mock.funcExec = f
	return mmExec.mock
}

// When sets expectation for the DB.Exec which will trigger the result defined by the following
// Then helper
func (mmExec *mDBMockExec) When(ctx context.Context, query string, args ...interface{}) *DBMockExecExpectation {
	if mmExec.mock.funcExec != nil {
		mmExec.mock.t.Fatalf("DBMock.Exec mock is already set by Set")
	}

	expectation := &DBMockExecExpectation{
		mock:   mmExec.mock,
		params: &DBMockExecParams{ctx, query, args},
	}
	mmExec.expectations = append(mmExec.expectations, expectation)
	return expectation
}

// Then sets up DB.Exec return parameters for the expectation previously defined by the When method
func (e *DBMockExecExpectation) Then(c2 pgconn.CommandTag, err error) *DBMock {
	e.results = &DBMockExecResults{c2, err}
	return e.mock
}

// Exec implements transactor.DB
func (mmExec *DBMock) Exec(ctx context.Context, query string, args ...interface{}) (c2 pgconn.CommandTag, err error) {
	mm_atomic.AddUint64(&mmExec.beforeExecCounter, 1)
	defer mm_atomic.AddUint64(&mmExec.afterExecCounter, 1)

	if mmExec.inspectFuncExec != nil {
		mmExec.inspectFuncExec(ctx, query, args...)
	}

	mm_params := &DBMockExecParams{ctx, query, args}

	// Record call args
	mmExec.ExecMock.mutex.Lock()
	mmExec.ExecMock.callArgs = append(mmExec.ExecMock.callArgs, mm_params)
	mmExec.ExecMock.mutex.Unlock()

	for _, e := range mmExec.ExecMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.c2, e.results.err
		}
	}

	if mmExec.ExecMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmExec.ExecMock.defaultExpectation.Counter, 1)
		mm_want := mmExec.ExecMock.defaultExpectation.params
		mm_got := DBMockExecParams{ctx, query, args}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmExec.t.Errorf("DBMock.Exec got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmExec.ExecMock.defaultExpectation.results
		if mm_results == nil {
			mmExec.t.Fatal("No results are set for the DBMock.Exec")
		}
		return (*mm_results).c2, (*mm_results).err
	}
	if mmExec.funcExec != nil {
		return mmExec.funcExec(ctx, query, args...)
	}
	mmExec.t.Fatalf("Unexpected call to DBMock.Exec. %v %v %v", ctx, query, args)
	return
}

// ExecAfterCounter returns a count of finished DBMock.Exec invocations
func (mmExec *DBMock) ExecAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExec.afterExecCounter)
}

// ExecBeforeCounter returns a count of DBMock.Exec invocations
func (mmExec *DBMock) ExecBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExec.beforeExecCounter)
}

// Calls returns a list of arguments used in each call to DBMock.Exec.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmExec *mDBMockExec) Calls() []*DBMockExecParams {
	mmExec.mutex.RLock()

	argCopy := make([]*DBMockExecParams, len(mmExec.callArgs))
	copy(argCopy, mmExec.callArgs)

	mmExec.mutex.RUnlock()

	return argCopy
}

// MinimockExecDone returns true if the count of the Exec invocations corresponds
// the number of defined expectations
func (m *DBMock) MinimockExecDone() bool {
	for _, e := range m.ExecMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExec != nil && mm_atomic.LoadUint64(&m.afterExecCounter) < 1 {
		return false
	}
	return true
}

// MinimockExecInspect logs each unmet expectation
func (m *DBMock) MinimockExecInspect() {
	for _, e := range m.ExecMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DBMock.Exec with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecCounter) < 1 {
		if m.ExecMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to DBMock.Exec")
		} else {
			m.t.Errorf("Expected call to DBMock.Exec with params: %#v", *m.ExecMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExec != nil && mm_atomic.LoadUint64(&m.afterExecCounter) < 1 {
		m.t.Error("Expected call to DBMock.Exec")
	}
}

type mDBMockQuery struct {
	mock               *DBMock
	defaultExpectation *DBMockQueryExpectation
	expectations       []*DBMockQueryExpectation

	callArgs []*DBMockQueryParams
	mutex    sync.RWMutex
}

// DBMockQueryExpectation specifies expectation struct of the DB.Query
type DBMockQueryExpectation struct {
	mock    *DBMock
	params  *DBMockQueryParams
	results *DBMockQueryResults
	Counter uint64
}

// DBMockQueryParams contains parameters of the DB.Query
type DBMockQueryParams struct {
	ctx   context.Context
	query string
	args  []interface{}
}

// DBMockQueryResults contains results of the DB.Query
type DBMockQueryResults struct {
	r1  pgx.Rows
	err error
}

// Expect sets up expected params for DB.Query
func (mmQuery *mDBMockQuery) Expect(ctx context.Context, query string, args ...interface{}) *mDBMockQuery {
	if mmQuery.mock.funcQuery != nil {
		mmQuery.mock.t.Fatalf("DBMock.Query mock is already set by Set")
	}

	if mmQuery.defaultExpectation == nil {
		mmQuery.defaultExpectation = &DBMockQueryExpectation{}
	}

	mmQuery.defaultExpectation.params = &DBMockQueryParams{ctx, query, args}
	for _, e := range mmQuery.expectations {
		if minimock.Equal(e.params, mmQuery.defaultExpectation.params) {
			mmQuery.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmQuery.defaultExpectation.params)
		}
	}

	return mmQuery
}

// Inspect accepts an inspector function that has same arguments as the DB.Query
func (mmQuery *mDBMockQuery) Inspect(f func(ctx context.Context, query string, args ...interface{})) *mDBMockQuery {
	if mmQuery.mock.inspectFuncQuery != nil {
		mmQuery.mock.t.Fatalf("Inspect function is already set for DBMock.Query")
	}

	mmQuery.mock.inspectFuncQuery = f

	return mmQuery
}

// Return sets up results that will be returned by DB.Query
func (mmQuery *mDBMockQuery) Return(r1 pgx.Rows, err error) *DBMock {
	if mmQuery.mock.funcQuery != nil {
		mmQuery.mock.t.Fatalf("DBMock.Query mock is already set by Set")
	}

	if mmQuery.defaultExpectation == nil {
		mmQuery.defaultExpectation = &DBMockQueryExpectation{mock: mmQuery.mock}
	}
	mmQuery.defaultExpectation.results = &DBMockQueryResults{r1, err}
	return mmQuery.mock
}

// Set uses given function f to mock the DB.Query method
func (mmQuery *mDBMockQuery) Set(f func(ctx context.Context, query string, args ...interface{}) (r1 pgx.Rows, err error)) *DBMock {
	if mmQuery.defaultExpectation != nil {
		mmQuery.mock.t.Fatalf("Default expectation is already set for the DB.Query method")
	}

	if len(mmQuery.expectations) > 0 {
		mmQuery.mock.t.Fatalf("Some expectations are already set for the DB.Query method")
	}

	mmQuery.mock.funcQuery = f
	return mmQuery.mock
}

// When sets expectation for the DB.Query which will trigger the result defined by the following
// Then helper
func (mmQuery *mDBMockQuery) When(ctx context.Context, query string, args ...interface{}) *DBMockQueryExpectation {
	if mmQuery.mock.funcQuery != nil {
		mmQuery.mock.t.Fatalf("DBMock.Query mock is already set by Set")
	}

	expectation := &DBMockQueryExpectation{
		mock:   mmQuery.mock,
		params: &DBMockQueryParams{ctx, query, args},
	}
	mmQuery.expectations = append(mmQuery.expectations, expectation)
	return expectation
}

// Then sets up DB.Query return parameters for the expectation previously defined by the When method
func (e *DBMockQueryExpectation) Then(r1 pgx.Rows, err error) *DBMock {
	e.results = &DBMockQueryResults{r1, err}
	return e.mock
}

// Query implements transactor.DB
func (mmQuery *DBMock) Query(ctx context.Context, query string, args ...interface{}) (r1 pgx.Rows, err error) {
	mm_atomic.AddUint64(&mmQuery.beforeQueryCounter, 1)
	defer mm_atomic.AddUint64(&mmQuery.afterQueryCounter, 1)

	if mmQuery.inspectFuncQuery != nil {
		mmQuery.inspectFuncQuery(ctx, query, args...)
	}

	mm_params := &DBMockQueryParams{ctx, query, args}

	// Record call args
	mmQuery.QueryMock.mutex.Lock()
	mmQuery.QueryMock.callArgs = append(mmQuery.QueryMock.callArgs, mm_params)
	mmQuery.QueryMock.mutex.Unlock()

	for _, e := range mmQuery.QueryMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.r1, e.results.err
		}
	}

	if mmQuery.QueryMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmQuery.QueryMock.defaultExpectation.Counter, 1)
		mm_want := mmQuery.QueryMock.defaultExpectation.params
		mm_got := DBMockQueryParams{ctx, query, args}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmQuery.t.Errorf("DBMock.Query got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmQuery.QueryMock.defaultExpectation.results
		if mm_results == nil {
			mmQuery.t.Fatal("No results are set for the DBMock.Query")
		}
		return (*mm_results).r1, (*mm_results).err
	}
	if mmQuery.funcQuery != nil {
		return mmQuery.funcQuery(ctx, query, args...)
	}
	mmQuery.t.Fatalf("Unexpected call to DBMock.Query. %v %v %v", ctx, query, args)
	return
}

// QueryAfterCounter returns a count of finished DBMock.Query invocations
func (mmQuery *DBMock) QueryAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmQuery.afterQueryCounter)
}

// QueryBeforeCounter returns a count of DBMock.Query invocations
func (mmQuery *DBMock) QueryBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmQuery.beforeQueryCounter)
}

// Calls returns a list of arguments used in each call to DBMock.Query.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmQuery *mDBMockQuery) Calls() []*DBMockQueryParams {
	mmQuery.mutex.RLock()

	argCopy := make([]*DBMockQueryParams, len(mmQuery.callArgs))
	copy(argCopy, mmQuery.callArgs)

	mmQuery.mutex.RUnlock()

	return argCopy
}

// MinimockQueryDone returns true if the count of the Query invocations corresponds
// the number of defined expectations
func (m *DBMock) MinimockQueryDone() bool {
	for _, e := range m.QueryMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.QueryMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterQueryCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcQuery != nil && mm_atomic.LoadUint64(&m.afterQueryCounter) < 1 {
		return false
	}
	return true
}

// MinimockQueryInspect logs each unmet expectation
func (m *DBMock) MinimockQueryInspect() {
	for _, e := range m.QueryMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DBMock.Query with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.QueryMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterQueryCounter) < 1 {
		if m.QueryMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to DBMock.Query")
		} else {
			m.t.Errorf("Expected call to DBMock.Query with params: %#v", *m.QueryMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcQuery != nil && mm_atomic.LoadUint64(&m.afterQueryCounter) < 1 {
		m.t.Error("Expected call to DBMock.Query")
	}
}

type mDBMockQueryRow struct {
	mock               *DBMock
	defaultExpectation *DBMockQueryRowExpectation
	expectations       []*DBMockQueryRowExpectation

	callArgs []*DBMockQueryRowParams
	mutex    sync.RWMutex
}

// DBMockQueryRowExpectation specifies expectation struct of the DB.QueryRow
type DBMockQueryRowExpectation struct {
	mock    *DBMock
	params  *DBMockQueryRowParams
	results *DBMockQueryRowResults
	Counter uint64
}

// DBMockQueryRowParams contains parameters of the DB.QueryRow
type DBMockQueryRowParams struct {
	ctx   context.Context
	query string
	args  []interface{}
}

// DBMockQueryRowResults contains results of the DB.QueryRow
type DBMockQueryRowResults struct {
	r1 pgx.Row
}

// Expect sets up expected params for DB.QueryRow
func (mmQueryRow *mDBMockQueryRow) Expect(ctx context.Context, query string, args ...interface{}) *mDBMockQueryRow {
	if mmQueryRow.mock.funcQueryRow != nil {
		mmQueryRow.mock.t.Fatalf("DBMock.QueryRow mock is already set by Set")
	}

	if mmQueryRow.defaultExpectation == nil {
		mmQueryRow.defaultExpectation = &DBMockQueryRowExpectation{}
	}

	mmQueryRow.defaultExpectation.params = &DBMockQueryRowParams{ctx, query, args}
	for _, e := range mmQueryRow.expectations {
		if minimock.Equal(e.params, mmQueryRow.defaultExpectation.params) {
			mmQueryRow.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmQueryRow.defaultExpectation.params)
		}
	}

	return mmQueryRow
}

// Inspect accepts an inspector function that has same arguments as the DB.QueryRow
func (mmQueryRow *mDBMockQueryRow) Inspect(f func(ctx context.Context, query string, args ...interface{})) *mDBMockQueryRow {
	if mmQueryRow.mock.inspectFuncQueryRow != nil {
		mmQueryRow.mock.t.Fatalf("Inspect function is already set for DBMock.QueryRow")
	}

	mmQueryRow.mock.inspectFuncQueryRow = f

	return mmQueryRow
}

// Return sets up results that will be returned by DB.QueryRow
func (mmQueryRow *mDBMockQueryRow) Return(r1 pgx.Row) *DBMock {
	if mmQueryRow.mock.funcQueryRow != nil {
		mmQueryRow.mock.t.Fatalf("DBMock.QueryRow mock is already set by Set")
	}

	if mmQueryRow.defaultExpectation == nil {
		mmQueryRow.defaultExpectation = &DBMockQueryRowExpectation{mock: mmQueryRow.mock}
	}
	mmQueryRow.defaultExpectation.results = &DBMockQueryRowResults{r1}
	return mmQueryRow.mock
}

// Set uses given function f to mock the DB.QueryRow method
func (mmQueryRow *mDBMockQueryRow) Set(f func(ctx context.Context, query string, args ...interface{}) (r1 pgx.Row)) *DBMock {
	if mmQueryRow.defaultExpectation != nil {
		mmQueryRow.mock.t.Fatalf("Default expectation is already set for the DB.QueryRow method")
	}

	if len(mmQueryRow.expectations) > 0 {
		mmQueryRow.mock.t.Fatalf("Some expectations are already set for the DB.QueryRow method")
	}

	mmQueryRow.mock.funcQueryRow = f
	return mmQueryRow.mock
}

// When sets expectation for the DB.QueryRow which will trigger the result defined by the following
// Then helper
func (mmQueryRow *mDBMockQueryRow) When(ctx context.Context, query string, args ...interface{}) *DBMockQueryRowExpectation {
	if mmQueryRow.mock.funcQueryRow != nil {
		mmQueryRow.mock.t.Fatalf("DBMock.QueryRow mock is already set by Set")
	}

	expectation := &DBMockQueryRowExpectation{
		mock:   mmQueryRow.mock,
		params: &DBMockQueryRowParams{ctx, query, args},
	}
	mmQueryRow.expectations = append(mmQueryRow.expectations, expectation)
	return expectation
}

// Then sets up DB.QueryRow return parameters for the expectation previously defined by the When method
func (e *DBMockQueryRowExpectation) Then(r1 pgx.Row) *DBMock {
	e.results = &DBMockQueryRowResults{r1}
	return e.mock
}

// QueryRow implements transactor.DB
func (mmQueryRow *DBMock) QueryRow(ctx context.Context, query string, args ...interface{}) (r1 pgx.Row) {
	mm_atomic.AddUint64(&mmQueryRow.beforeQueryRowCounter, 1)
	defer mm_atomic.AddUint64(&mmQueryRow.afterQueryRowCounter, 1)

	if mmQueryRow.inspectFuncQueryRow != nil {
		mmQueryRow.inspectFuncQueryRow(ctx, query, args...)
	}

	mm_params := &DBMockQueryRowParams{ctx, query, args}

	// Record call args
	mmQueryRow.QueryRowMock.mutex.Lock()
	mmQueryRow.QueryRowMock.callArgs = append(mmQueryRow.QueryRowMock.callArgs, mm_params)
	mmQueryRow.QueryRowMock.mutex.Unlock()

	for _, e := range mmQueryRow.QueryRowMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.r1
		}
	}

	if mmQueryRow.QueryRowMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmQueryRow.QueryRowMock.defaultExpectation.Counter, 1)
		mm_want := mmQueryRow.QueryRowMock.defaultExpectation.params
		mm_got := DBMockQueryRowParams{ctx, query, args}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmQueryRow.t.Errorf("DBMock.QueryRow got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmQueryRow.QueryRowMock.defaultExpectation.results
		if mm_results == nil {
			mmQueryRow.t.Fatal("No results are set for the DBMock.QueryRow")
		}
		return (*mm_results).r1
	}
	if mmQueryRow.funcQueryRow != nil {
		return mmQueryRow.funcQueryRow(ctx, query, args...)
	}
	mmQueryRow.t.Fatalf("Unexpected call to DBMock.QueryRow. %v %v %v", ctx, query, args)
	return
}

// QueryRowAfterCounter returns a count of finished DBMock.QueryRow invocations
func (mmQueryRow *DBMock) QueryRowAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmQueryRow.afterQueryRowCounter)
}

// QueryRowBeforeCounter returns a count of DBMock.QueryRow invocations
func (mmQueryRow *DBMock) QueryRowBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmQueryRow.beforeQueryRowCounter)
}

// Calls returns a list of arguments used in each call to DBMock.QueryRow.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmQueryRow *mDBMockQueryRow) Calls() []*DBMockQueryRowParams {
	mmQueryRow.mutex.RLock()

	argCopy := make([]*DBMockQueryRowParams, len(mmQueryRow.callArgs))
	copy(argCopy, mmQueryRow.callArgs)

	mmQueryRow.mutex.RUnlock()

	return argCopy
}

// MinimockQueryRowDone returns true if the count of the QueryRow invocations corresponds
// the number of defined expectations
func (m *DBMock) MinimockQueryRowDone() bool {
	for _, e := range m.QueryRowMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.QueryRowMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterQueryRowCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcQueryRow != nil && mm_atomic.LoadUint64(&m.afterQueryRowCounter) < 1 {
		return false
	}
	return true
}

// MinimockQueryRowInspect logs each unmet expectation
func (m *DBMock) MinimockQueryRowInspect() {
	for _, e := range m.QueryRowMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DBMock.QueryRow with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.QueryRowMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterQueryRowCounter) < 1 {
		if m.QueryRowMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to DBMock.QueryRow")
		} else {
			m.t.Errorf("Expected call to DBMock.QueryRow with params: %#v", *m.QueryRowMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcQueryRow != nil && mm_atomic.LoadUint64(&m.afterQueryRowCounter) < 1 {
		m.t.Error("Expected call to DBMock.QueryRow")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *DBMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockBeginTxInspect()

		m.MinimockExecInspect()

		m.MinimockQueryInspect()

		m.MinimockQueryRowInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *DBMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *DBMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockBeginTxDone() &&
		m.MinimockExecDone() &&
		m.MinimockQueryDone() &&
		m.MinimockQueryRowDone()
}