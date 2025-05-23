generate: 
	uv run openapi-python-client generate --path ../api/v1.yaml --output-path oapi --overwrite 

test: generate
	uv run main.py
