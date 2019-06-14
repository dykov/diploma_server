package util

import "time"

var (
	host     = "localhost"
	port     = 5432
	user     = "pygo_admin" //os.Getenv("PlacardDatabase_user")
	password = "12345"      // os.Getenv("PlacardDatabase_password")
	dbname   = "pygo_db"    //os.Getenv("PlacardDatabase_dbname")
)

const (
	AccessTokenExpiration  = 1 * time.Hour
	RefreshTokenExpiration = 24 * time.Hour
)

const (
	cacheDefaultExpiration = 30 * time.Minute
	cacheCleanupInterval   = 1 * time.Hour
)

const DurationToDelete = 30 * 24 * time.Hour

//var SecretKeyForToken = os.Getenv("secret_token_key")
var SecretKeyForToken = []byte("jwt_secret_key")

//var Salt2 = os.Getenv("password_salt2")
const Salt = "Salt"

const DatabaseConnectionFailed = "Database connection failed "
const DatabaseSessionFailed = "Database session failed "
const SuccessfulDatabaseConnection = "Successful database connection"
const NoDatabaseConnection = "No database connection "
const UnableToDecodeJson = "Unable to decode JSON in request body "
const UnableToFindObject = "Unable to find object in database "
const UnableToInsertObject = "Unable to insert object into database "
const UnableToUpdateObject = "Unable to update object in database "
const UnableToDeleteObject = "Unable to delete object in database "
const UnableToParsePageNumber = "Unable to parse page number "
const UnableToParseVersionValue = "Unable to parse version value"
const UnableToParseData = "Unable to parse data"
const SuccessfulArchiving = "Successful archiving "
const UnsuccessfulArchiving = "Unsuccessful archiving "
const SuccessfulDeleting = "Successful deleting "

const NotFound = "Not found "
const CityNotFound = "City not found "
const UserNotFound = "User not found "
const LocationNotFound = "Location not found "
const EventNotFound = "Event not found "
const CategoryNotFound = "Category not found "
const SomethingHappenedWrong = "Something happened wrong "

var CategoryTypeValue = map[string]uint64{
	"tag":       0,
	"event":     1,
	"location":  2,
	"organizer": 4,
}

const (
	NameMaxLength        = 100
	DescriptionMaxLength = 1500
	AddressMaxLength     = 100
)

/*
const (
	EntityUser          = iota // 0
	EntityEvent                // 1
	EntityPlace                // 2
	EntityEventsInPlace        // 3
	EntityComment              // 4
	EntitySettings             // 5
)
*/

const (
	GrantGuest = 2 << iota // 2
	GrantUser  = 2 << iota // 4
	GrantModer = 2 << iota // 8
	GrantAdmin = 2 << iota // 16
)
