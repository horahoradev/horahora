from jinja2 import Environment, FileSystemLoader
from Crypto.PublicKey import RSA
import os

if __name__ == '__main__':
    origin = input("What is the origin FQDN? Domain only, no scheme (e.g. www.horahora.org)")

    # compose file
    env = Environment(loader=FileSystemLoader("./templates"))
    template = env.get_template("docker-compose.yaml.j2")
    f_template = open("docker-compose.yaml", "w")
    f_template.write(template.render(build_images=True))
    print("Wrote docker-compose.yaml")

    # nginx config
    template = env.get_template("nginx.conf.j2")
    f_template = open("./configs/nginx.conf", "w")
    f_template.write(template.render(origin=origin))
    print("Wrote nginx.conf")


    # env file
    key = RSA.generate(2048, os.urandom)
    pem = key.export_key("PEM")
    template = env.get_template("env.j2")
    f_template = open(".env", "w")
    f_template.write(template.render(origin=origin, keypair=str(pem)[2:-1]))
    print("Wrote .env")

    # webapp env file *vomiting*
    template = env.get_template("env.webapp.j2")
    f_template = open("./webapp/.env", "w")
    f_template.write(template.render(origin=origin))
    print("Wrote ./webapp/.env")
