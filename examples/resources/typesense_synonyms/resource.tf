resource "typesense_synonyms" "my_synonyms" {
  name            = "my-synonyms"
  collection_name = typesense_collection.my_collection.name

  synonyms = [
    "blazer",
    "coat",
    "jacket"
  ]
}
