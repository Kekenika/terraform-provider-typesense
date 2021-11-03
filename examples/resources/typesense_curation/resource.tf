resource "typesense_curation" "my_curation" {
  name            = "my-curation"
  collection_name = typesense_collection.my_collection.name

  rule {
    query = "apple"
    match = "exact"
  }

  include {
    id = 4
  }


  include {
    id       = 10
    position = 100
  }

  exclude = {
    id = 100
  }
}
