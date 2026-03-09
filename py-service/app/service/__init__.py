from grpc.aio import server
from grpc import ssl_server_credentials
from os import environ
from pathlib import Path

from payslip.v1.payslip_pb2_grpc import add_DocumentServiceServicer_to_server

from app.service.payslip import DocumentServer

async def main():
    svc = server()

    if tls_path := environ.get("TLS_PATH"):

        certs = Path(tls_path)

        private_key = (certs / "server.key").read_bytes()
        cert = (certs / "server.crt").read_bytes()
        

        creds = ssl_server_credentials(
            [(private_key, cert)]
        )

        svc.add_secure_port(
            "[::]:55443",
            creds
        )
        add_DocumentServiceServicer_to_server(
            DocumentServer(),
            svc
        )

        
        await svc.start()
        print("gRPC TLS server started on :50051")
        await svc.wait_for_termination()

    else:
        raise ValueError("Cannot find tls certificates")