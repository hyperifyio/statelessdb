# @hyperifyio/statelessdb

**StatelessDB** is a stateless, object-oriented cloud database server built to
securely process encrypted data stored on the client side.

With a zero-trust architecture, StatelessDB allows data to be stored by the 
client or third-party services while compute servers perform secure operations
without retaining any state. Using a shared secret for decryption, it ensures 
sensitive data remains protected and inaccessible to the client. 

StatelessDB is a scalable, flexible solution for cloud-native environments, 
offering secure computation on both public and private object properties—without 
persistent server-side storage.

## Starting the development server

```bash
docker-compose build && docker-compose up
```

Once started, the server is available at http://localhost:8080/api/v1

## Starting the server from localhost

You can start the server locally like this:

```
PRIVATE_KEY=9ca549e8e80e363cb92b99936dd869c65eca7f474d2b595a72d5e9a2d79eff61 ./statelessdb
```

## Manual testing with Curl

### Creating a resource without public data

Request body:

```json
{}
```

Command: 

```bash
curl -i -d '{}' http://localhost:8080/api/v1
```

Response headers:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 08 Oct 2024 23:48:53 GMT
Content-Length: 440
```

Response body:

```json
{
  "id":"d626cac1-da23-4c67-9001-7bb03a40e90e",
  "owner":"91b1ab41-4a73-488f-89fd-c3119b349361",
  "created":"2024-10-08T23:48:53Z",
  "updated":"2024-10-08T23:48:53Z",
  "public":null,
  "private":"8N1svYP/KbElP84uLI2Ch3wck8jBdQIa+4QUW1G6O+QabgWkupNM99NkxSlw5n0dvum+7lMrGwtrFDvIJgh2bXZLMI6vyDX6VKl4XZSds5z/zCH0QjNG+sSVN+nLY6GR1iJctJwRJzuNtpe4mGl+IYBR5xrnV3VGQ9/BrEkhoErtDuxsQd2ES0yd7JiP6JAnnZVH3V95/MZqfNcJhfHYViXKKk3OF8rbCGROcfzsFlhPany0LiUgkHJl9A+a1MM3"
}
```

Notice, that the `public` is null. This is a performance optimization: we will 
not automatically allocate an object unless you provide one.


### Creating a resource without public data

Request body:

```json
{
  "hello": "world"
}
```

Command: 

```bash
curl -i -d '{"public":{"hello":"world"}}' http://localhost:8080/api/v1
```

Response headers:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 09 Oct 2024 12:36:49 GMT
Content-Length: 479
```

Response body:

```json
{
  "id":"d626cac1-da23-4c67-9001-7bb03a40e90e",
  "owner":"91b1ab41-4a73-488f-89fd-c3119b349361",
  "created":"2024-10-08T23:48:53Z",
  "updated":"2024-10-08T23:48:53Z",
  "public": {
    "hello": "world"
  },
  "private":"8N1svYP/KbElP84uLI2Ch3wck8jBdQIa+4QUW1G6O+QabgWkupNM99NkxSlw5n0dvum+7lMrGwtrFDvIJgh2bXZLMI6vyDX6VKl4XZSds5z/zCH0QjNG+sSVN+nLY6GR1iJctJwRJzuNtpe4mGl+IYBR5xrnV3VGQ9/BrEkhoErtDuxsQd2ES0yd7JiP6JAnnZVH3V95/MZqfNcJhfHYViXKKk3OF8rbCGROcfzsFlhPany0LiUgkHJl9A+a1MM3"
}
```

### Using the resource

Request body:

```json
{
  "payload":{
    "id":"d626cac1-da23-4c67-9001-7bb03a40e90e",
    "owner":"91b1ab41-4a73-488f-89fd-c3119b349361",
    "created":"2024-10-08T23:48:53Z",
    "updated":"2024-10-08T23:48:53Z",
    "public": null,
    "private":"8N1svYP/KbElP84uLI2Ch3wck8jBdQIa+4QUW1G6O+QabgWkupNM99NkxSlw5n0dvum+7lMrGwtrFDvIJgh2bXZLMI6vyDX6VKl4XZSds5z/zCH0QjNG+sSVN+nLY6GR1iJctJwRJzuNtpe4mGl+IYBR5xrnV3VGQ9/BrEkhoErtDuxsQd2ES0yd7JiP6JAnnZVH3V95/MZqfNcJhfHYViXKKk3OF8rbCGROcfzsFlhPany0LiUgkHJl9A+a1MM3"
  }
}
```

Command:

```bash
curl -i \
  -d '{"payload":{"id":"d626cac1-da23-4c67-9001-7bb03a40e90e","owner":"91b1ab41-4a73-488f-89fd-c3119b349361","created":"2024-10-08T23:48:53Z","updated":"2024-10-08T23:48:53Z","public":{},"private":"8N1svYP/KbElP84uLI2Ch3wck8jBdQIa+4QUW1G6O+QabgWkupNM99NkxSlw5n0dvum+7lMrGwtrFDvIJgh2bXZLMI6vyDX6VKl4XZSds5z/zCH0QjNG+sSVN+nLY6GR1iJctJwRJzuNtpe4mGl+IYBR5xrnV3VGQ9/BrEkhoErtDuxsQd2ES0yd7JiP6JAnnZVH3V95/MZqfNcJhfHYViXKKk3OF8rbCGROcfzsFlhPany0LiUgkHJl9A+a1MM3"}}' \
  http://localhost:8080/api/v1
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 08 Oct 2024 23:49:15 GMT
Content-Length: 440
```

```json
{
  "id":"d626cac1-da23-4c67-9001-7bb03a40e90e",
  "owner":"91b1ab41-4a73-488f-89fd-c3119b349361",
  "created":"2024-10-08T23:48:53Z",
  "updated":"2024-10-08T23:49:15Z",
  "public": null,
  "private":"cIkAbZ/rnUbafQSbUiDcdWE8DHGQfbPMU8QuWPot6JTTehppqdGFkR9NLYE/ctpNkpHtumcI88WNIO+DSuhTTmzFkr1jIaL6eF6/tbp98/nHVYHXDg/+txGqhkjnylVOi5VNqgPNLfJI6qoxow3AcsdlL89MJmtr28ocPijAH29ZXmDQMn5+EEFFovHpJ0jBuTjB/tmMsls1NW9FgxL7wWlqWVk8R3/9gBS+GMSdAvmcg9NLSzSw9nR5YkDEckO6"
}
```
