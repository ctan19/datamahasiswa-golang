package main

func getDBConfig() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "kampus",
		SSLMode:  "disable",
	}
}