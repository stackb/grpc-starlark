"""repositories.bzl declares dependencies for the workspace
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def repositories():
    """repositories loads all dependencies for the workspace
    """
    rules_proto()  # via <TOP>
    io_bazel_rules_go()  # via bazel_gazelle
    bazel_gazelle()  # via <TOP>
    build_stack_rules_proto()
    protobuf_core_deps()
    io_bazel_rules_docker()
    hermetic_cc_toolchain()

def protobuf_core_deps():
    bazel_skylib()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via <TOP>

def io_bazel_rules_go():
    # Release: v0.39.1
    # TargetCommitish: release-0.39
    # Date: 2023-04-20 04:35:08 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_go/releases/tag/v0.39.1
    # Size: 1759832 (1.8 MB)
    _maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "473a064d502e89d11c497a59f9717d1846e01515a3210bd169f22323161c076e",
        strip_prefix = "rules_go-0.39.1",
        urls = ["https://github.com/bazelbuild/rules_go/archive/v0.39.1.tar.gz"],
    )

def bazel_gazelle():
    # Branch: master
    # Commit: 2d1002926dd160e4c787c1b7ecc60fb7d39b97dc
    # Date: 2022-11-14 04:43:02 +0000 UTC
    # URL: https://github.com/bazelbuild/bazel-gazelle/commit/2d1002926dd160e4c787c1b7ecc60fb7d39b97dc
    #
    # fix updateStmt makeslice panic (#1371)
    # Size: 1859745 (1.9 MB)
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "5ebc984c7be67a317175a9527ea1fb027c67f0b57bb0c990bac348186195f1ba",
        strip_prefix = "bazel-gazelle-2d1002926dd160e4c787c1b7ecc60fb7d39b97dc",
        urls = ["https://github.com/bazelbuild/bazel-gazelle/archive/2d1002926dd160e4c787c1b7ecc60fb7d39b97dc.tar.gz"],
    )

def local_bazel_gazelle():
    _maybe(
        native.local_repository,
        name = "bazel_gazelle",
        path = "/Users/i868039/go/src/github.com/bazelbuild/bazel-gazelle",
    )

def rules_proto():
    # Commit: f7a30f6f80006b591fa7c437fe5a951eb10bcbcf
    # Date: 2021-02-09 14:25:06 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_proto/commit/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf
    #
    # Merge pull request #77 from Yannic/proto_descriptor_set_rule
    #
    # Create proto_descriptor_set
    # Size: 14397 (14 kB)
    _maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
        strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
        urls = ["https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz"],
    )

def build_stack_rules_proto():
    # Branch: master
    # Commit: aa380e4421057b35228544bc234f816bb6b72c1c
    # Date: 2022-12-08 05:19:32 +0000 UTC
    # URL: https://github.com/stackb/rules_proto/commit/aa380e4421057b35228544bc234f816bb6b72c1c
    #
    # use distinct impLang for scala proto exports (#304)
    #
    # * use distinct impLang for scala proto exports
    # * fix test
    # Size: 2074364 (2.1 MB)
    http_archive(
        name = "build_stack_rules_proto",
        sha256 = "820dc71f2e265a50104671d323caba53790dfe20e9f7249a0e6beeaee39b4597",
        strip_prefix = "rules_proto-aa380e4421057b35228544bc234f816bb6b72c1c",
        urls = ["https://github.com/stackb/rules_proto/archive/aa380e4421057b35228544bc234f816bb6b72c1c.tar.gz"],
    )

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz",
        ],
    )

def rules_python():
    _maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz",
        ],
    )

def zlib():
    _maybe(
        http_archive,
        name = "zlib",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        build_file = "@build_stack_rules_proto//third_party:zlib.BUILD",
    )

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
        strip_prefix = "protobuf-3.14.0",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        ],
    )

def io_bazel_rules_docker():
    # Branch: master
    # Commit: 8e70c6bcb584a15a8fd061ea489b933c0ff344ca
    # Date: 2023-04-27 20:06:36 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_docker/commit/8e70c6bcb584a15a8fd061ea489b933c0ff344ca
    #
    # The OCI distribution spec only allows lower case letter in container repository (#2252)
    #
    # name
    # (https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests),
    # this doesn't bode well with Bazel package path using upper case.
    #
    # To be clear: docker itself is ok with upper case:
    # https://ktomk.github.io/pipelines/doc/DOCKER-NAME-TAG.html.
    #
    # Picking the common denominator here, and force lower case on pkg path.
    # Size: 601209 (601 kB)
    _maybe(
        http_archive,
        name = "io_bazel_rules_docker",
        sha256 = "c27b53d53a5704fb676078843f1a674ff196ab4fb9d7f6b74cf7748b47c9374f",
        strip_prefix = "rules_docker-8e70c6bcb584a15a8fd061ea489b933c0ff344ca",
        urls = ["https://github.com/bazelbuild/rules_docker/archive/8e70c6bcb584a15a8fd061ea489b933c0ff344ca.tar.gz"],
    )

def hermetic_cc_toolchain():
    # Commit: a9d87b21a5dddd691336c6c0004fa5dcfe5b9b48
    # Date: 2023-06-05 07:01:43 +0000 UTC
    # URL: https://github.com/uber/hermetic_cc_toolchain/commit/a9d87b21a5dddd691336c6c0004fa5dcfe5b9b48
    #
    # UBSAN: strip paths
    #
    # When a sanitizer (say, UBSAN via `-fsanitize=undefined`) is turned on, `zig cc` will include absolute paths to some header files, making the artifacts non-reproducible.
    #
    # This commit strips dirnames from such header files.
    #
    # Signed-off-by: Motiejus Jakštys <motiejus@uber.com>
    # Size: 45271 (45 kB)
    _maybe(
        http_archive,
        name = "hermetic_cc_toolchain",
        sha256 = "92f42183aaaa4c05610f4a6b37e30d54bd52020ec34144fbc2a5cabc02656612",
        strip_prefix = "hermetic_cc_toolchain-a9d87b21a5dddd691336c6c0004fa5dcfe5b9b48",
        urls = ["https://github.com/uber/hermetic_cc_toolchain/archive/a9d87b21a5dddd691336c6c0004fa5dcfe5b9b48.tar.gz"],
    )
