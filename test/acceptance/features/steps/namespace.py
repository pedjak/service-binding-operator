import re
from environment import ctx

from command import Command


class Namespace(object):
    def __init__(self, name):
        self.name = name
        self.cmd = Command()

    def create(self):
        output, exit_code = self.cmd.run(f"{ctx.cli} create namespace {self.name}")
        if re.search(f'namespace/{self.name} created', output) is not None or \
                re.search(rf'namespaces "{self.name}" already exists', output) is not None:
            return True
        else:
            print(f"Unexpected output when creating namespace: '{output}'")
        return False

    def is_present(self):
        _, exit_code = self.cmd.run(f'{ctx.cli} get ns {self.name}')
        return exit_code == 0
