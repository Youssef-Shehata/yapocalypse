# Yapocalypse

## Overview

Yapocalypse is a microblogging platform inspired by Twitter, where users can "yap" (tweet) and follow other users.  
The platform currently focuses on text-based content and does not support videos, images, or likes. Users receive "yaps" from their followers in their feed.

## Status

This project is for practicing purposes only and is not ready for production use. It is also not fully functional yet and may contain bugs or security vulnerabilities.

if you encounter any issues or have any suggestions, please feel free to open an *issue* or *pull request*.

## Installation
1. Install ( Go, PostgreSQL, Redis, Kafka, Docker, Docker Compose) 
2. Clone the repository
3. Run `docker-compose up -d` to start the containers
4. Run `docker-compose exec postgres psql -U postgres -d postgres` to connect to the PostgreSQL database  
5. Read the [Env example](https://github.com/Youssef-Shehata/yapocalypse/blob/main/cmd/web/env_example) to create a `.env` file
6. Read the [Api Docs](https://github.com/Youssef-Shehata/yapocalypse/blob/main/cmd/web/API_DOCS.md) to understand the API endpoints

## Architecture

### Components

- **API Layer**: Handles HTTP requests and responses.
- **Kafka Producer**: Publishes "yap" events to Kafka.
- **Kafka Consumers**: Consumes "yap" events and processes them.
- **Redis**: Used for caching to improve performance.
- **JWT**: Used for authentication.
- **PostgreSQL**: The primary database for storing user data, "yaps," and relationships.

### Kafka Integration

#### Kafka Producer

- **Topic**: `yaps`

- **Action**: Publishes "yaps" to the `yaps` topic.

- **Key**: yap_id

- **Value**: yap_json

#### Kafka Consumers  

 *Feed Consumer*

- **Topic**: `yaps`

- **Action** : Updates the user's feed with new `yaps` from followed users.

 *Analytics Consumer*

- **Topic**: `yaps`

- **Action**: Analysis of `yaps` to determine popular topics.(not implemented yet)


### Redis Caching

- **user** : `{user_id}:feed`:  Caches the user's feed.

- **yap** : `{yap_id}:yap`: Caches individual `yaps`.

- **Cache Expiry** : Set to 5 minutes to balance between performance and data freshness.

### PostgreSQL Database

- **Users** : `id` , `email`, `username` , `password` , `created_at` , `updated_at`

- **Yaps** : `id` , `user_id` , `content` , `created_at` , `updated_at`

- **Followers** : `follower_id` , `following_id` , `created_at` , `updated_at`

- **Indexes**:  
  index on `Followers(followee_id) `
  index on `Feed(user_id)`
  index on `Yaps(user_id)`


## Security
- **Password Hashing**: BCrypt

- **JWT**: JSON Web Tokens
