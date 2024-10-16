package pg_query_test

import (
	"testing"

	pg_query "github.com/fukasawah/go-pgquery"
)

// https://github.com/pganalyze/libpg_query/blob/16-latest/test/normalize_tests.c
var libpgqueryNormalizeTests = []string{
	"SELECT 1",
	"SELECT $1",
	"SELECT $1, 1",
	"SELECT $1, $2",
	"CREATE ROLE postgres PASSWORD 'xyz'",
	"CREATE ROLE postgres PASSWORD $1",
	"CREATE ROLE postgres ENCRYPTED PASSWORD 'xyz'",
	"CREATE ROLE postgres ENCRYPTED PASSWORD $1",
	"ALTER ROLE foo WITH PASSWORD 'bar' VALID UNTIL 'infinity'",
	"ALTER ROLE foo WITH PASSWORD $1 VALID UNTIL $2",
	"ALTER ROLE postgres LOGIN SUPERUSER ENCRYPTED PASSWORD 'xyz'",
	"ALTER ROLE postgres LOGIN SUPERUSER ENCRYPTED PASSWORD $1",
	"SELECT a, SUM(b) FROM tbl WHERE c = 'foo' GROUP BY 1, 'bar' ORDER BY 1, 'cafe'",
	"SELECT a, SUM(b) FROM tbl WHERE c = $1 GROUP BY 1, $2 ORDER BY 1, $3",
	"select date_trunc($1, created_at at time zone $2), count(*) from users group by date_trunc('day', created_at at time zone 'US/Pacific')",
	"select date_trunc($1, created_at at time zone $2), count(*) from users group by date_trunc($1, created_at at time zone $2)",
	"select count(1), date_trunc('day', created_at at time zone 'US/Pacific'), 'something', 'somethingelse' from users group by date_trunc('day', created_at at time zone 'US/Pacific'), date_trunc('day', created_at), 'foobar', 'abcdef'",
	"select count($1), date_trunc($2, created_at at time zone $3), $4, $5 from users group by date_trunc($2, created_at at time zone $3), date_trunc($6, created_at), $4, $5",
	"SELECT CAST('abc' as varchar(50))",
	"SELECT CAST($1 as varchar(50))",
	"CREATE OR REPLACE FUNCTION pg_temp.testfunc(OUT response \"mytable\", OUT sequelize_caught_exception text) RETURNS RECORD AS $func_12345$ BEGIN INSERT INTO \"mytable\" (\"mycolumn\") VALUES ('myvalue') RETURNING * INTO response; EXCEPTION WHEN unique_violation THEN GET STACKED DIAGNOSTICS sequelize_caught_exception = PG_EXCEPTION_DETAIL; END $func_12345$ LANGUAGE plpgsql; SELECT (testfunc.response).\"mycolumn\", testfunc.sequelize_caught_exception FROM pg_temp.testfunc(); DROP FUNCTION IF EXISTS pg_temp.testfunc();",
	"CREATE OR REPLACE FUNCTION pg_temp.testfunc(OUT response \"mytable\", OUT sequelize_caught_exception text) RETURNS RECORD AS $1 LANGUAGE plpgsql; SELECT (testfunc.response).\"mycolumn\", testfunc.sequelize_caught_exception FROM pg_temp.testfunc(); DROP FUNCTION IF EXISTS pg_temp.testfunc();",
	"CREATE PROCEDURE insert_data(a integer, b integer) LANGUAGE SQL AS $$ INSERT INTO tbl VALUES (a); INSERT INTO tbl VALUES (b); $$",
	"CREATE PROCEDURE insert_data(a integer, b integer) LANGUAGE SQL AS $1",
	"DO $$DECLARE r record; BEGIN FOR r IN SELECT table_schema, table_name FROM information_schema.tables WHERE table_type = 'VIEW' AND table_schema = 'public' LOOP EXECUTE 'GRANT ALL ON ' || quote_ident(r.table_schema) || '.' || quote_ident(r.table_name) || ' TO webuser'; END LOOP; END$$",
	"DO $1",
	"CREATE SUBSCRIPTION mysub CONNECTION 'host=192.168.1.50 port=5432 user=foo dbname=foodb' PUBLICATION mypublication, insert_only",
	"CREATE SUBSCRIPTION mysub CONNECTION $1 PUBLICATION mypublication, insert_only",
	"ALTER SUBSCRIPTION mysub SET PUBLICATION insert_only",
	"ALTER SUBSCRIPTION mysub SET PUBLICATION insert_only",
	"ALTER SUBSCRIPTION mysub CONNECTION 'host=192.168.1.50 port=5432 user=foo dbname=foodb'",
	"ALTER SUBSCRIPTION mysub CONNECTION $1",
	"CREATE USER MAPPING FOR bob SERVER foo OPTIONS (user 'bob', password 'secret')",
	"CREATE USER MAPPING FOR bob SERVER foo OPTIONS (user $1, password $2)",
	"ALTER USER MAPPING FOR bob SERVER foo OPTIONS (SET password 'public')",
	"ALTER USER MAPPING FOR bob SERVER foo OPTIONS (SET password $1)",
	"MERGE into measurement m USING new_measurement nm ON (m.city_id = nm.city_id and m.logdate=nm.logdate) WHEN MATCHED AND nm.peaktemp IS NULL THEN DELETE WHEN MATCHED THEN UPDATE SET peaktemp = greatest(m.peaktemp, nm.peaktemp), unitsales = m.unitsales + coalesce(nm.unitsales, 0) WHEN NOT MATCHED THEN INSERT (city_id, logdate, peaktemp, unitsales) VALUES (city_id, logdate, peaktemp, unitsales)",
	"MERGE into measurement m USING new_measurement nm ON (m.city_id = nm.city_id and m.logdate=nm.logdate) WHEN MATCHED AND nm.peaktemp IS NULL THEN DELETE WHEN MATCHED THEN UPDATE SET peaktemp = greatest(m.peaktemp, nm.peaktemp), unitsales = m.unitsales + coalesce(nm.unitsales, $1) WHEN NOT MATCHED THEN INSERT (city_id, logdate, peaktemp, unitsales) VALUES (city_id, logdate, peaktemp, unitsales)",
	// These below are as expected, though questionable if upstream shouldn't be
	// fixed as this could bloat pg_stat_statements
	"DECLARE cursor_b CURSOR FOR SELECT * FROM x WHERE id = 123",
	"DECLARE cursor_b CURSOR FOR SELECT * FROM x WHERE id = $1",
	"FETCH 1000 FROM cursor_a",
	"FETCH 1000 FROM cursor_a",
	"CLOSE cursor_a",
	"CLOSE cursor_a",
}

// https://github.com/pganalyze/libpg_query/blob/15-latest/test/normalize.c
func TestLibPgqueryNormalize(t *testing.T) {
	tests := libpgqueryNormalizeTests
	for i := 0; i < len(tests); i += 2 {
		input := tests[i]
		expected := tests[i+1]
		t.Run(input, func(t *testing.T) {
			actual, err := pg_query.Normalize(input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if actual != expected {
				t.Errorf("expected %q, got %q", expected, actual)
			}
		})
	}
}
