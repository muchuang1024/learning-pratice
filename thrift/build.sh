#!bin/bash

if [ ! -d "/lib" ]; then
    mkdir -p lib
fi

thrift -gen go -out lib  service.thrift
thrift -gen go -out lib  request.thrift
thrift -gen go -out lib  types.thrift