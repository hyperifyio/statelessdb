# @hyperifyio/statelessdb

**StatelessDB** is a stateless, object-oriented cloud database server built to
securely process encrypted data stored on the client side. With a zero-trust 
architecture, StatelessDB allows data to be stored by the client or third-party 
services while compute servers perform secure operations without retaining any
state. Using a shared secret for decryption, it ensures sensitive data remains 
protected and inaccessible to the client. StatelessDB is a scalable, flexible 
solution for cloud-native environments, offering secure computation on both 
public and private object properties—without persistent server-side storage.

## Starting the development server

```bash
docker-compose build && docker-compose up
```

Once started, the server is available at http://localhost:8080
and API is available at http://localhost:3001

## Starting the server from localhost

You can start the server locally like this:

```
PRIVATE_KEY=9ca549e8e80e363cb92b99936dd869c65eca7f474d2b595a72d5e9a2d79eff61 ./statelessdb
```

## Manual testing with Curl

### Requesting data

Request body:

```json
{
  "nextIndex": 0, 
}
```

Command: 

```bash
curl -i -d '{"nextIndex": 0}' http://localhost:3001
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 07 Apr 2024 23:41:23 GMT
Content-Length: 436
```

```json
{
  "score":0,
  "private":"Y1aU4hhRmV1Puc+05E/apHN5gdAaaXI9px4fOEcvkiGYG4Po6Gp8eO+ZzUduanmX0yWPY6ChVZk+TW7QPO99XgtcoqJjHAwJy0EcV5v54elYflk1Ltr9kBCQqQEP5Tf3WsmB+zinXaFxr6Jkc+mDLLY/VKqMVmkP/qELOLVnMOnuxkCXdzXONYAOYU0u7IEMRtB2lC6fvAjoy6s9wWJFWvp526aFcAnTUN31gIWJbWI6nJu92WJIVu0+wxs9E8AbOBhEG0hXpM72hmH8bBXml5s8Z9S9UxMLpv8ZqaZd5fzCLN1G4ctuLmUC/f5fKhJLAGMhHyMnYLL6zgaf8FkdbvQn3DL/9F1dMmdb",
}
```

### Continuing a data

Request body:

```json
{
  "nextIndex": 1, 
  "private":"Y1aU4hhRmV1Puc+05E/apHN5gdAaaXI9px4fOEcvkiGYG4Po6Gp8eO+ZzUduanmX0yWPY6ChVZk+TW7QPO99XgtcoqJjHAwJy0EcV5v54elYflk1Ltr9kBCQqQEP5Tf3WsmB+zinXaFxr6Jkc+mDLLY/VKqMVmkP/qELOLVnMOnuxkCXdzXONYAOYU0u7IEMRtB2lC6fvAjoy6s9wWJFWvp526aFcAnTUN31gIWJbWI6nJu92WJIVu0+wxs9E8AbOBhEG0hXpM72hmH8bBXml5s8Z9S9UxMLpv8ZqaZd5fzCLN1G4ctuLmUC/f5fKhJLAGMhHyMnYLL6zgaf8FkdbvQn3DL/9F1dMmdb",
}
```

Command:
```bash
curl -i -d '{"nextIndex": 1, "gameState": {"score":0,"cards":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"private":"Y1aU4hhRmV1Puc+05E/apHN5gdAaaXI9px4fOEcvkiGYG4Po6Gp8eO+ZzUduanmX0yWPY6ChVZk+TW7QPO99XgtcoqJjHAwJy0EcV5v54elYflk1Ltr9kBCQqQEP5Tf3WsmB+zinXaFxr6Jkc+mDLLY/VKqMVmkP/qELOLVnMOnuxkCXdzXONYAOYU0u7IEMRtB2lC6fvAjoy6s9wWJFWvp526aFcAnTUN31gIWJbWI6nJu92WJIVu0+wxs9E8AbOBhEG0hXpM72hmH8bBXml5s8Z9S9UxMLpv8ZqaZd5fzCLN1G4ctuLmUC/f5fKhJLAGMhHyMnYLL6zgaf8FkdbvQn3DL/9F1dMmdb","lastCard":1,"lastIndex": 15}}' http://localhost:3001
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 07 Apr 2024 23:42:07 GMT
Content-Length: 436
```

```json
{
  "score":0,
  "cards":[0,4,0,0,0,0,0,0,0,0,0,0,0,0,0,4],
  "private":"YTUi5aJQjyCqbEBVf1gxNB+KZmkXnzQg0jiqfkJfECE18esg+q+hODgwz3s0lNb5v7oTLRanO/VK22Ppl3zAcCEk6aObarfBJMeGgcw7RdWXUZ19f6pBgz9rp1baVM9CkmPyc/kqqdEZOKFms89dzefhKbY/tUqP6IwLQ5Se3zHT6DsI0YAjbx2JLWcbwUW17vRWMkNibuNpFVgC4H6UwPQLnvNJkJGRjq8Zl6t/xaUIhaEyLMRMF0nVuO6aHQfgel6W/tDMDN2e8CJIlOMpu9zWcJDaRQb+p9Ojk2GtSsQUm90ectElWy/gQ66Rgi8B6mi5hby3kGS6Y8KSBNKbR05F9Sr/sEf196ff",
  "lastCard": 4,
  "lastIndex": 1
}
```
