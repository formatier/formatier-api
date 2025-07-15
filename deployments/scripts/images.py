import subprocess
import sys
import utils

infrastructure_servers = ["gateway-server", "orchestrator-server", "auth-server"]
service_servers = ["document-service", "plugin-service", "user-service", "workspace-service"]

current_version = "v1"

def build_docker_images():
    """
    Builds Docker images for infrastructure and service servers.
    """

    print("Starting Docker image build process")

    # Build infrastructure server images
    for server in infrastructure_servers:
        dockerfile_path = f"./{server}/Dockerfile"
        image_tag = f"ghcr.io/formatier/formatier-api/{server}:{current_version}"
        command = ["docker", "build", "-t", image_tag, "-f", dockerfile_path, "."]
        print(f"Building image: {image_tag} from {dockerfile_path}")
        try:
            subprocess.run(command, check=True, shell=True)
            print(f"Successfully built {image_tag}")
        except subprocess.CalledProcessError as e:
            utils.printError(e)
            sys.exit(1)

    # Build service server images
    for service in service_servers:
        dockerfile_path = f"./services/{service}/Dockerfile"
        image_tag = f"ghcr.io/formatier/formatier-api/{service}:{current_version}"
        command = ["docker", "build", "-t", image_tag, "-f", dockerfile_path, "."]
        print(f"Building image: {image_tag} from {dockerfile_path}")
        try:
            subprocess.run(command, check=True, shell=True)
            print(f"Successfully built {image_tag}")
        except subprocess.CalledProcessError as e:
            utils.printError(e)
            sys.exit(1)

    print("Successfully built all server images!")

def push_docker_images():
    for server in infrastructure_servers: pass
    for service in service_servers: pass