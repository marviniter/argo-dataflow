# Git Step

This intended as a convenient way to write steps without having to build and publish images.

When a steps starts, the code is checked out from Git, and then run using `./entrypoint.sh`.

```yaml
git:
  branch: main
  path: examples/git
  image: golang:1.16
  url: https://github.com/argoproj-labs/argo-dataflow
```

* [Example pipeline](examples/106-git-pipeline.yaml)
* [Source code](examples/git)