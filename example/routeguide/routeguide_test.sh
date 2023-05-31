#!/bin/bash

set -euo pipefail

server_pid=0
server_address="localhost:3444"
server_exe="./cmd/grpc-starlark/grpc-starlark_/grpc-starlark"
protoset_file="./example/routeguide/routeguide_proto_descriptor.pb"
handlers_file="./example/routeguide/routeguide.grpc.star"
grpcurl_exe="./example/routeguide/grpcurl.exe"

function test_get_feature {
    "${grpcurl_exe}" \
        -vv \
        -d '{}' \
        -plaintext \
        -protoset "${protoset_file}" \
        "${server_address}" \
        'example.routeguide.RouteGuide.GetFeature'
    exit_code=$?
    test $exit_code -eq 0 && echo "PASS" || echo "FAIL"
    return $exit_code
}

function test_list_features {
    "${grpcurl_exe}" \
        -vv \
        -d '{}' \
        -plaintext \
        -protoset "${protoset_file}" \
        "${server_address}" \
        'example.routeguide.RouteGuide.ListFeatures'
    exit_code=$?
    test $exit_code -eq 0 && echo "PASS" || echo "FAIL"
    return $exit_code
}

function test_record_route {
    "${grpcurl_exe}" \
        -vv \
        -d @ \
        -plaintext \
        -protoset "${protoset_file}" \
        "${server_address}" \
        'example.routeguide.RouteGuide.RecordRoute' \
        <<-EOM
{ "latitude": 413628156, "longitude": -749015468 }
{ "latitude": 413628156, "longitude": -749015468 }
{ "latitude": 413628156, "longitude": -749015468 }
{ "latitude": 413628156, "longitude": -749015468 }    
EOM
    exit_code=$?
    test $exit_code -eq 0 && echo "PASS" || echo "FAIL"
    return $exit_code
}

function test_route_chat {
    "${grpcurl_exe}" \
        -vv \
        -d @ \
        -plaintext \
        -protoset "${protoset_file}" \
        "${server_address}" \
        'example.routeguide.RouteGuide.RouteChat' \
        <<-EOM
{ "message": "First" }
{ "message": "Second" }
{ "message": "Third" }
EOM
    exit_code=$?
    test $exit_code -eq 0 && echo "PASS" || echo "FAIL"
    return $exit_code
}

function start_server {
    ROUTEGUIDE_ADDRESS="${server_address}" \
    "${server_exe}" \
        --protoset="${protoset_file}" \
        "${handlers_file}" \
        &
    server_pid=$!
}

function stop_server {
    kill $server_pid
}

function main {
    start_server

    echo "========================================"
    test_get_feature
    # echo "========================================"
    # test_list_features
    # echo "========================================"
    # test_record_route
    # echo "========================================"
    # test_route_chat
    exit 1
    stop_server
}

main
