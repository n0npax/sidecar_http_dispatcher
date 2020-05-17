destroy:
	kubectl delete -f deployment.yaml
	kubectl delete cm nginx-conf

deploy:
	kubectl create cm nginx-conf --from-file configs/nginx.conf || true
	kubectl apply -f deployment.yaml
