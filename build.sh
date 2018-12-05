#!/bin/bash 

cd v3
go build -o server
mv server ../server
cd ..

rm meta/input/*
./server