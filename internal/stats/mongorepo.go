package stats

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const COLLECTION = "verifications"

type MongoRepository struct {
	db         *mongo.Database
	collection string
}

func marshall(verification DNAVerification) bson.M {
	N := len(verification.DNA)
	var b strings.Builder
	b.Grow(N*N)
	for _, s := range verification.DNA {
		b.WriteString(s)
	}

	return bson.M{
		"dna": b.String(),
		"result": verification.Result,
	}
}

func NewMongoRepository(host, dbName string) (MongoRepository, error) {
	// Setup database instance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return MongoRepository{}, err
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return MongoRepository{}, err
	}
	db := client.Database(dbName)
	collection := COLLECTION

	// Create index on dna to ensure unique dna
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = db.Collection(collection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"dna": 1},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return MongoRepository{}, err
	}
	return MongoRepository {
		db:         db,
		collection: collection,
	}, nil
}

func (repo MongoRepository) Persist(verification DNAVerification) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := repo.db.Collection(repo.collection).InsertOne(ctx, marshall(verification))
	if !isDupKeyErr(err) {
		return err
	}

	return nil
}

func (repo MongoRepository) CountByResult(value bool) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count, err := repo.db.Collection(repo.collection).CountDocuments(ctx, bson.M{
		"result": value,
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

// isDupKeyErr checks if an error is a duplicate key error
func isDupKeyErr(err error) bool {
	we, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}
	for _, e := range we.WriteErrors {
		if e.Code == 11000 {
			return true
		}
	}
	return false
}