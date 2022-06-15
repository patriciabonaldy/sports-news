package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

const eventCollectionName = "event"

// Repository is a mongo EventRepository implementation.
type Repository struct {
	databaseName string
	db           *mongo.Client
	log          logger.Logger
}

// NewDBStorage initializes a mongo-based implementation of Storage.
func NewDBStorage(ctx context.Context, cfg *config.Database, log logger.Logger) (*Repository, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(cfg.URI).
			SetAuth(options.Credential{Username: cfg.User,
				Password: cfg.Password}))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Repository{
		databaseName: cfg.DatabaseName,
		db:           client,
		log:          log,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, articleID string) (internal.ArticleNews, error) {
	var result ArticleNews

	err := r.getCollection(eventCollectionName).
		FindOne(ctx, bson.M{"answer_id": articleID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return internal.ArticleNews{}, internal.ErrArticleNotFound
		}

		return internal.ArticleNews{}, err
	}

	return parseToBusinessArticleNews(result), nil
}

func (r *Repository) Save(ctx context.Context, article internal.ArticleNews) error {
	articleDB := parseToArticleNewsDB(article)
	_, err := r.getCollection(eventCollectionName).InsertOne(ctx, articleDB)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) getCollection(collectionName string) *mongo.Collection {
	return r.db.Database(r.databaseName).Collection(collectionName, nil)
}
