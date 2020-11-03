import os


class Environment(object):
    cli = "oc"

    def __init__(self, cli):
        self.set_cli(cli)

    def set_cli(self, cli):
        self.cli = cli

    def get_cli(self):
        return self.cli


# This is a global context (complementing behave's context)
# to be accesible from any place, even where behave's context is not available.
global ctx
ctx = Environment(os.getenv("TEST_ACCEPTANCE_CLI", "oc"))
