package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeHandler struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewRecipeHandler(ctx context.Context, collection *mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		ctx:        ctx,
		collection: collection,
	}
}
