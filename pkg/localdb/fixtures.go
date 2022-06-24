package localdb

import (
	"context"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var FixtureSkus = []interface{}{
	StorageUnit{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		SKU:       "sku_1",
		Name:      "Storage Unit 1",
	},
	StorageUnit{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		SKU:       "sku_2",
		Name:      "Storage Unit 2",
	},
	StorageUnit{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		SKU:       "sku_3",
		Name:      "Storage Unit 3",
	},
}

var FixtureProducts = []interface{}{
	Product{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		Category:  "Shoe",
		Name:      "Shoe Red Spirit",
		Model:     "That one Marilyn Monroe had",
		Price:     6112.37,
		Brand:     "Unique Brands",
	},
	Product{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		Category:  "Shoe",
		Name:      "Shoe Blue Bird",
		Model:     "That one Peter Griffin had",
		Price:     92.37,
		Brand:     "Unique Brands",
	},
	Product{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		Category:  "Tools",
		Name:      "Screwdriver",
		Model:     "Tramontina Facility",
		Price:     2.70,
		Brand:     "Tramontina",
	},
	Product{
		ID:        primitive.NewObjectID(),
		Timestamp: primitive.Timestamp{},
		Category:  "Smartphones",
		Name:      "Xiaomi POCO x4 PRO",
		Model:     "Xiaomi POCO x4 PRO 6gb/128gb",
		Price:     690.88,
		Brand:     "XIAOMI",
	},
}

func InstallFixtures(client *mongo.Client) {
	skuCollection := client.Database("store").Collection("StorageUnit")
	skuCollection.InsertMany(context.TODO(), FixtureSkus)

	productCollection := client.Database("store").Collection("Product")
	productCollection.InsertMany(context.TODO(), FixtureProducts)

	cursorSku, err1 := skuCollection.Find(context.TODO(), bson.D{})
	cursorProduct, err2 := productCollection.Find(context.TODO(), bson.D{})

	var skus []SkuProduct
	if err1 = cursorSku.All(context.TODO(), &skus); err1 != nil {
		panic(err1)
	}

	var products []Product
	if err2 = cursorProduct.All(context.TODO(), &products); err2 != nil {
		panic(err2)
	}

	skuProdCollection := client.Database("store").Collection("SKUProducts")
	for _, sku := range skus {
		fmt.Println(sku.ID)
		for _, pro := range products {
			fmt.Println("\t", pro.ID)
			skuProducts := []interface{}{
				SkuProduct{
					ID:        primitive.NewObjectID(),
					Timestamp: primitive.Timestamp{},
					Sku:       string(sku.ID.Hex()),
					Product:   string(pro.ID.Hex()),
					Available: float64(rand.Intn(100)),
				},
			}
			skuProdCollection.InsertMany(context.TODO(), skuProducts)

		}
	}
}
