package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepositry struct {
	collection *mongo.Collection
}

func NewProductRepository(connectionstring string) *ProductRepositry {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionstring))

	if err != nil {
		panic(err)
	}

	collection := client.Database("product_db").Collection("product_collection")

	fmt.Println("Connected Successfully")

	return &ProductRepositry{
		collection: collection,
	}
}

func (p *ProductRepositry) InsertProduct(product *models.Product) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.ID = primitive.NilObjectID
	result, err := p.collection.InsertOne(ctx, &product)
	if err != nil {
		return "", err
	}

	return result.InsertedID, nil
}
