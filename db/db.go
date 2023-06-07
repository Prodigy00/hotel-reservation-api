package db

const (
	DbName     = "hotel-reservation"
	DbUri      = "mongodb://localhost:27017"
	TestDbName = "hotel-reservation-test"
)

//func ToObjectID(id string) primitive.ObjectID {
//	oid, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		panic(err)
//	}
//	return oid
//}
