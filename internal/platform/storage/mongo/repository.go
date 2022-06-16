package mongo

import (
	"context"
	"log"

	"github.com/patriciabonaldy/sports-news/cmd/bootstrap/config"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/patriciabonaldy/sports-news/internal"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

const collectionName = "article"

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

func (r *Repository) GetArticles(ctx context.Context) ([]internal.ArticleNews, error) {
	cursor, err := r.getCollection(collectionName).Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var results []internal.ArticleNews
	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem ArticleNews
		err = cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, parseToBusinessArticleNews(elem))
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cursor.Close(ctx)

	return results, nil
}

func (r *Repository) GetArticleByID(ctx context.Context, articleID string) (*internal.ArticleNews, error) {
	var result ArticleNews

	err := r.getCollection(collectionName).
		FindOne(ctx, bson.M{"article_id": articleID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal.ErrArticleNotFound
		}

		return nil, err
	}

	art := parseToBusinessArticleNews(result)
	return &art, nil
}

func (r *Repository) Save(ctx context.Context, article internal.ArticleNews) error {
	articleDB := parseToArticleNewsDB(article)
	_, err := r.getCollection(collectionName).InsertOne(ctx, articleDB)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) getCollection(collectionName string) *mongo.Collection {
	return r.db.Database(r.databaseName).Collection(collectionName, nil)
}
