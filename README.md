# ASSIGNMENT

Write a Golang based HTTP server which accepts GET requests with input parameter as “sortKey”
and “limit”. The server queries three URLs mentioned below, combines the results from all three
URLs, sorts them by the sortKey and returns the response. The Server should also limit the number
of items in the API response to input parameter “limit”.

There are many ways by which you can run this HTTP server.

## One way
1. If you have kubernetes locally on you machine, just deploy the deployment using command:
> kubectl apply -f deployment.yml

2. Now you can acces the api at http://localhost:80/getData?sortKey=views&limit=10

3. For cleanup:
> kubectl delete -f deployment.yml

## Second way
1. If you have docker, you can build the Dockerfile using command:
> docker build -t <name_you_wish> . 

2. Run the docker container
> docker run -it -p 80:8000 <name_of_image_in_step_1>

3. Now you can access the api at http://localhost:80/getData?sortKey=views&limit=10

## Third way
1. If you have nothing you can just run the executable:
> ./assignment

2. Or else if you have Go installed, you can run the code:
> go run main.go 

3. Now you can access the api at http://localhost:8000/getData?sortKey=views&limit=10