from jinja2 import Environment, FileSystemLoader
from Crypto.PublicKey import RSA

if __name__ == '__main__':
    origin = input("What is the origin FQDN? Domain only, no scheme (e.g. www.horahora.org)")

    # compose file
    env = Environment(loader=FileSystemLoader("./templates"))
    template = env.get_template("docker-compose.yaml.j2")
    f_template = open("docker-compose.yaml, "W")
    f_template.write(template.render(build_images=True))

    # nginx config
    template = env.get_template("nginx.config.j2")
    f_template = open("docker-compose.yaml, "W")
    f_template.write(template.render(origin=origin))

    key = RSA.generate(2048, os.urandom)

    # env file
    template = env.get_template("env.j2")
    f_template = open(".env", "W")
    f_template.write(template.render(origin=origin, keypair=key.export_key("PEM")))

    # webapp env file *vomiting*
    template = env.get_template("env.webapp.j2")
    f_template = open("./webapp/.env", "W")
    f_template.write(template.render(origin=origin))
