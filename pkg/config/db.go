package config

type Database struct {
	ConnectionString string `yaml:"connection_string" env:"DB_CONNECTION_STRING"`
	MigrationsDir    string `yaml:"migration_dir"     env:"DB_MIGRATION_DIR"`
	PoolMaxConns     uint   `yaml:"pool_max_conns"    env:"DB_POOL_MAX_CONNS"    env-default:"20"`
	PoolMinConns     uint   `yaml:"pool_min_conns"    env:"DB_POOL_MIN_CONNS"    env-default:"5"`

	PoolHealthCheckPeriodMilliseconds uint `yaml:"pool_health_check_period_milliseconds" env:"DB_POOL_HEALTH_CHECK_PERIOD_MILLISECONDS" env-default:"60000"`
	PoolMaxConnIdleTimeMilliseconds   uint `yaml:"pool_max_conn_idle_time_milliseconds"  env:"DB_POOL_MAX_CONN_IDLE_TIME_MILLISECONDS"  env-default:"60000"`
	PoolMaxConnLifetimeMilliseconds   uint `yaml:"pool_max_conn_lifetime_milliseconds"   env:"DB_POOL_MAX_CONN_LIFETIME_MILLISECONDS"   env-default:"180000"`
}
