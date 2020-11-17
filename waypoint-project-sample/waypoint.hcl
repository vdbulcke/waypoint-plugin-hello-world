# The name of your project. A project typically maps 1:1 to a VCS repository.
# This name must be unique for your Waypoint server. If you're running in
# local mode, this must be unique to your machine.
project = "custom-plugin-project-sample"

# Labels can be specified for organizational purposes.
# labels = { "foo" = "bar" }

# An application to deploy.
app "exec-app" {
    # Build specifies how an application should be deployed. In this case,
    # we'll build using a Dockerfile and keeping it in a local registry.
    build {
        use "hello-world-builder" {
            command = ["./build.sh"]
        }
        
        # Uncomment below to use a remote docker registry to push your built images.
        #
        # registry {
        #   use "docker" {
        #     image = "registry.example.com/image"
        #     tag   = "latest"
        #   }
        # }

    }

    # Deploy to Docker
    deploy {
        use "hello-world-builder" {
            command = ["echo", "hello"]
        }
    }
}
