
# <h1 align="center"> Olx clone(Backend) </h1>
___
<p align="center">
<a href="Java url">
    <img alt="Java" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
</a>
<a href="License url" >
        <img alt="BSD Clause 3" src="https://camo.githubusercontent.com/3dbcfa4997505c80ef928681b291d33ecfac2dabf563eb742bb3e269a5af909c/68747470733a2f2f696d672e736869656c64732e696f2f6769746875622f6c6963656e73652f496c65726961796f2f6d61726b646f776e2d6261646765733f7374796c653d666f722d7468652d6261646765"/>
    </a>
</p>

---

<p align="left">

## Overview
Olx clone, with Golang, PostgreSQL, Redis, docker, AWS S3, AWS SES, AWS SQS to achieve blazing performance, with features like product, user, seller, review, report-product, favorite and many more, this product compress of 25+ working APIs, also health API.


## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
   - [Environment Variables](#environment-variables)
- [Contributing](#contributing)
- [Contact](#contact)
- [License](#license)

## Technologies <a name="technologies"></a>

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Gin](https://gin-gonic.com)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [AWS S3](https://aws.amazon.com/)
- [AWS SES](https://aws.amazon.com/)
- [AWS SQS](https://aws.amazon.com/)

## Prerequisites

- [Docker](https://www.docker.com/)

## Features

- Olx clone with 25+ working APIs
- For scale we used AWS SES for email service
- For scale we used AWS S3 for file storage
- Performance DB we used PostgresQL 
- Dockerized for easy deployment

## Environment Variables

Table bellow shows the obligatory environment variables for mariadb container. You should set them based on what was also set for backend container.

Environment variable  | Default value | Optional
--- | --- | ---
STAGE | "" | `YES`
DB_HOST | http://127.0.0.1 | `YES`
DB_PORT | 5432 | `YES`
DB_USER | user | `YES`
DB_PASSWORD | postgres | `YES`
DB_NAME | olx-clone | `YES`
REDIS_HOST |  | `NO`
REDIS_PORT |  | `NO`
REDIS_USER |  | `NO`
REDIS_PASSWORD |  | `NO`
SENTRY_DSN |  | `NO`
DD_AGENT_HOST |  | `NO`
S3_BUCKET |  | `NO`

## Getting started
   ```
   First of all, correctly configure the Golang development environment on your machine, see https://go.dev/doc/install
   
   - Clone this repository:
   $ git clone https://github.com/swarajkumarsingh/olx-clone

   - Enter in directory:
   $ cd olx-clone

   - For install dependencies(optional):
   $ make install

   - Run the app: 
   $ ./run.sh
   ```

---

## DataBase Design
     
#### User Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | INT (Primary Key) | Unique identifier for each user |
| username    | VARCHAR(255) |  Unique identifier for each user                |
| fullname    | VARCHAR(255) | User's fullname                  |
| avatar       | TEXT | User's profile picture             |
| email       | VARCHAR(255) | User's email address             |
| password    | VARCHAR(255) | Securely hashed password          |
| location       | TEXT | User's profile picture             |
| coordinates       | TEXT | User's profile picture             |
| otp       | TEXT | User's profile picture             |
| otp_expiration  | TIMESTAMP | OTP timestamp    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
| updated_at  | TIMESTAMP | Timestamp of account modification|
     
#### Seller Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | INT (Primary Key) | Unique identifier for each seller |
| username    | VARCHAR(255) |  Unique identifier for each seller           |
| fullname    | VARCHAR(255) | Seller's fullname                  |
| description       | TEXT | Seller description             |
| is_verified       | BOOLEAN | Checks if the seller is verified             |
| avatar       | TEXT | Seller's profile picture             |
| phone       | VARCHAR(12) | Seller's phone             |
| email       | VARCHAR(100) | Seller's email address             |
| password    | VARCHAR(255) | Securely hashed password          |
| city    | VARCHAR(50) | Seller's city      |
| state    | VARCHAR(50) | Seller's state          |
| country    | VARCHAR(50) | Seller's country          |
| zip_code    | VARCHAR(50) | Seller's zip-code          |
| location       | TEXT | Seller's location             |
| coordinates       | TEXT | Seller's coordinates             |
| rating       | TEXT | Seller's rating             |
| account_status       | TEXT | Seller's account_status             |
| otp       | TEXT | Seller's otp           |
| otp_expiration  | TIMESTAMP | OTP timestamp    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
     
#### Product Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | INT (Primary Key) | Unique identifier for each seller |
| title    | VARCHAR(255) | Title's for the product |
| description       | TEXT | Seller description         |
| location       | TEXT | Seller's location       |
| coordinates       | TEXT | Seller's coordinates             |
| views       | BIGINT | Views count  |
| price       | VARCHAR(100) | Views count  |
| seller_id  | ID | Seller's ID    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
     
#### Review Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | ID (Primary Key) | Unique identifier for each seller |
| user_id  | ID | User's ID    |
| product_id  | ID | Products's ID    |
| rating  | TEXT | Products's rating    |
| comment  | TEXT | Products's comment    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
     
#### Seller Reports Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | ID (Primary Key) | Unique identifier for each seller |
| user_id  | ID | User's ID    |
| product_id  | ID | Products's ID    |
| message  | VARCHAR(100) | Report message    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
     
#### Favorites Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | ID (Primary Key) | Unique identifier for each seller |
| user_id  | ID | User's ID    |
| product_id  | ID | Products's ID    |
| created_at  | TIMESTAMP | Timestamp of account creation    |
     
#### Product Views Table

| Column Name | Data Type | Description                       |
| ----------- | --------- | --------------------------------- |
| id          | ID (Primary Key) | Unique identifier for each seller |
| user_id  | ID | User's ID    |
| product_id  | ID | Products's ID    |
| created_at  | TIMESTAMP | Timestamp of account creation    |

### Contributing
Contributions are welcome! If you'd like to contribute to this project, please follow these guidelines:
```
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and test thoroughly.
4. Commit your changes with clear commit messages.
5. Create a pull request against the main branch.
```


## Disclaimer

This application is a personal project built with educational and learning purposes in mind. It is neither affiliated nor endorsed by Amazon in any way. While the app features product details and images inspired by Amazon, these are solely for demonstration purposes and may not represent actual products. All rights to these elements belong to their respective owners. We are using them for educational purposes only and have no intention of commercial exploitation.

Additionally, be aware that any attempts to place orders within this prototype are purely for testing purposes and will not result in actual product deliveries or charges in the real-world. This environment is designated exclusively for simulation and development purposes

## Contact

- Swaraj Singh <br> <br>
  <a  href="https://www.linkedin.com/in/swarajkumarsingh/" target="_blank"><img alt="LinkedIn" src="https://img.shields.io/badge/linkedin%20-%230077B5.svg?&style=for-the-badge&logo=linkedin&logoColor=white" /></a>
  <a href="sswaraj169@gmail.com"><img  alt="Gmail" src="https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white" />

  feel free to contact me!

## License

> You can check out the full license [here](https://github.com/swarajkumarsingh/olx-clone/blob/main/LICENSE)

This project is licensed under the terms of the **MIT** license