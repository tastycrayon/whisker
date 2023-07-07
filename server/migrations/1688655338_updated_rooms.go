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

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yun28rcq",
			"name": "name",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 3,
				"max": 250,
				"pattern": "^[a-zA-Z0-9- ]+$"
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qn65wt47fm67idr")
		if err != nil {
			return err
		}

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yun28rcq",
			"name": "name",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 3,
				"max": 250,
				"pattern": "^[a-z0-9- ]+$"
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		return dao.SaveCollection(collection)
	})
}
