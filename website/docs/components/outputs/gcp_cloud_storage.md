---
title: gcp_cloud_storage
type: output
status: experimental
categories: ["Services","GCP"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/output/gcp_cloud_storage.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

EXPERIMENTAL: This component is experimental and therefore subject to change or removal outside of major version releases.


Sends message parts as objects to a Google Cloud Storage bucket. Each object is
uploaded with the path specified with the `path` field.

Introduced in version 3.43.0.


<Tabs defaultValue="common" values={[
  { label: 'Common', value: 'common', },
  { label: 'Advanced', value: 'advanced', },
]}>

<TabItem value="common">

```yaml
# Common config fields, showing default values
output:
  gcp_cloud_storage:
    bucket: ""
    path: ${!count("files")}-${!timestamp_unix_nano()}.txt
    content_type: application/octet-stream
    max_in_flight: 1
    batching:
      count: 0
      byte_size: 0
      period: ""
      check: ""
```

</TabItem>
<TabItem value="advanced">

```yaml
# All config fields, showing default values
output:
  gcp_cloud_storage:
    bucket: ""
    path: ${!count("files")}-${!timestamp_unix_nano()}.txt
    content_type: application/octet-stream
    content_encoding: ""
    chunk_size: 16777216
    max_in_flight: 1
    timeout: 5s
    batching:
      count: 0
      byte_size: 0
      period: ""
      check: ""
      processors: []
```

</TabItem>
</Tabs>

In order to have a different path for each object you should use function
interpolations described [here](/docs/configuration/interpolation#bloblang-queries), which are
calculated per message of a batch.

### Metadata

Metadata fields on messages will be sent as headers, in order to mutate these values (or remove them) check out the [metadata docs](/docs/configuration/metadata).

### Credentials

By default Benthos will use a shared credentials file when connecting to GCP
services. You can find out more [in this document](/docs/guides/gcp).

### Batching

It's common to want to upload messages to Google Cloud Storage as batched
archives, the easiest way to do this is to batch your messages at the output
level and join the batch of messages with an
[`archive`](/docs/components/processors/archive) and/or
[`compress`](/docs/components/processors/compress) processor.

For example, if we wished to upload messages as a .tar.gz archive of documents
we could achieve that with the following config:

```yaml
output:
  gcp_cloud_storage:
    bucket: TODO
    path: ${!count("files")}-${!timestamp_unix_nano()}.tar.gz
    batching:
      count: 100
      period: 10s
      processors:
        - archive:
            format: tar
        - compress:
            algorithm: gzip
```

Alternatively, if we wished to upload JSON documents as a single large document
containing an array of objects we can do that with:

```yaml
output:
  gcp_cloud_storage:
    bucket: TODO
    path: ${!count("files")}-${!timestamp_unix_nano()}.json
    batching:
      count: 100
      processors:
        - archive:
            format: json_array
```

## Performance

This output benefits from sending multiple messages in flight in parallel for
improved performance. You can tune the max number of in flight messages with the
field `max_in_flight`.

This output benefits from sending messages as a batch for improved performance.
Batches can be formed at both the input and output level. You can find out more
[in this doc](/docs/configuration/batching).

## Fields

### `bucket`

The bucket to upload messages to.


Type: `string`  
Default: `""`  

### `path`

The path of each message to upload.
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `"${!count(\"files\")}-${!timestamp_unix_nano()}.txt"`  

```yaml
# Examples

path: ${!count("files")}-${!timestamp_unix_nano()}.txt

path: ${!meta("kafka_key")}.json

path: ${!json("doc.namespace")}/${!json("doc.id")}.json
```

### `content_type`

The content type to set for each object.
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `"application/octet-stream"`  

### `content_encoding`

An optional content encoding to set for each object.
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `""`  

### `chunk_size`

An optional chunk size which controls the maximum number of bytes of the object that the Writer will attempt to send to the server in a single request. If ChunkSize is set to zero, chunking will be disabled.


Type: `number`  
Default: `16777216`  

### `max_in_flight`

The maximum number of messages to have in flight at a given time. Increase this to improve throughput.


Type: `number`  
Default: `1`  

### `timeout`

The maximum period to wait on an upload before abandoning it and reattempting.


Type: `string`  
Default: `"5s"`  

### `batching`

Allows you to configure a [batching policy](/docs/configuration/batching).


Type: `object`  

```yaml
# Examples

batching:
  byte_size: 5000
  count: 0
  period: 1s

batching:
  count: 10
  period: 1s

batching:
  check: this.contains("END BATCH")
  count: 0
  period: 1m
```

### `batching.count`

A number of messages at which the batch should be flushed. If `0` disables count based batching.


Type: `number`  
Default: `0`  

### `batching.byte_size`

An amount of bytes at which the batch should be flushed. If `0` disables size based batching.


Type: `number`  
Default: `0`  

### `batching.period`

A period in which an incomplete batch should be flushed regardless of its size.


Type: `string`  
Default: `""`  

```yaml
# Examples

period: 1s

period: 1m

period: 500ms
```

### `batching.check`

A [Bloblang query](/docs/guides/bloblang/about/) that should return a boolean value indicating whether a message should end a batch.


Type: `string`  
Default: `""`  

```yaml
# Examples

check: this.type == "end_of_transaction"
```

### `batching.processors`

A list of [processors](/docs/components/processors/about) to apply to a batch as it is flushed. This allows you to aggregate and archive the batch however you see fit. Please note that all resulting messages are flushed as a single batch, therefore splitting the batch into smaller batches using these processors is a no-op.


Type: `array`  
Default: `[]`  

```yaml
# Examples

processors:
  - archive:
      format: lines

processors:
  - archive:
      format: json_array

processors:
  - merge_json: {}
```

