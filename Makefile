test:
	echo '{"name": "Test", "age": 4}' | Authentication=test go run example/example.go store

test2:
	Authentication=test go run example/example.go store
