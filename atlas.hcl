// read internal/database/dialer/cockroachdb/README.md for more information

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/database/dialer/cockroachdb/table",
    "--dialect", "postgres"
  ]
}

// development environment
env "dev" {
  // pointing to gorm models to generate schemas
  src = data.external_schema.gorm.url
  // the database which is holding the actual tables
  url = "postgres://root:@localhost:26257/?sslmode=disable"
  // a temporary database for Atlas to do migrations
  dev = "postgres://root:@localhost:26258/?sslmode=disable"
  // location of migration files
  migration {
    dir = "file://internal/database/dialer/cockroachdb/migration"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

// production environment
env "prod" {

}
