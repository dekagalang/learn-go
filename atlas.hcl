env "local" {
  src = "file://schema.hcl"
  dev = "docker://postgres/16?search_path=public"
  migration {
    dir = "file://migrations"
    format = atlas
  }
  format {
    migrate {
      diff = "{{ sql . }}"
    }
  }
}

env "docker" {
  src = "file://schema.hcl"
  dev = "docker://postgres/16/crud_api_dev?user=postgres&password=postgres&search_path=public"
  migration {
    dir = "file://migrations"
    format = atlas
  }
}
