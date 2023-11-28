package config

// Project const for compiled Boiler build with project config
const Project = "{\n  \"root\": {\n    \"files\": [\n      {\n        \"name\": \"readme\",\n        \"extension\": \"md\"\n      },\n      {\n        \"name\": \"Makefile\",\n        \"extension\": \"\",\n        \"template\": {\n          \"name\": \"makefile\"\n        }\n      },\n      {\n        \"name\": \"go\",\n        \"extension\": \"mod\",\n        \"template\": {\n          \"name\": \"gomod\"\n        }\n      },\n      {\n        \"name\": \"go\",\n        \"extension\": \"sum\"\n      }\n    ],\n    \"dirs\": [\n      {\n        \"name\": {\n          \"value\": \"api\"\n        }\n      },\n      {\n        \"name\": {\n          \"value\": \"cmd\"\n        },\n        \"dirs\": [\n          {\n            \"name\": {\n              \"value\": \"app\"\n            },\n            \"files\": [\n              {\n                \"name\": \"main\",\n                \"extension\": \"go\",\n                \"template\": {\n                  \"name\": \"main\"\n                }\n              }\n            ]\n          }\n        ]\n      },\n      {\n        \"name\": {\n          \"value\": \"internal\"\n        },\n        \"dirs\": [\n          {\n            \"name\": {\n              \"value\": \"app\"\n            }\n          },\n          {\n            \"name\": {\n              \"value\": \"pkg\"\n            }\n          }\n        ]\n      },\n      {\n        \"name\": {\n          \"value\": \"pkg\"\n        }\n      },\n      {\n        \"name\": {\n          \"value\": \".cicd\"\n        },\n        \"files\": [\n          {\n            \"name\": \"docker_compose\",\n            \"extension\": \"yaml\",\n            \"template\": {\n              \"name\": \"docker_compose\"\n            }\n          }\n        ]\n      }\n    ]\n  }\n}"

// Gqlgen const for compiled Boiler build with gqlgen project config
const Gqlgen = "{\n  \"root\": {\n    \"dirs\": [\n      {\n        \"name\": {\n          \"value\": \"api\"\n        },\n        \"dirs\": [\n          {\n            \"name\": {\n              \"value\": \"graphql\"\n            },\n            \"dirs\": [\n              {\n                \"name\": {\n                  \"value\": \"directives\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"directives\",\n                    \"extension\": \"graphql\"\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"mutation\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"mutation\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_mutation\"\n                    }\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"query\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"query\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_query\"\n                    }\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"scalars\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"scalars\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_scalars\"\n                    }\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"schema\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"schema\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_schema\"\n                    }\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"types\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"types\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_types\"\n                    }\n                  }\n                ]\n              },\n              {\n                \"name\": {\n                  \"value\": \"enums\"\n                },\n                \"files\": [\n                  {\n                    \"name\": \"enums\",\n                    \"extension\": \"graphql\",\n                    \"template\": {\n                      \"name\": \"gqlgen_enums\"\n                    }\n                  }\n                ]\n              }\n            ]\n          }\n        ]\n      }\n    ]\n  }\n}"
