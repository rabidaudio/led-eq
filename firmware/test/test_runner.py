import importlib.util
import sys
from pathlib import Path

from framework import test_runners


def import_from_path(module_name, file_path):
    spec = importlib.util.spec_from_file_location(module_name, file_path)
    module = importlib.util.module_from_spec(spec)
    sys.modules[module_name] = module
    spec.loader.exec_module(module)
    return module


def load_all_test_files():
    test_dir = Path(__file__).resolve().parent
    for test_module_path in test_dir.glob("test_*.py"):
        if str(test_module_path) == __file__:
            continue  # ignore current file
        module_name = test_module_path.name.removesuffix(".py")
        import_from_path(module_name, test_module_path)


if __name__ == "__main__":
    load_all_test_files()
    for runner in test_runners.values():
        runner()
