<pre>                                             
    )                                            )  
 ( /(     )        (          (  (     (      ( /(  
 )\()) ( /(  `  )  )\   (     )\))(   ))\ (   )\()) 
((_)\  )(_)) /(/( ((_)  )\ ) ((_))\  /((_))\ (_))/  
| |(_)((_)_ ((_)_\ (_) _(_/(  (()(_)(_)) ((_)| |_   
| ' \ / _` || '_ \)| || ' \))/ _` | / -_)(_-<|  _|  
|_||_|\__,_|| .__/ |_||_||_| \__, | \___|/__/ \__|  
            |_|              |___/                  
</pre>

# go-hapingest

[![Go Report Card](https://goreportcard.com/badge/github.com/Phillezi/kthcloud-cli?style=social)](https://goreportcard.com/report/github.com/CM1007-FHIR-DATA/go-hapingest)

## Table of Contents

- [Overview](#overview)
- [Configuration](#configuration)

## Overview

`go-hapingest` is a specialized file server designed to simplify bulk importing of FHIR data into a HAPI server. It makes sure that data is bulk imported in the correct order by creating multiple manifests based on what the resources depends on. It is only made to handle ndjson files, which are commonly used in FHIR bulk imports.

You can run `go-hapingest` as a Docker container, mounting your FHIR data files to the `/fhir-data` directory within the container. This approach simplifies access to large datasets and facilitates their import directly into the HAPI server by letting the HAPI server fetch the content from the fileserver.

To enable bulk importing on your HAPI server, make sure to set the `HAPI_FHIR_BULK_IMPORT_ENABLED` flag to `true`. This configuration is essential for allowing the server to process bulk import requests properly.

## Configuration

`go-hapingest` can be configured using environment variables. Below is a table of configurable environment variables:

| Environment Variable        | Description                                                                          | Default Value                           |
| --------------------------- | ------------------------------------------------------------------------------------ | --------------------------------------- |
| `DATA_DIR`                  | Directory for storing FHIR data files                                                | `./fhir-data`                           |
| `URL_BASE`                  | Base URL for the application                                                         | `http://host.docker.internal`           |
| `FHIR_SERVER_URL`           | URL for the FHIR server                                                              | `http://host.docker.internal:8080/fhir` |
| `PORT`                      | Port on which the fileserver runs                                                    | `8001`                                  |
| `BLOCKING_PING_FHIR_SERVER` | Blocks the main thread until the FHIR server responds with an ok code on /fhir/$meta | `true`                                  |
| `SLEEP_WHEN_DONE`           | Sleeps infinitly when done                                                           | `true`                                  |
