# Insider trades
Service that provides recent and structured insider trades information.

## Documentation
- [trades-receiver-service](trades-receiver-service/README.md)


## Getting started

1. Run ``bazel run //:gazelle`` to update dependencies
2. Run ``bazel build //...`` to build the whole project


You can run tests using:

```bazel test //...```

Run this command to build dependencies from go.mod:

``bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies``
