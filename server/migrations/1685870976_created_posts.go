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
			"id": "eegt3gpod8aufp7",
			"created": "2023-06-04 09:29:36.319Z",
			"updated": "2023-06-04 09:29:36.319Z",
			"name": "posts",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "dtxq7aeo",
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
					"id": "4klmexnc",
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
					"id": "x0jizgbp",
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

		collection, err := dao.FindCollectionByNameOrId("eegt3gpod8aufp7")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
