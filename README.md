# poc-csv-uploader
Application to upload a csv file and persist with another service
This repository is nothing to be really used in production. Is something just to demonstrate my Golang Skills
This intended to accept large files. used with [this repository](https://github.com/VictorPrado99/poc-csv-persistence) to persist the data in a database.  
Is possible to send a csv file, and wait the processing happen. The process is asynchronous.  

### `POST` /upload
This endpoint receive as body a csv file with key called myCsv. The columns of this csv should be as follow

| id | email  | phone_number  | parcel_weight  |
|----|--------|---------------|----------------|

As you can see in ./dev

The process is async, so unless you have some problem in the header. You will receive 202.

## Architecture

Whis microservice is pretty simple, he will work multithreaded to parse teh csv into objects, after that, will transform into a array, and this array into a json. This json is sent to [poc-csv-persistence](https://github.com/VictorPrado99/poc-csv-persistence) to persist the data at a database.

Similar to the [persistence](https://github.com/VictorPrado99/poc-csv-persistence), the docker-compose in ./dev is enough to run this app isolated, running persistence to have funcional database.  
In the root directory is exist a docker-compose, to run a MySQL instance, [poc-csv-persistence](https://github.com/VictorPrado99/poc-csv-persistence) and [poc-csv-uploader](https://github.com/VictorPrado99/poc-csv-uploader)  

All images can be found at my docker hub 
- [docker-hub poc-csv-uploader](https://hub.docker.com/repository/docker/victorprado99/poc-csv-uploader)
- [docker-hub poc-csv-persistence](https://hub.docker.com/repository/docker/victorprado99/poc-csv-persistence)

As usual, the code itself have comments explaining my reasoning at each step