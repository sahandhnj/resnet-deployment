#!/bin/bash

venv/bin/locust -f locust/locustfile.py --host=http://localhost:3001/api/model/imagedetector/v1