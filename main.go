package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstance wraps the MongoDB client and database connection.
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"                            // Database name
const mongoURI = "mongodb://localhost:27017/" + dbName // MongoDB connection URI

// Employee struct represents an employee entity.
type Employee struct {
	Name   string  `json:"name"`                              // Employee name
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"` // Employee ID
	Salary float64 `json:"salary"`                            // Employee salary
	Age    float64 `json:"age"`                               // Employee age
}

// Connect initializes the MongoDB connection.
func Connect() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     client.Database(dbName),
	}

	return nil
}

func main() {
	// Connect to the MongoDB database.
	if err := Connect(); err != nil {
		log.Fatal("Error while connecting to DB : %v ", err)
	}

	// Create a new Fiber app instance.
	app := fiber.New()

	// Get all employees.
	app.Get("/api/v1/employee", func(c *fiber.Ctx) error {
		query := bson.D{{}} // Empty query to fetch all documents

		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employees []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees) // Return the list of employees in JSON format
	})

	// Add a new employee.
	app.Post("/api/v1/employee", func(c *fiber.Ctx) error {
		employee := new(Employee)

		// Parse the request body to get employee details.
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee.ID = "" // Reset the ID to let MongoDB generate a new one.

		// Insert the new employee record into the database.
		insertionResult, err := mg.Db.Collection("employees").InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee)

		return c.JSON(createdEmployee) // Return the created employee record
	})

	// Update an existing employee by ID.
	app.Put("/api/v1/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id") // Get the employee ID from the request parameters

		employeeID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee := new(Employee)

		// Parse the request body to get updated employee details.
		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeID}} // Query to find the employee by ID
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "salary", Value: employee.Salary},
					{Key: "age", Value: employee.Age},
				},
			},
		}

		// Update the employee record in the database.
		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404) // Return 404 if no record was found
			}
			return c.SendStatus(500)
		}

		employee.ID = idParam // Set the employee ID in the response

		return c.Status(200).JSON(employee) // Return the updated employee record
	})

	// Delete an employee by ID.
	app.Delete("/api/v1/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id") // Get the employee ID from the request parameters

		employeeID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeID}} // Query to find the employee by ID

		// Delete the employee record from the database.
		result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), &query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(400) // Return 400 if no record was deleted
		}

		return c.Status(200).JSON("Record deleted") // Confirm the deletion
	})

	// Start the Fiber application on port 3000.
	log.Fatal(app.Listen(":3000"))
}
