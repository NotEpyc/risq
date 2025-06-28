package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("=== Railway Environment Check ===")

	// Check DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		// Mask sensitive parts
		if strings.Contains(dbURL, "@") {
			parts := strings.Split(dbURL, "@")
			fmt.Printf("DATABASE_URL: postgres://...@%s\n", parts[len(parts)-1])
		} else {
			fmt.Println("DATABASE_URL: [FOUND but no @ symbol]")
		}
	} else {
		fmt.Println("DATABASE_URL: [NOT SET]")
	}

	// Check individual DB variables
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))
	fmt.Printf("DB_SSL_MODE: %s\n", os.Getenv("DB_SSL_MODE"))

	// Check PostgreSQL Railway variables
	fmt.Printf("PGHOST: %s\n", os.Getenv("PGHOST"))
	fmt.Printf("PGPORT: %s\n", os.Getenv("PGPORT"))
	fmt.Printf("PGUSER: %s\n", os.Getenv("PGUSER"))
	fmt.Printf("PGDATABASE: %s\n", os.Getenv("PGDATABASE"))

	// Check Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		fmt.Printf("REDIS_URL: %s\n", redisURL[:min(len(redisURL), 20)]+"...")
	} else {
		fmt.Println("REDIS_URL: [NOT SET]")
	}

	// Check app variables
	fmt.Printf("APP_ENV: %s\n", os.Getenv("APP_ENV"))
	fmt.Printf("PORT: %s\n", os.Getenv("PORT"))

	fmt.Println("=== End Environment Check ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
