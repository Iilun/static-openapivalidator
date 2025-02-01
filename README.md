# static-openapivalidator

## Idea

This project aims to validate API Request/Responses against an OpenAPI spec.

Oftentimes, validating an OpenAPI spec can be cumbersome. Tools such as [Prism](https://docs.stoplight.io/docs/prism) or others offer proxies to validate requests. However, including these in a CICD pipeline can make the CI environment stray away from production conditions.

API testing tooling such as Postman also offers OpenAPI validation, but it is difficult here to distinguish between functional failures and OpenAPI discrepancies. 

The goal here is to run the API test in a production-like environment, and then analyze the report files for OpenAPI discrepancies.

## Installation

```
go install github.com/Iilun/static-openapivalidator
```

## Usage

```
static-openapivalidator --spec <openapi file> --report <report containing API results>
```

### Options

| Flag         | Aliases | Description                                                      |
|--------------|---------|------------------------------------------------------------------|
| spec         | s       | Path to the openapi spec                                         |
| report       | r       | Path to the report containing the API results                    |
| format       | f       | [Format](#supported-formats) of the report file (default: bruno) |
| report-html  | -       | Path to output an HTML report to                                 |
| report-json  | -       | Path to output a JSON report to                                  |
| report-junit | -       | Path to output a JUNIT report to                                 |
| help         | h       | Display help information                                         |

### Configuration file

A configuration file path can be given through the `CONFIG_FILE` environment variable.

This YAML file allows for routes/requests/responses exclusion as well as some more advanced configuration options.

```yaml
ignore:
  # Ignore the servers section from the OpenAPI spec, allowing for hosts not given in the OpenAPI to match the route anyway
  servers: true
  # A list of requests name (can be extracted from html report) where the result should be ignored
  requests:
   - "glob for test name - depends on parser"
  # A list of responses name (can be extracted from html report) where the result should be ignored
  responses:
    - "glob for test name - depends on parser"
  # A list of relative routes where the result should be ignored
  routes:
   - "glob for relative path"
```


## Supported formats

Here are the formats available to use the results from

### Bruno

Flag value: `bruno`

Use with the file produced by the `--reporter-json` option of bru.

### Postman

Flag value: `postman`

Use with the file produced by the `-r json --reporter-json-export` option of newman.

## Troubleshooting

## Could not find route

If all your requests end up with the following error: 

    could not find route for <route>: no matching operation was found

Make sure to enable the option `ignore.servers: true` in the [configuration file](#configuration-file)