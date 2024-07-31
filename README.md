# MAE (Mirror Accelerated Engine)
This is a Docker-based image acceleration engine, which is mainly used to accelerate the download of Docker hub images in K8S

## Features
- It only requires you to deploy a Pod and a strategy to achieve acceleration
- It can help us quickly solve the recent problem that Docker Hub cannot pull images

## How to use?
- Deploy the service to a K8S network communication environment
- Deploy the corresponding WebHook strategy in K8S

## Notice
- This project is still in the development stage, and the current version is only for internal testing

## Project Instruction
- This project involves Webhook intercepting data from Kubernetes API Server
- So I hope you can understand the code before considering using it