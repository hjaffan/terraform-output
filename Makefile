build:
	docker build --platform linux/amd64 -t hjaffan/terraform-output .
	docker push hjaffan/terraform-output