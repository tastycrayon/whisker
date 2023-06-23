package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("2x7apdpw")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// add
		del_room := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "2x7apdpw",
			"name": "room",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "qn65wt47fm67idr",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": []
			}
		}`), del_room)
		collection.Schema.AddField(del_room)

		return dao.SaveCollection(collection)
	})
}
