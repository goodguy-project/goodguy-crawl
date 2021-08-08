import os
import sys


def build():
    root = os.path.dirname(os.path.abspath(__file__))
    os.chdir(root)
    pb_path = os.path.join(root, 'crawl_service')
    crawl_service_proto_path = os.path.join(pb_path, 'crawl_service.proto')
    os.system(f'{sys.executable} -m grpc_tools.protoc -I{root} --python_out={root} --grpc_python_out={root} '
              f'{crawl_service_proto_path}')


if __name__ == '__main__':
    build()
