package config

import "testing"

func TestConfig(t *testing.T) {
	cfg := LoadConfig("../../config")
	if cfg.PostgresDB.ReadDB.Database != "statusy" {
		t.Errorf("failed to parse postgres readdb database. expected=%s got=%s", "statusy", cfg.PostgresDB.ReadDB.Database)
	}

	if cfg.PostgresDB.ReadDB.Host != "localhost" {
		t.Errorf("failed to parse postgres readdb database. expected=%s got=%s", "localhost", cfg.PostgresDB.ReadDB.Host)
	}

	if cfg.PostgresDB.ReadDB.Password != "password" {
		t.Errorf("failed to parse postgres readdb database. expected=%s got=%s", "password", cfg.PostgresDB.ReadDB.Password)
	}

	if cfg.PostgresDB.ReadDB.User != "statusy" {
		t.Errorf("failed to parse postgres readdb database. expected=%s got=%s", "statusy", cfg.PostgresDB.ReadDB.User)
	}

	if cfg.PostgresDB.ReadDB.Port != 5432 {
		t.Errorf("failed to parse postgres readdb database. expected=%s got=%d", "5432", cfg.PostgresDB.ReadDB.Port)
	}

}
