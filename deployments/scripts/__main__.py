import argparse
import images

argsParser = argparse.ArgumentParser(description="Build config flags")

argsParser.add_argument("--all", action="store_true", help="Run all scripts")
argsParser.add_argument("--build", "-b", action="store_true", help="Only build docker image")
argsParser.add_argument("--push", "-p", action="store_true", help="Only push docker image")

args = argsParser.parse_args()

if __name__ == "__main__":
    if args.all:
        images.build_docker_images()
        images.push_docker_images()
    else:
        if args.build: images.build_docker_images()
        if args.push: images.push_docker_images()