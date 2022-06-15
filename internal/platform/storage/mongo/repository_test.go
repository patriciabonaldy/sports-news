package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

var db *mongo.Client
var uri string

func TestMain(m *testing.M) {
	// setup
	pool, resource := mongoContainer()
	// run tests
	exitCode := m.Run()
	// kill and remove the container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	// disconnect mongodb client
	if err := db.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func TestNewDBStorage(t *testing.T) {
	_, err := NewDBStorage(context.Background(), &config.Database{}, logger.New())
	require.Error(t, err)

	_, err = NewDBStorage(context.Background(), &config.Database{
		URI:          uri,
		DatabaseName: "",
		User:         "root",
		Password:     "password",
	},
		logger.New())
	require.NoError(t, err)
}

func TestRepository_Save(t *testing.T) {
	repo := &Repository{
		databaseName: "test",
		db:           db,
	}

	article := mockArticle()
	ctx := context.Background()
	err := repo.Save(ctx, article)
	assert.NoError(t, err)

	got, err := repo.GetByID(ctx, article.NewsID)
	assert.NoError(t, err)

	want, err := repo.GetByID(ctx, got.NewsID)
	assert.NoError(t, err)
	assert.Equal(t, reflect.DeepEqual(want, got), true)
}

func mockArticle() internal.ArticleNews {
	article := internal.NewArticle()

	return article
}

func mongoContainer() (*dockertest.Pool, *dockertest.Resource) {
	const MONGO_INITDB_ROOT_USERNAME = "root"
	const MONGO_INITDB_ROOT_PASSWORD = "password"

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=" + MONGO_INITDB_ROOT_USERNAME,
		"MONGO_INITDB_ROOT_PASSWORD=" + MONGO_INITDB_ROOT_PASSWORD,
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env:        environmentVariables,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		var err error
		uri = fmt.Sprintf("mongodb://%s:%s@localhost:%s",
			MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD,
			resource.GetPort("27017/tcp"))
		db, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				uri,
			),
		)
		if err != nil {
			return err
		}
		return db.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	return pool, resource
}
