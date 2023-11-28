#! /bin/bash

docker image build -f Dockerfile -t ascii-art-dockerize .

docker container run -p 8080:8080 --detach --name ascii-dockerize-container ascii-art-dockerize

powershell.exe /c start http://localhost:8080
