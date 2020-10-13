# databaseRestApi

  databaseRestApi/cmd/servicemachine is a CLI tool that starts an API endpoint (/cpu, /storage or /memory) that provides access to cpu, storage or memory of the machine in a RESTful manner

# terraform
To deploy in terraform you need:
1. "terraform.tfvars" file with:
  access_key = value
  secret_key = value
2. terraform cli
3. terraform apply -auto-approve
4. Access the ip from the output of terraform command
  elastic-ip-for-web-server = $ip
5. Go to $ip:8989/cpu or $ip:8989/storage or $ip:8989/memory
