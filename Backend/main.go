package main

import (
	"3DQuest/dbdriver"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// _ "github.com/go-kivik/couchdb/v4" // Kivik CouchDB Driver
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Couldn't load the .env file")
	}
}

func main() {
	fmt.Printf("Hello World @ PORT = %s\n", os.Getenv("PORT"))
	client, _ := dbdriver.CreateClient()
	// fmt.Println(client.URL.String())
	dbs, _ := dbdriver.GetAllDBs(client)
	db_info, _ := dbdriver.ConnectToDB(client, dbs[2])
	// fmt.Println(dbs)
	// db_info, _ := dbdriver.GetDBInfo(client, dbs[2])
	fmt.Println(db_info.Name, db_info.Sizes.External)
	// doc, _ := dbdriver.GetDesignDocument(client, dbs[2], "test")
	// fmt.Printf("%+v\n", doc.Views["admin_user_view"])
	// docInfo, _ := dbdriver.GetDesignDocumentInfo(client, dbs[2], "test")
	// fmt.Println(docInfo.ViewIndex.Sizes.File)

	// newDoc, _ := dbdriver.GetDocument(client, dbs[2], "b7a10978cdb6f2c3b3fbbc91660007a2")
	// fmt.Printf("%+v\n", newDoc["type"])
	opts := dbdriver.CreateDesignViewOptions()
	opts.IncludeDocs = true
	// @TODO must test json query
	view, _ := dbdriver.GetDesignView(client, "test", "all_user_view", opts)
	fmt.Printf("%+v\n", view.Rows[2].Key["name"])
	// fmt.Printf("THIS SHOULD NOT BE EMPTY %+v\n", view.Rows[1].Doc)
	opts.Sorted = false
	opts.IncludeDocs = false
	view, _ = dbdriver.GetDesignView(client, "test", "all_user_view", opts)
	// fmt.Printf("THIS SHOULD NOT EXIST %+v\n", view.TotalRows)
	// fmt.Println(opts.EndKey, opts.StartKeyDocID, opts.Update)
	new_id, _ := dbdriver.GetUUIDFromCouchDB(client)
	fmt.Println(new_id)
	// documentito := `{
	// 	"_id": "%s",
	// 	"type": "user",
	// 	"name": "Sr. Patata Potato",
	// 	"nif": "AAAAA",
	// 	"credits": 69,
	// 	"email": "mr.potato@patata.es",
	// 	"password_hash": "lmaolmaolmaohasherino"
	// }`
	// documentito = fmt.Sprintf(documentito, new_id)
	// fmt.Println(documentito)
	// data, _ := dbdriver.PutDocument(client, dbs[2], new_id, documentito)
	// fmt.Printf("%+v\n", *data)

	// fmt.Println("REV ", rev)
	// documentito := make(map[string]interface{}, 0)
	// documentito2 := &map[string]interface{}{
	// 	"type":          "admin",
	// 	"name":          "COOLIO The Artist",
	// 	"nif":           "676676676Q",
	// 	"credits":       777,
	// 	"email":         "gangstas@paradi.se",
	// 	"password_hash": "Everybody's runnin', but half of them ain't lookin'",
	// }
	// responserino, _ := dbdriver.CreateNewDocument(client, dbs[2], documentito2)
	// fmt.Printf("%+v\n", responserino)
	// res := dbdriver.CreateNewDatabase(client, "lmaodb")
	// fmt.Println(res)
	findopts := &dbdriver.FindOptions{}

	jsonStr := `{
			"credits": {
				"$lt": 666
			},
            "type": "user"
		}`

	findopts.Limit = 3
	json.Unmarshal([]byte(jsonStr), &findopts.Selector)
	fmt.Printf("%v\n", findopts.Selector)

	jsonfile, _ := json.Marshal(&findopts)
	fmt.Println(string(jsonfile))
	found, _ := dbdriver.FindInDatabase(client, findopts)
	fmt.Printf("%v\n", found.Docs)
}

func hdnl_hello_world(ectx echo.Context) error {
	return ectx.String(http.StatusOK, "hello world!")
}

// e := echo.New()
// 	e.GET("/", hdnl_hello_world)

// db := client.DB("questdb")

// usr := User{}
// row := db.Get(context.TODO(), "b7a10978cdb6f2c3b3fbbc91660007a2")
// ctx := context.TODO()
// models.QueryAllUsers(db, ctx)
// fmt.Println("All Users queried")

// models.QueryBasicUsers(db, ctx)
// fmt.Println("Admin Users queried")

// models.QueryUserByEmail(db, ctx)
// fmt.Println("Basic Users queried")

// e.Start(fmt.Sprintf(":%s", os.Getenv("PORT")))
