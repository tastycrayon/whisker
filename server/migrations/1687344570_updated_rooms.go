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
		collection.Schema.RemoveField("ccgsbd7u")

		// add
		new_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "niwunqw2",
			"name": "description",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
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
			"id": "ccgsbd7u",
			"name": "description",
			"type": "editor",
			"required": false,
			"unique": false,
			"options": {}
		}`), del_description)
		collection.Schema.AddField(del_description)

		// remove
		collection.Schema.RemoveField("niwunqw2")

		return dao.SaveCollection(collection)
	})
}
