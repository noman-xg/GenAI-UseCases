# DrawTF

DrawTF is a go-based application that launches that allows you to interact with your data and generate terraform configs from json inputs. This codebase serves as the backend for terraform config generation through generative AI.

## Requirements
To run DrawTF, you need to have the following files in a folder named `scripts` next to the binary file:

- `chatportal.py`: This is the main script that handles the chat logic and the communication with the go server.
- `embeddings.py`: This is a helper script that generates embeddings for your data using a pre-trained model.
- `prompts.py`: This is a helper script that contains predefined prompts for different types of data.

You also need to have GO Python 3.7 and go1.21.0 or higher installed on your system:

You can install these packages using the command:

`pip install langchain[all] pandas openai requests`

## Usage
To create a binary for the drawTF, you can use the following command:
`go build https://github.com/noman-xg/GenAI-UseCases.git`
`go mod tidy`

To use DrawTF, you need to run the binary file with the following command:

`./drawTF`

This will start the go server on port 8082. You can then access the API Server by making HTTP POST requests to the following endpoints:

- `/message`: This endpoint accepts natural language queries as input and returns responses based on your data. You need to provide your data file (in CSV or JSON format) and your query as parameters in the request body. The response will be a JSON object with the following fields:
    - `status`: This indicates whether the query was successful or not.
    - `message`: This contains the natural language response to your query.

- `/start-gradio`: This endpoint launches a web-based chat portal that allows you to chat with your data using natural language. You need to provide your data file (in CSV or JSON format) as a parameter in the request body. The response will be a JSON object with the following field:
    - `url`: This contains a URL to access the chat portal.

- `/tfconfig`: This endpoint generates terraform configs from your data. You need to provide your data (in JSON format) and your request as parameters in the request body. The response will be a JSON object with the following fields:
    - `status`: This indicates whether the request was successful or not.
    - `content`: This contains the terraform code snippet generated from your request.
