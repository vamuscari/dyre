# \[Dy\]namic \[Re\]quests

sql query match

DyRe is a request builder for a middleware  service. The intent of DyRe is to help automate selections of multiple fields without removing base functionality. If there is the desire to build additional functions on your request or handle things in a specific way, then it should more or less be business a usual when using Gin.    



## Setting up JSON config

First you will need json file to build with.
The name of the request is the thing being called.
Fields can be a string or an object with multiple params. 
If a type is not declared make sure to setup a default type that works for you. 
Groups are collections of fields to help simplify requests. A group must have a name and fields. 
The fields in a group are the same as in the base fields.

```json
[
  {
    "name": "Customers",
    "tableName": "Customers",
    "fields": [
      {
        "name": "CustomerID",
        "required": true
      },
      "Name",
      "Phone"
    ],
    "groups": [
      {
        "name": "Address",
        "fields": [
          {
            "name": "Street"
          },
          {
            "name": "City"
          },
          {
            "name": "State",
            "type": "null.Int"
          },
          {
            "name": "Zipcode",
            "type": "null.String"
          },
          {
            "name": "Zipint",
            "type": "null.Int"
          }
        ]
      }
    ]
  }
]

```


## Setting up middleware
```go

// Global Var for fetching request info
var Re map[string]dyre.DyRe_Request

// "gopkg.in/guregu/null.v4"
func main() {

	var sqlTypes = map[string]interface{}{
		"null.String": null.String{},
		"null.Int":    null.Int{},
		"null.Bool":   null.Bool{},
		"null.Time":   null.Time{},
		"null.Float":  null.Float{},
	}
    // Add the sql Types to the know types
    // The name does not have to match
	dyre.AddTypes(sqlTypes)
	dyre.DefaultType = "null.String"

	var dyre_err error
	Re, dyre_err = dyre.Init("./dyre.json")
	if dyre_err != nil {
		log.Panicf("dyre init failed: %v", dyre_err)
	}
}
```
