# calc_s3_mp_etag
Standalone utility to calculate AWS multi-part uploaded S3 file Etags

Inspired from in particular r03 & Antonio Espinosa's StackOverflow answers to:
https://stackoverflow.com/questions/6591047/etag-definition-changed-in-amazon-s3

Example invocations:
```
$ touch 0
$ calc_s3_mp_etag 0
d41d8cd98f00b204e9800998ecf8427e-0

$ echo a > myfn
$ calc_s3_mp_etag -fn=1 myfn
myfn: 6d9d51f5ac5484b9001c319dbb39c139-1

$ dd if=/dev/zero bs=1048576 count=10 of=10M
$ calc_s3_mp_etag 10M
669fdad9e309b552f1e9cf7b489c1f73-2
$ calc_s3_mp_etag -chunksize=15 10M
9fbaeee0ccc66f9a8e3d3641dca37281-1
```

### Background
I wanted to validate before cleaning up the local sources of archiving done using AWS's StorageGateway utility.

As per the StackPeople have figured out AWS calculates S3 file Etag metadata by MD5 summing the MD5 sums of each multi-part part uploaded.  There is an existing Bash script + other scripts as answers to the question, but I decided I wanted a compiled program, and thought this may come in useful for others.  2 problems I encountered with the Go function posted by r03 included:
1. it assumed files under the chunk size were uploaded non-multipart - StorageGateway doesn't appear to do this and uses multi-part for everything (with the [default 8MB](https://docs.aws.amazon.com/cli/latest/topic/s3-config.html#multipart-chunksize) chunk size)
2. it reads the whole file into memory and can crash due to this
