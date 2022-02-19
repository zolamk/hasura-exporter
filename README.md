# Hasura Prometheus Exporter

Prometheus exporter for metrics about your [Hasura](https://hasura.io) GraphQL Server.

## Supported Hasura Versions

This exporter has been tested on hasura versions starting from `v2.0.0`

## Requirements

Your hasura graphql server needs to have the [metadata](https://hasura.io/docs/latest/graphql/core/api-reference/metadata-api/index.html) and [schema](https://hasura.io/docs/latest/graphql/core/api-reference/schema-api/index.html) APIs enabled, the APIs are enabled by default unless have been enabled

## Installation

For pre-built binaries please take a look at the [releases](https://github.com/zolamk/hasura-exporter)
## Configuration

| ENV Variable                | Description                                                               |
| --------------------------- | ------------------------------------------------------------------------- |
| DEBUG                       | If set to true also debug information will be logged, otherwise only info |
| HASURA_GRAPHQL_ADMIN_SECRET | Admin secret for hasura graphql server                                    |
| HASURA_GRAPHQL_ENDPOINT     | URL to hasura graphql server e.g `https://graph.example.com`              |
| WEB_ADDR                    | Address for this exporter to run, default: `:9921`                        |


## Metrics

| Name                               | Type  | Help                                                                                             |
| ---------------------------------- | ----- | ------------------------------------------------------------------------------------------------ |
| hasura_metadata_consistency_status | gauge | metadata consistency status of your hasura graphql server, 1 means metadata is consistent, 0 not |
| hasura_pending_event_triggers      | gauge | the number of pending event triggers labeld with trigger name                                    |
| hasura_processed_event_triggers    | gauge | the number of processed event triggers labeled with trigger name                                 |
| hasura_successful_event_triggers   | gauge | the number of successfully processed event triggers labeled with trigger name                    |
| hasura_failed_event_triggers       | gauge | the number of failed event triggers labeled with trigger name                                    |
| hasura_pending_cron_triggers       | gauge | the number of pending cron triggers labeled with trigger name                                    |
| hasura_processed_cron_triggers     | gauge | the number of processed cron triggers labeled with trigger name                                  |
| hasura_successful_cron_triggers    | gauge | the number of successfully processed event triggers labeled with trigger name                    |
| hasura_failed_cron_triggers        | gauge | the number of failed cron triggers labeled with trigger name                                     |
| hasura_pending_one_off_events      | gauge | the number of pending one off scheduled events                                                   |
| hasura_processed_one_off_events    | gauge | the number of processed one off scheduled events                                                 |
| hasura_successful_one_off_events   | gauge | the number of successfully processed on off scheduled events                                     |
| hasura_failed_one_off_events       | gauge | the number of failed one off scheduled events                                                    |