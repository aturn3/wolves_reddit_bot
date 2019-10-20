#!/bin/bash

. gcp.config

gcloud functions call $testCloudFunction --data '{"topic":"$testCloudFunctionTopic","message":""}'
sleep 10s
gcloud functions logs read $testCloudFunction
