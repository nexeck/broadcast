load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_library", "go_test")

gazelle(
    name = "gazelle",
    external = "vendored",
    prefix = "github.com/nexeck/multicast",
)

go_library(
    name = "go_default_library",
    srcs = ["multicast.go"],
    importpath = "github.com/nexeck/multicast",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["multicast_test.go"],
    importpath = "github.com/nexeck/multicast",
    library = ":go_default_library",
)
