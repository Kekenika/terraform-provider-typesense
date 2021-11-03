resource "typesense_collection" "my_collection" {
  name = "my-collection"

  fields {
    name = "name"
    type = "string"
  }
}
