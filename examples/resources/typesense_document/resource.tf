resource "typesense_document" "doc" {
  id              = "doc"
  collection_name = typesense_collection.my_collection.name

  document = {
    "company_name"  = "Stark Industries"
    "num_employees" = 5215
    "country"       = "USA"
  }
}
