# sports-news  👋
This application gets data from different providers and transforms from xml into a consistent 
and desirable format that can consume.

It also provides stability so that, as is often the case, when the external vendor has issues, 
we can still provide data to the applications, even if the data is outdated.

The requirements of this project have the following:

At regular intervals poll the endpoint for new news articles
● Transform the XML feeds of news articles into appropriate model(s) and save them in
the database.
● Provide two REST endpoints which return JSON:
    a. Retrieve a list of all articles
    b. Retrieve a single article by ID



### Setup:

 1.- Run follow command:

~~~bash
make setup
~~~


2.- Create a topic
~~~bash
make create_topics
~~~

### Documentation API

~~~bash
"/health"       --> return health of app
"/articles"     --> return a list of articles
"/articles/:id" --> return an article
~~~

### Testing

~~~bash
make test
~~~


#### 👨‍💻 Full list what has been used:
* [Kafka](https://github.com/segmentio/kafka-go) - Kafka library in Go
* [gin](https://github.com/gin-gonic/gin) - Web framework
* [MongoDB](https://github.com/mongodb/mongo-go-driver) - The Go driver for MongoDB
* [Docker](https://www.docker.com/) - Docker
* [Docker test](https://github.com/ory/dockertest/) - Docker test
* [uuid](https://github.com/google/uuid/) - uuid
* [cron](https://github.com/robfig/cron) - cron
* [testify](https://github.com/stretchr/testify) - testify
* [big_queue](https://github.com/patriciabonaldy/big_queue/) - My own library to Publisher/consumer in kafka or sqs

