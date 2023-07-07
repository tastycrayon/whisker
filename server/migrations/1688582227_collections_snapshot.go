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
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-06-01 09:32:51.602Z",
				"updated": "2023-07-04 15:15:04.243Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": 3,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": true
				}
			},
			{
				"id": "qn65wt47fm67idr",
				"created": "2023-06-03 19:21:02.317Z",
				"updated": "2023-06-28 17:29:36.914Z",
				"name": "rooms",
				"type": "base",
				"system": false,
				"schema": [
					{
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
					},
					{
						"system": false,
						"id": "2a4sho5z",
						"name": "slug",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": 3,
							"max": 250,
							"pattern": "^[a-z0-9-]+$"
						}
					},
					{
						"system": false,
						"id": "uj7x6ifn",
						"name": "cover",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 2097152,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/webp",
								"image/gif"
							],
							"thumbs": [
								"120x120"
							],
							"protected": false
						}
					},
					{
						"system": false,
						"id": "yc9o0kqf",
						"name": "createdBy",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "5csvygh5",
						"name": "description",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": 400,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xle4dlzn",
						"name": "type",
						"type": "select",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"public",
								"personal"
							]
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_nK5ZpLq` + "`" + ` ON ` + "`" + `rooms` + "`" + ` (` + "`" + `slug` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_KMebiiC` + "`" + ` ON ` + "`" + `rooms` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "@request.auth.id != \"\" && @request.data.createdBy = @request.auth.id",
				"updateRule": "@request.auth.id != \"\" && createdBy = @request.auth.id",
				"deleteRule": "@request.auth.id != \"\" && createdBy = @request.auth.id",
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
