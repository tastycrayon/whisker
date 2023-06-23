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

		collection, err := dao.FindCollectionByNameOrId("qn65wt47fm67idr")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("dwsojqdg")

		// add
		new_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ccgsbd7u",
			"name": "description",
			"type": "editor",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_description)
		collection.Schema.AddField(new_description)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qn65wt47fm67idr")
		if err != nil {
			return err
		}

		// add
		del_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "dwsojqdg",
			"name": "description",
			"type": "editor",
			"required": false,
			"unique": false,
			"options": {}
		}`), del_description)
		collection.Schema.AddField(del_description)

		// remove
		collection.Schema.RemoveField("ccgsbd7u")

		return dao.SaveCollection(collection)
	})
}
