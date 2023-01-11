# Terraform Backend Surreal
This is a Terraform backend that uses SurrealDB to store state.

## Usage
### Server
This backend is configured via a `config.yaml` file that it expects to be in the root of the project.

Here's an example of one
    
```yaml
---
port: 8032
surreal:
  user: root
  password: root
  address: "ws://localhost:8000/rpc"
  namespace: terraformbackend
  database: terraformbackend
users:
  terraform: alligator3
```

Port: The port that the backend will listen on
Surreal: The SurrealDB connection information
Users: A map of username to password that will be used to authenticate with the backend

### Terraform
After the server is setup, you can include the backend in your Terraform configuration.

In Terraform...
```hcl
terraform {
  backend "http" {
    address = "http://localhost:8032/myproject/mystack"
    lock_address = "http://localhost:8032/myproject/mystack"
    unlock_address = "http://localhost:8032/myproject/mystack"
    username = "terraform"
    password = "alligator3"
    skip_cert_verification = true
  }
}
```

In CDKTF...
```typescript
new HttpBackend(this, {
    address: "http://localhost:8032/myproject/mystack",
    lockAddress: "http://localhost:8032/myproject/mystack",
    unlockAddress: "http://localhost:8032/myproject/mystack",
    username: "terraform",
    password: "alligator3",
    skipCertVerification: true,
})
```

The uri is a combination of the project name and the stack name. This is so there aren't potential conflicts between stacks across different projects.

## Development
This is still in development. The base functionality is here but there could be some changes introduced in the future depending on needs.

### Todo List
- [ ] Add tests
- [ ] Build and publish container
- [ ] CDKTF Construct for backend
- [ ] Deployment examples (Kubernetes, Docker, etc.)
- [ ] More usage examples (Terraform, CDKTF, etc.)

### Maybe Someday
- [ ] TLS support
