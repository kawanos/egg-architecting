package testutil

import (
	"context"
	"log"
	"os"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func InitData(ctx context.Context, db string, files []string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		log.Fatalf("Invalid database id %s", db)
	}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	var createTablesSQL []string
	for _, file := range files {
		sqlData, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		createTablesSQL = append(createTablesSQL, string(sqlData))
	}

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: createTablesSQL,
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func DropData(ctx context.Context, db string) error {

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		log.Fatalf("Invalid database id %s", db)
	}
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	err = adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{
		Database: db,
	})
	if err != nil {
		return err
	}
	return nil

}
