protobuf:
	protoc	--go_out=./pkg/proto  --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/proto --go-grpc_opt=paths=source_relative \
		--proto_path=./proto/proto \
		$$(find ./proto/proto -type f -iname "*.proto")

docker-protobuf:
	docker run --rm \
		-v $$(pwd)/proto:/proto \
		-v $$(pwd)/pkg:/pkg \
		generate-go-code \
			protoc	--go_out=/pkg/  --go_opt=paths=source_relative \
				--go-grpc_out=/pkg/ --go-grpc_opt=paths=source_relative \
				--proto_path=proto/proto \
				$$(find proto/proto -type f -iname "*.proto")

update-go-proto:
	git -C ./proto pull  \
		$(MAKE) protobuf \
		echo $(shell git show --pretty=format:%s -s HEAD -C ./proto)
		# git add . \
		# git commit -m "${git show --pretty=format:%s -s HEAD -C ./proto }"
		# git show --pretty=format:%s -s HEAD -C ./proto 
