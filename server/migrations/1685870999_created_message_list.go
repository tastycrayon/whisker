package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "5qxpiq26zh75xfa",
			"created": "2023-06-04 09:29:59.334Z",
			"updated": "2023-06-04 09:29:59.334Z",
			"name": "message_list",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "mdljvabr",
					"name": "roomId",
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
				},
				{
					"system": false,
					"id": "iyl2emvf",
					"name": "sender",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "p8wnkhw9",
					"name": "content",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("5qxpiq26zh75xfa")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
