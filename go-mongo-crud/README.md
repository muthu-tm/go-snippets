# Go MongoDB CRUD #
Insert, Update, Delete & View in User details using "Golang - Gin - MongoDB"

## RUN
* Initialize the go dependency manager at first
> $ dep init

* Update and ensure the dependencies were added
> $ dep ensure --add go.mongodb.org/mongo-driver/mongo go.mongodb.org/mongo-driver/bson go.mongodb.org/mongo-driver/mongo/options

## Database details
please go through the .env file and create database and collection in mongodb

* To create database
> use test

* To create collectionn
> db.createCollection(collection_name)

## Date Formatter
The output is JSON with a date that has been added 12 hours like (Result) and it is available in *_/api/v2/date_*; you can use the simple ui (open the date-formatter-ajax.html) file

Request:
{
date: "2017-07-22T01:50"
}

Result:
{
date: "2017-07-22T13:50"
}

## Perlin Noise 2D

It has service to create perlin noise and it is available in *_/api/v3/perlin/2D_*; you can use the simple ui (open the perlin-ajax.html) file
* Ajax call to generate the Perlin Noise image. 
* Input 2 numbers x and y as image size (can use GET or POST request). 
* Then the system will calculate the calculation result from coordinate 0,0 up to x, y with result of calculation of Perlin noise in each coordinate.
* Use Thread (or) Goroutine to perform parallel calculations of each coordinate.
