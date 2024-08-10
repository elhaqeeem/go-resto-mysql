# Project Title
<h2> Resto Go (echo) mysql</h2>
## Description

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)


## Installation
To install and set up the project, follow these steps:

1. Clone the repository:
   ```sh
   git clone https://github.com/username/repository.git

## Usage

2. Copy env file 
   ```sh
    cp .env .env.backup

3. Running local --> delete comment in 
   ```sh

   //loadEnv() 
   
   And 

   ```sh
   
   /func loadEnv() {
	/err := godotenv.Load()
	//if err != nil {
	//	panic("Failed load env file")
	//}
   //}

4. Deploy to aws or etc --> upload or bulk environment in setting deployment

5. Command to running 
   ```sh
   go mod tidy
   go run main.go




