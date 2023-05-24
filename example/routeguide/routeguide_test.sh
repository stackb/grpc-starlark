#!/bin/bash

set -euo pipefail

server_pid=0
server_address="localhost:1234"
server_address_file=$(mktemp "${TEST_TMPDIR}/tmp.stderr.XXXXXX")
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
    "${server_exe}" \
        --port=0 \
        --bind_address_file="${server_address_file}" \
        --load="${handlers_file}" \
        --protoset="${protoset_file}" \        &
    server_pid=$!
}

function get_server_address {
    server_address=$(cat "${server_address_file}")
    if [[ "${server_address}" ]]
    then
        echo "server_address: ${server_address}"
    else
        echo "failed to get server address"
        exit 1
    fi
}

function stop_server {
    kill $server_pid
}

function main {
    start_server
    sleep 0.1
    get_server_address

    echo "========================================"
    test_get_feature
    echo "========================================"
    test_list_features
    echo "========================================"
    test_record_route
    echo "========================================"
    test_route_chat

    stop_server
}

main
