resource "typesense_collection_alias" "my_alias" {
  name            = "my-alias"
  collection_name = typesense_collection.my_collection.name
}
