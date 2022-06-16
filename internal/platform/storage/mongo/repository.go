package mongo

import (
	"context"
	"log"

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

func (r *Repository) GetArticles(ctx context.Context) ([]internal.ArticleNews, error) {
	opts := options.Find()
	cursor, err := r.getCollection(eventCollectionName).Find(context.TODO(), bson.M{}, opts)
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

func (r *Repository) GetArticleByID(ctx context.Context, articleID string) (internal.ArticleNews, error) {
	var result ArticleNews

	err := r.getCollection(eventCollectionName).
		FindOne(ctx, bson.M{"article_id": articleID}).Decode(&result)
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
