# assignment_demo_2023

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This is a demo and template for backend assignment of 2023 TikTok Tech Immersion.

## Installation

Requirement:

- golang 1.18+
- docker

To install dependency tools:

```bash
make pre
```

## Run

```bash
docker-compose up -d
```

Check if it's running:

```bash
curl localhost:8080/ping
```



#   Sean's mini report

This assignment was done mainly with the help of :
- https://o386706e92.larksuite.com/docx/QE9qdhCmsoiieAx6gWEuRxvWsRc


I had some trouble with configurations and kept getting errors while creating the docker containers using GoLand, while using docker-compose up -d I was able to create my containers.
Ultimately the Makefile helped by using 


```bash
make generate
```


While the containers are able to be created and messages are able to be sent, my Pull function does not work for some reason.

In the handler.go file I have recreated the getRoomID function using postman with this GET request

```bash
localhost:8080/api/pull?chat=s1:s2&cursor=0&limit=2&reverse=true
```
It seems I am getting a 500 error with my function returning the variable 'chat' as an empty string.
However I do not understand why it is so.

I am mindful that I was granted a small time extension to complete the assignment and I do not think it would be fair if I was given a later deadline than the 2 days already.

Therefore I will upload and submit this code as part of the assignment despite some parts not working.